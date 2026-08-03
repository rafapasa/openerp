// scripts/remove_audit_hooks.go
package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	modelDir := "internal/models"

	fmt.Println("🔍 Removendo hooks de auditoria dos models...")

	// Backup
	backupDir := fmt.Sprintf("internal/models_backup_%s", time.Now().Format("20060102_150405"))
	fmt.Printf("📁 Criando backup em: %s\n", backupDir)

	if err := copyDir(modelDir, backupDir); err != nil {
		fmt.Printf("❌ Erro ao criar backup: %v\n", err)
		os.Exit(1)
	}

	filesProcessed := 0
	filesModified := 0

	err := filepath.Walk(modelDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		filesProcessed++
		fmt.Printf("📝 Analisando: %s\n", path)

		modified, err := processFile(path)
		if err != nil {
			fmt.Printf("⚠️  Erro ao processar %s: %v\n", path, err)
			return nil
		}

		if modified {
			filesModified++
			fmt.Printf("✅ Removido: %s\n", path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("❌ Erro: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✨ Processo concluído!\n")
	fmt.Printf("📊 Arquivos analisados: %d\n", filesProcessed)
	fmt.Printf("📊 Arquivos modificados: %d\n", filesModified)
	fmt.Printf("📁 Backup criado em: %s\n", backupDir)
	fmt.Println("\n⚠️  Verifique os arquivos e execute 'go build ./...' para testar")
}

func processFile(filePath string) (bool, error) {
	// Ler o conteúdo do arquivo
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	// Verificar se tem hooks de auditoria
	hasHook := false
	if strings.Contains(string(content), "BeforeCreate") &&
		strings.Contains(string(content), "CreatedBy") {
		hasHook = true
	}

	if !hasHook {
		return false, nil
	}

	// Parse do arquivo
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return false, fmt.Errorf("erro ao parsear: %w", err)
	}

	// Remover os hooks
	newNode := removeAuditHooksFromNode(node)

	// Escrever arquivo
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, newNode); err != nil {
		return false, fmt.Errorf("erro ao formatar: %w", err)
	}

	if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
		return false, fmt.Errorf("erro ao escrever: %w", err)
	}

	return true, nil
}

func removeAuditHooksFromNode(node *ast.File) *ast.File {
	var newDecls []ast.Decl
	removedCount := 0

	for _, decl := range node.Decls {
		// Verifica se é uma função
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			newDecls = append(newDecls, decl)
			continue
		}

		// Verifica se é BeforeCreate ou BeforeUpdate (qualquer receiver)
		if fn.Name.Name == "BeforeCreate" || fn.Name.Name == "BeforeUpdate" {
			// Verificar se a função trata CreatedBy/UpdatedBy
			isAuditHook := false
			if fn.Body != nil {
				ast.Inspect(fn.Body, func(n ast.Node) bool {
					// Procura por CreatedBy ou UpdatedBy
					if sel, ok := n.(*ast.SelectorExpr); ok {
						// Verifica se é m.CreatedBy, p.CreatedBy, etc.
						if _, ok := sel.X.(*ast.Ident); ok {
							if sel.Sel.Name == "CreatedBy" || sel.Sel.Name == "UpdatedBy" {
								isAuditHook = true
								return false
							}
						}
					}
					// Procura por new(int) - padrão comum nos hooks
					if call, ok := n.(*ast.CallExpr); ok {
						if fun, ok := call.Fun.(*ast.Ident); ok {
							if fun.Name == "new" && len(call.Args) > 0 {
								if sel, ok := call.Args[0].(*ast.Ident); ok {
									if sel.Name == "int" {
										// Verifica se está sendo atribuído a CreatedBy/UpdatedBy
										if parent, ok := n.(*ast.AssignStmt); ok {
											for _, lhs := range parent.Lhs {
												if sel2, ok := lhs.(*ast.SelectorExpr); ok {
													if sel2.Sel.Name == "CreatedBy" || sel2.Sel.Name == "UpdatedBy" {
														isAuditHook = true
														return false
													}
												}
											}
										}
									}
								}
							}
						}
					}
					return true
				})
			}

			if isAuditHook {
				fmt.Printf("   🗑️  Removendo função: %s (receiver: %s)\n",
					fn.Name.Name, fn.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name)
				removedCount++
				continue // Pula esta declaração
			}
		}

		newDecls = append(newDecls, decl)
	}

	node.Decls = newDecls
	return node
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(destPath, data, info.Mode())
	})
}

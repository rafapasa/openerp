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
	err := copyDir(modelDir, backupDir)
	if err != nil {
		fmt.Printf("❌ Erro ao criar backup: %v\n", err)
		os.Exit(1)
	}

	// Processar arquivos
	err = filepath.Walk(modelDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.Contains(path, "_test.go") {
			processFile(path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("❌ Erro: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✨ Processo concluído!")
}

func processFile(filePath string) {
	// Parse do arquivo
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("❌ Erro ao parsear %s: %v\n", filePath, err)
		return
	}

	// Verificar se tem hooks de auditoria
	hasAuditHooks := false
	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if strings.HasPrefix(fn.Name.Name, "BeforeCreate") ||
				strings.HasPrefix(fn.Name.Name, "BeforeUpdate") {
				// Verificar se contém appcontext.GetUserID
				hasGetUserID := false
				ast.Inspect(fn.Body, func(n ast.Node) bool {
					if call, ok := n.(*ast.CallExpr); ok {
						if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
							if ident, ok := sel.X.(*ast.Ident); ok {
								if ident.Name == "appcontext" && sel.Sel.Name == "GetUserID" {
									hasGetUserID = true
									return false
								}
							}
						}
					}
					return true
				})
				if hasGetUserID {
					hasAuditHooks = true
					return false
				}
			}
		}
		return true
	})

	if hasAuditHooks {
		fmt.Printf("📝 Processando: %s\n", filePath)

		// Remover os hooks
		newNode := removeAuditHooksFromNode(node)

		// Escrever arquivo
		var buf bytes.Buffer
		if err := format.Node(&buf, fset, newNode); err != nil {
			fmt.Printf("❌ Erro ao formatar: %v\n", err)
			return
		}

		if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
			fmt.Printf("❌ Erro ao escrever: %v\n", err)
			return
		}

		fmt.Printf("✅ Removido: %s\n", filePath)
	}
}

func removeAuditHooksFromNode(node *ast.File) *ast.File {
	// Remover funções BeforeCreate e BeforeUpdate
	var newDecls []ast.Decl
	for _, decl := range node.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if strings.HasPrefix(fn.Name.Name, "BeforeCreate") ||
				strings.HasPrefix(fn.Name.Name, "BeforeUpdate") {
				// Verificar se contém appcontext.GetUserID
				hasGetUserID := false
				ast.Inspect(fn.Body, func(n ast.Node) bool {
					if call, ok := n.(*ast.CallExpr); ok {
						if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
							if ident, ok := sel.X.(*ast.Ident); ok {
								if ident.Name == "appcontext" && sel.Sel.Name == "GetUserID" {
									hasGetUserID = true
									return false
								}
							}
						}
					}
					return true
				})
				if hasGetUserID {
					continue // Pular este decl
				}
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

		relPath, _ := filepath.Rel(src, path)
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

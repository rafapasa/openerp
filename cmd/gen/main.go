package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	// Configuração da conexão
	dsn := "openerp:1234@tcp(localhost:3306)/etools_openerp?charset=utf8mb4&parseTime=True&loc=Local"

	// Conectar ao banco
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco: %v", err)
	}

	// Configurar o gerador
	g := gen.NewGenerator(gen.Config{
		OutPath:           "./internal/models",
		Mode:              gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable:     false,
		FieldCoverable:    false,
		FieldSignable:     false,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	// Usar a conexão
	g.UseDB(db)

	// Gerar todos os modelos
	g.ApplyBasic(
		// Tabelas principais
		g.GenerateModel("usuario"),
		g.GenerateModel("grupo_usuario"),
		g.GenerateModel("entidade"),
		g.GenerateModel("entidade_endereco"),
		g.GenerateModel("entidade_formacontato"),
		g.GenerateModel("produto"),
		g.GenerateModel("produto_grupo"),
		g.GenerateModel("produto_estoque"),
		g.GenerateModel("documento_venda"),
		g.GenerateModel("documento_venda_item"),
		g.GenerateModel("documento_venda_pagamento"),
		g.GenerateModel("forma_pagamento"),
		g.GenerateModel("condicao_pagamento"),
		g.GenerateModel("tabela_preco"),
		g.GenerateModel("empresa_filial"),
		g.GenerateModel("processo"),
		g.GenerateModel("operacaofiscal"),

		// Modelos específicos para autenticação
		g.GenerateModel("usuario_filial"),
		g.GenerateModel("seguranca_recurso_grupo"),
	)

	// Executar o gerador
	g.Execute()

	fmt.Println("✅ Models gerados com sucesso!")
}

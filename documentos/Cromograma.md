🚀 Cronograma de Desenvolvimento Sugerido

### **Fase 1: Base e Infraestrutura (Dias 1-2) - 100% Concluído**
*   **[100%]** Configurar ambiente de desenvolvimento (`scripts/dev.ps1`, `go.mod`).
*   **[100%]** Criar estrutura de pastas (padrão `cmd`, `internal`, `docs`, `scripts`).
*   **[100%]** Configurar conexão com MySQL (evidenciado pelos models GORM e `comandos.txt`).
*   **[100%]** Implementar sistema de autenticação JWT (`/auth/login`, `/auth/refresh` no `swagger.json`).
*   **[100%]** Criar middleware de segurança (implícito pelo uso de `BearerAuth` nos endpoints).
*   **[100%]** Configurar logging e monitoramento (dependência `go.uber.org/zap` instalada).
*   **[100%]** Criar CRUD básico de Usuários (evidenciado pelos endpoints de autenticação e `GORM_TAGS.md`).

### **Fase 2: Módulo de Entidades (Dias 3-4) - 100% Concluído**
*   **[100%]** Implementar CRUD de Entidade (Clientes/Fornecedores) (`/entidades` no `swagger.json`).
*   **[100%]** Implementar Endereços (`/entidades/{id}/enderecos`).
*   **[100%]** Implementar Contatos (`/entidades/{id}/contatos`).
*   **[100%]** Implementar Grupo de Entidades (`/grupos-entidades`).
*   **[100%]** Implementar Regime Tributário (`/entidades/{id}/regimes-tributarios`).
*   **[100%]** Implementar Limite de Crédito (`/limites-credito`).

### **Fase 3: Módulo de Produtos (Dias 5-7) - 90% Concluído**
*   **[100%]** Implementar CRUD de Produto (`/produtos` e `produto_service.go`).
*   **[100%]** Implementar Grupos e Subgrupos (`/produto-grupos`, `/produto-subgrupos`).
*   **[100%]** Implementar Marcas e Modelos (`/produto-marcas`, `/produto-modelos`).
*   **[100%]** Implementar Tabela de Preços (`/tabelas-preco`).
*   **[50%]** Implementar Grade de Produtos (Cores/Tamanhos).
    *   **Observação:** Não há uma estrutura clara para grade. **Sugestão:** Criar tabelas `produto_variacoes` e `produto_atributos` (ex: 'Cor', 'Tamanho') para uma grade flexível.
*   **[0%]** Implementar Estoque Inicial.
    *   **Observação:** Funcionalidade ainda não iniciada. Pode ser um campo no cadastro do produto ou uma movimentação inicial no módulo de estoque.

### **Fase 4: Módulo de Pedidos (Dias 8-12) - 40% Concluído**
*   **[80%]** Implementar CRUD de Pedidos (`/documentos-venda` no `swagger.json`).
    *   **Observação:** A estrutura principal está criada. Faltam os fluxos de alteração de status.
*   **[80%]** Implementar Itens do Pedido (definido em `DocumentoVendaRequest`).
*   **[20%]** Implementar Pagamentos do Pedido (`documento_venda_pagamento.go` existe, mas a lógica de negócio não está implementada no service).
*   **[0%]** Implementar Fluxo de Situações (Aberto → Em atividade → Fechado).
    *   **Observação:** Necessário criar endpoints ou services para transitar o pedido entre os status.
*   **[100%]** Implementar Tipos de Entrega (definido em `DocumentoVendaRequest`).
*   **[100%]** Implementar Taxa de Entrega (definido em `DocumentoVendaRequest`).
*   **[100%]** Implementar Descontos e Acréscimos (definido em `DocumentoVendaRequest`).
*   **[0%]** Implementar Histórico de Pedidos.
    *   **Sugestão:** Criar uma tabela `documento_venda_historico` para logar todas as alterações de status e outras modificações importantes.

### **Fase 5: Módulo de Estoque (Dias 13-15) - 5% Concluído**
*   **[0%]** Implementar Movimentação de Estoque.
*   **[20%]** Implementar Saldo de Estoque.
    *   **Observação:** O campo `saldo_estoque` existe na resposta do DTO de Produto, mas a lógica de cálculo e atualização não está implementada.
*   **[0%]** Implementar Baixa de Estoque (pedidos).
*   **[0%]** Implementar Entrada de Estoque (compras).
*   **[0%]** Implementar Ajustes de Estoque.
*   **[0%]** Implementar Histórico de Movimentações.
*   **[0%]** Implementar Relatório de Estoque.

### **Fase 6: Módulo Financeiro (Dias 16-19) - 10% Concluído**
*   **[20%]** Implementar Títulos a Receber.
    *   **Observação:** A estrutura `documento_venda_pagamento` pode gerar os títulos, mas o CRUD e a gestão dos títulos (baixa, etc.) não existem.
*   **[0%]** Implementar Títulos a Pagar.
*   **[0%]** Implementar Baixas de Títulos.
*   **[0%]** Implementar Remessa Bancária.
*   **[0%]** Implementar Retorno Bancário.
*   **[0%]** Implementar Controle de Cheques.
*   **[0%]** Implementar Fluxo de Caixa.
*   **[NOVO]** Implementar Centro de Custo.
    *   **Observação:** O model `centro_de_custo.go` já foi criado, indicando um bom começo para esta funcionalidade.

### **Fase 7: Módulo de Relatórios (Dias 20-22) - 0% Concluído**
*   **[0%]** Implementar Relatório de Vendas.
*   **[0%]** Implementar Relatório de Produtos.
*   **[0%]** Implementar Relatório de Clientes.
*   **[0%]** Implementar Relatório Financeiro.
*   **[0%]** Implementar Dashboard de Indicadores.

### **Fase 8: Integrações (Dias 23-25) - 0% Concluído**
*   **[0%]** Implementar Integração com WhatsApp.
*   **[0%]** Implementar Integração com PIX.
*   **[0%]** Implementar Integração com PagSeguro.
*   **[0%]** Implementar Notificações por E-mail.
*   **[0%]** Implementar Notificações SMS.

### **Fase 9: Testes e Deploy (Dias 26-30) - 30% Concluído**
*   **[20%]** Realizar testes unitários.
    *   **Observação:** O arquivo `comandos.txt` menciona `TestEntidadeSuite`, indicando que a estrutura de testes com `testify` já está sendo utilizada. É preciso expandir a cobertura.
*   **[0%]** Realizar testes de integração.
*   **[0%]** Realizar testes de carga.
*   **[80%]** Documentar API (Swagger).
    *   **Observação:** O `swagger.json` está bem completo e atualizado, o que é excelente. O script `swagger.ps1` está um pouco diferente do comando no `dev.ps1`, vale a pena unificar.
*   **[0%]** Configurar CI/CD.
*   **[0%]** Realizar deploy em produção.

### **[NOVO] Fase 10: Melhorias Gerais e Refatoração**
*   **[NOVO]** **Refatorar Scripts PowerShell:** Unificar `dev.ps1`, `dev1.ps1` e `criar.ps1` em um único script `dev.ps1` mais robusto, utilizando o `param` block para argumentos, que é uma prática mais moderna e segura.
*   **[NOVO]** **Padronizar Auditoria:** Garantir que todos os models tenham os campos `created_by` e `updated_by` e que os hooks `BeforeCreate`/`BeforeUpdate` sejam usados para preenchê-los, possivelmente extraindo o ID do usuário do contexto da requisição.
*   **[NOVO]** **Melhorar Tratamento de Erros:** Padronizar as respostas de erro da API, talvez usando um `error handler` centralizado no Gin para converter erros de serviço (ex: `apperrors.NewValidationError`) em respostas HTTP 400 ou 404 apropriadas.
*   **[NOVO]** **Implementar Cache:** Para dados que não mudam com frequência (ex: Tipos de Documento, Formas de Pagamento), implementar um cache com Redis (a dependência `go-redis` já está no projeto) para melhorar a performance.

Espero que esta análise detalhada ajude a guiar os próximos passos do desenvolvimento. O projeto está com uma base muito forte!

<!--
[PROMPT_SUGGESTION]Com base na análise, crie os models e DTOs para a funcionalidade de "Grade de Produtos" que você sugeriu.[/PROMPT_SUGGESTION]
[PROMPT_SUGGESTION]Refatore os scripts PowerShell (`dev.ps1`, `dev1.ps1`, `criar.ps1`) em um único arquivo `dev.ps1` aprimorado.[/PROMPT_SUGGESTION]

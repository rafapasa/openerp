// ============================================================
// PACOTE E IMPORTAÇÕES
// ============================================================

// 1. Declaração do pacote
// Todo arquivo Go começa com 'package nome'
// 'models' é o nome do pacote onde este arquivo está
package models

// 2. Importações
// Importamos os pacotes que vamos usar neste arquivo
import (
	"time" // Pacote para trabalhar com datas e horas
)

// ============================================================
// DEFINIÇÃO DO STRUCT (MODEL)
// ============================================================

// 3. Definição do struct Usuario
// Um struct é como uma "receita" que define a estrutura de um objeto
// O nome do struct deve começar com letra maiúscula para ser exportado
// (visível para outros pacotes)
type Usuario struct {
    // 4. Campos do struct
    // Cada campo tem: Nome, Tipo e Tags (opcionais)
    // As tags são metadados entre crases ``
    
    // ============================================================
    // CAMPOS PRINCIPAIS
    // ============================================================
    
    // ID - Chave primária da tabela
    // `gorm:"column:usu_id;primaryKey;autoIncrement"`
    //   - column: nome da coluna no banco de dados
    //   - primaryKey: indica que é chave primária
    //   - autoIncrement: o banco gera automaticamente
    // `json:"id"` - nome usado quando o objeto é convertido para JSON
    ID int `gorm:"column:usu_id;primaryKey;autoIncrement" json:"id"`
    
    // GrupoUsuarioID - Chave estrangeira para grupo_usuario
    // `not null` - campo obrigatório
    // `gpu_id` - nome da coluna no banco
    GrupoUsuarioID int `gorm:"column:gpu_id;not null" json:"grupo_usuario_id"`
    
    // Nome - Nome do usuário
    // `type:varchar(100)` - define o tipo e tamanho no banco
    Nome string `gorm:"column:usu_nome;type:varchar(100);not null" json:"nome"`
    
    // Login - Nome de usuário para autenticação
    // `unique` - não pode haver dois usuários com o mesmo login
    Login string `gorm:"column:usu_login;type:varchar(20);not null;unique" json:"login"`
    
    // Senha - Senha do usuário (será armazenada com hash)
    // `-` na tag json: este campo NÃO aparece nas respostas JSON
    // Isso é importante para não expor a senha
    Senha string `gorm:"column:usu_senha;type:varchar(100);not null" json:"-"`
    
    // Situacao - Status do usuário (1-ativo, 2-inativo, etc)
    // `situiacao` está escrito assim no banco original (com "i")
    Situacao int `gorm:"column:usu_situiacao;not null" json:"situacao"`
    
    // Observacoes - Campo de texto livre
    // `*string` - ponteiro para string (pode ser nil = null no banco)
    // `omitempty` - se for nil, não aparece no JSON
    Observacoes *string `gorm:"column:usu_observacoes;type:text" json:"observacoes,omitempty"`
    
    // SenhaExclusao - Senha para exclusão (segurança extra)
    // `-` no JSON - nunca exposto nas respostas
    SenhaExclusao *string `gorm:"column:usu_senhaexclusao;type:varchar(100)" json:"-"`
    
    // ============================================================
    // CAMPOS DE AUDITORIA (presentes em todas as tabelas)
    // ============================================================
    
    // CreatedAt - Data de criação do registro
    // `default:CURRENT_TIMESTAMP` - valor padrão é a data/hora atual
    CreatedAt time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
    
    // UpdatedAt - Data da última atualização
    // `ON UPDATE CURRENT_TIMESTAMP` - atualiza automaticamente quando o registro é alterado
    UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
    
    // DeletedAt - Data de exclusão lógica (soft delete)
    // `*time.Time` - ponteiro (pode ser nil = não deletado)
    // `index` - cria um índice no banco para consultas mais rápidas
    DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime;index" json:"deleted_at,omitempty"`
    
    // CreatedBy - ID do usuário que criou este registro
    CreatedBy *int `gorm:"column:created_by" json:"created_by,omitempty"`
    
    // UpdatedBy - ID do usuário que atualizou este registro
    UpdatedBy *int `gorm:"column:updated_by" json:"updated_by,omitempty"`
    
    // ============================================================
    // RELACIONAMENTOS (associações com outras tabelas)
    // ============================================================
    
    // GrupoUsuario - Relacionamento com a tabela grupo_usuario
    // `gorm:"foreignKey:gpu_id;references:gpu_id"`
    //   - foreignKey: campo nesta tabela que faz a referência
    //   - references: campo na tabela relacionada
    // `omitempty` - se for nil, não aparece no JSON
    GrupoUsuario *GrupoUsuario `gorm:"foreignKey:gpu_id;references:gpu_id" json:"grupo_usuario,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

// 5. Método TableName
// Define o nome da tabela no banco de dados
// O GORM usa isso para saber qual tabela usar
// Se não definir, ele usa o plural do nome do struct: "usuarios"
func (Usuario) TableName() string {
    return "usuario" // Nome exato da tabela no banco
}

// 6. Método BeforeCreate (Hook do GORM)
// Executado ANTES de criar um novo registro no banco
// Útil para definir valores padrão como CreatedBy
func (u *Usuario) BeforeCreate() error {
    // Se CreatedBy não foi definido, usar 0
    if u.CreatedBy == nil {
        u.CreatedBy = new(int) // Cria um ponteiro para int
        *u.CreatedBy = 0        // Define o valor como 0
    }
    // Se UpdatedBy não foi definido, usar 0
    if u.UpdatedBy == nil {
        u.UpdatedBy = new(int)
        *u.UpdatedBy = 0
    }
    return nil
}

// 7. Método BeforeUpdate (Hook do GORM)
// Executado ANTES de atualizar um registro no banco
func (u *Usuario) BeforeUpdate() error {
    // Se UpdatedBy não foi definido, usar 0
    if u.UpdatedBy == nil {
        u.UpdatedBy = new(int)
        *u.UpdatedBy = 0
    }
    return nil
}

// ============================================================
// MÉTODOS AUXILIARES (funções de conveniência)
// ============================================================

// 8. Método IsActive
// Verifica se o usuário está ativo
// Retorna true se situacao == 1 (ativo)
func (u *Usuario) IsActive() bool {
    return u.Situacao == 1
}

// 9. Método IsDeleted
// Verifica se o usuário foi deletado logicamente
// Retorna true se DeletedAt != nil
func (u *Usuario) IsDeleted() bool {
    return u.DeletedAt != nil
}

// 10. Método SoftDelete
// Realiza a exclusão lógica (soft delete)
// Define a data atual em DeletedAt
func (u *Usuario) SoftDelete() {
    now := time.Now()
    u.DeletedAt = &now
}
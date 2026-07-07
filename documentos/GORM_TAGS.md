# Guia Rápido - Tags GORM

## Tags Mais Comuns

| Tag | Exemplo | Descrição |
|-----|---------|-----------|
| `column` | `column:usu_id` | Nome da coluna no banco |
| `type` | `type:varchar(100)` | Tipo SQL da coluna |
| `not null` | `not null` | Campo obrigatório |
| `default` | `default:1` | Valor padrão |
| `primaryKey` | `primaryKey` | Chave primária |
| `autoIncrement` | `autoIncrement` | Auto incremento |
| `unique` | `unique` | Valor único |
| `index` | `index` | Cria índice |
| `uniqueIndex` | `uniqueIndex` | Índice único |

## Tags para Relacionamentos

| Tag | Exemplo | Descrição |
|-----|---------|-----------|
| `foreignKey` | `foreignKey:gpu_id` | Campo chave estrangeira |
| `references` | `references:gpu_id` | Campo referenciado |
| `many2many` | `many2many:usuario_permissao` | Relação muitos para muitos |

## Tags para JSON

| Tag | Exemplo | Descrição |
|-----|---------|-----------|
| `json:"nome"` | `json:"id"` | Nome no JSON |
| `json:"-"` | `json:"-"` | Esconde o campo no JSON |
| `omitempty` | `json:"nome,omitempty"` | Esconde se for vazio |

## Exemplo Completo

```go
type Usuario struct {
    ID             int       `gorm:"column:usu_id;primaryKey;autoIncrement" json:"id"`
    GrupoUsuarioID int       `gorm:"column:gpu_id;not null" json:"grupo_usuario_id"`
    Nome           string    `gorm:"column:usu_nome;type:varchar(100);not null" json:"nome"`
    Login          string    `gorm:"column:usu_login;type:varchar(20);not null;unique" json:"login"`
    Senha          string    `gorm:"column:usu_senha;type:varchar(100);not null" json:"-"`
    CreatedAt      time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
    DeletedAt      *time.Time `gorm:"column:deleted_at;type:datetime;index" json:"deleted_at,omitempty"`
}
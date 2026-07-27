package repository

import (
	"fmt"

	"github.com/openerp/backend/internal/models"
	"gorm.io/gorm"
)

// ProdutoVariacaoRepository gerencia operações de CRUD para ProdutoVariacao.
type ProdutoVariacaoRepository struct {
	db *gorm.DB
}

// NewProdutoVariacaoRepository cria uma nova instância de ProdutoVariacaoRepository.
func NewProdutoVariacaoRepository(db *gorm.DB) *ProdutoVariacaoRepository {
	return &ProdutoVariacaoRepository{db: db}
}

// Create cria uma nova variação de produto no banco de dados.
func (r *ProdutoVariacaoRepository) Create(variacao *models.ProdutoVariacao) error {
	return r.db.Create(variacao).Error
}

// FindByID busca uma variação de produto pelo ID.
func (r *ProdutoVariacaoRepository) FindByID(id int) (*models.ProdutoVariacao, error) {
	var variacao models.ProdutoVariacao
	err := r.db.
		Preload("Produto").
		Preload("EmpresaFilial").
		Preload("Cor").
		Preload("Tamanho").
		First(&variacao, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("variação de produto com ID %d não encontrada", id)
		}
		return nil, fmt.Errorf("erro ao buscar variação de produto: %w", err)
	}
	return &variacao, nil
}

// Update atualiza uma variação de produto existente no banco de dados.
func (r *ProdutoVariacaoRepository) Update(variacao *models.ProdutoVariacao) error {
	return r.db.Save(variacao).Error
}

// Delete realiza a exclusão lógica de uma variação de produto pelo ID.
func (r *ProdutoVariacaoRepository) Delete(id int) error {
	variacao := &models.ProdutoVariacao{}
	err := r.db.First(variacao, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("variação de produto com ID %d não encontrada para exclusão", id)
		}
		return fmt.Errorf("erro ao buscar variação de produto para exclusão: %w", err)
	}

	variacao.SoftDelete() // Chama o método SoftDelete do modelo
	return r.db.Save(variacao).Error
}

// List lista variações de produto com paginação e filtros.
func (r *ProdutoVariacaoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoVariacao, int64, error) {
	var variacoes []models.ProdutoVariacao
	var total int64
	query := r.db.Model(&models.ProdutoVariacao{}).
		Preload("Produto").
		Preload("EmpresaFilial").
		Preload("Cor").
		Preload("Tamanho")

	// Aplicar filtros
	for key, value := range filters {
		if value != "" {
			switch key {
			case "produto_id":
				query = query.Where("pro_id = ?", value)
			case "empresa_filial_id":
				query = query.Where("emf_id = ?", value)
			case "cor_id":
				query = query.Where("cor_id = ?", value)
			case "tamanho_id":
				query = query.Where("ptam_id = ?", value)
			case "sku":
				query = query.Where("provar_sku LIKE ?", "%"+value.(string)+"%")
			}
		}
	}

	// Contar o total de registros antes de aplicar o limit/offset
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("erro ao contar variações de produto: %w", err)
	}

	// Aplicar paginação
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&variacoes).Error; err != nil {
		return nil, 0, fmt.Errorf("erro ao listar variações de produto: %w", err)
	}

	return variacoes, total, nil
}

// FindBySKU busca uma variação de produto pelo SKU e EmpresaFilialID.
func (r *ProdutoVariacaoRepository) FindBySKU(sku string, empresaFilialID int) (*models.ProdutoVariacao, error) {
	var variacao models.ProdutoVariacao
	err := r.db.Where("provar_sku = ? AND emf_id = ?", sku, empresaFilialID).First(&variacao).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Não encontrado, não é um erro
		}
		return nil, fmt.Errorf("erro ao buscar variação de produto por SKU: %w", err)
	}
	return &variacao, nil
}

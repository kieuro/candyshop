package product

import (
	entity "candyshop/internal/product/entity"
	"candyshop/pkg/response"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type ProductRepository interface {
	GetAllProduct(offset, limit int) ([]entity.Product, *response.Error)
	GetProductByID(id uuid.UUID) (*entity.Product, *response.Error)
	GetProductBySKU(sku string) (*entity.Product, *response.Error)
	CreateProduct(data entity.Product) (*entity.Product, *response.Error)
	UpdateProduct(data entity.Product) *response.Error
	DeleteProduct(id uuid.UUID, deletedAt time.Time) *response.Error
}

type productRepository struct {
	db *sqlx.DB
}

// GetProductBySKU implements ProductRepository.
func (p *productRepository) GetProductBySKU(sku string) (*entity.Product, *response.Error) {
	var product entity.Product

	query := `
		SELECT id, sku, type, name, brand, sugar_level, production_year, distributor, status, created_at
		FROM products
		WHERE sku = $1
	`
	err := p.db.Get(&product, query, sku)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error().Err(err).Int("status", 404).Str("function", "get product by sku").Msg("failed to get product by sku")
			return nil, &response.Error{
				StatusCode: 404,
				Message:    "failed to fetch product",
				Error:      err,
			}
		}

		log.Error().Err(err).Int("status", 500).Str("function", "get product by sku").Msg("failed to get product by sku")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to fetch product",
			Error:      err,
		}
	}

	return &product, nil
}

// CreateProduct implements ProductRepository.
func (p *productRepository) CreateProduct(data entity.Product) (*entity.Product, *response.Error) {
	tx, err := p.db.Beginx()
	if err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "create product").Msg("failed to create product")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to start transaction",
			Error:      err,
		}
	}

	defer tx.Rollback()

	query := `
		INSERT INTO products (id, sku, type, name, brand, sugar_level, production_year, distributor, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, sku, type, name, brand, sugar_level, production_year, distributor, status, created_at
	`

	var model entity.Product

	errInsert := tx.QueryRowx(query,
		data.ID,
		data.SKU,
		data.Type,
		data.Name,
		data.Brand,
		data.SugarLevel,
		data.ProductionYear,
		data.Distributor,
		data.Status).Scan(&model.ID,
		&model.SKU,
		&model.Type,
		&model.Name,
		&model.Brand,
		&model.SugarLevel,
		&model.ProductionYear,
		&model.Distributor,
		&model.Status,
		&model.CreatedAt)

	if errInsert != nil {
		log.Error().Err(errInsert).Int("status", 500).Str("function", "create product").Msg("failed to create product")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to create product",
			Error:      errInsert,
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "create product").Msg("failed to create product")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to commit transaction",
			Error:      err,
		}
	}

	return &model, nil
}

// DeleteProduct implements ProductRepository.
func (p *productRepository) DeleteProduct(id uuid.UUID, deletedAt time.Time) *response.Error {
	tx, err := p.db.Beginx()
	if err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "delete product").Msg("failed to delete product")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to start transaction",
			Error:      err,
		}
	}

	defer tx.Rollback()

	query := `UPDATE products SET deleted_at = $2, status = false WHERE id = $1`

	_, errExec := tx.Exec(query, id, deletedAt)
	if errExec != nil {
		log.Error().Err(errExec).Int("status", 500).Str("function", "delete product").Msg("failed to delete product")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to delete product",
			Error:      errExec,
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "delete product").Msg("failed to delete product")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to commit transaction",
			Error:      err,
		}
	}

	return nil
}

// GetAllProduct implements ProductRepository.
func (p *productRepository) GetAllProduct(offset int, limit int) ([]entity.Product, *response.Error) {
	var products []entity.Product

	query := `
		SELECT id, sku, type, name, brand, sugar_level, production_year, distributor, status, created_at
		FROM products
		LIMIT $1 OFFSET $2
		`

	err := p.db.Select(&products, query, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error().Err(err).Int("status", 404).Str("function", "get all product").Msg("failed to get all product")
			return nil, &response.Error{
				StatusCode: 404,
				Message:    "failed to fetch products",
				Error:      err,
			}
		}

		log.Error().Err(err).Int("status", 500).Str("function", "get all product").Msg("failed to get all product")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to fetch products",
			Error:      err,
		}
	}

	return products, nil
}

// GetProductByID implements ProductRepository.
func (p *productRepository) GetProductByID(id uuid.UUID) (*entity.Product, *response.Error) {
	var product entity.Product

	query := `
		SELECT id, sku, type, name, brand, sugar_level, production_year, distributor, status, created_at
		FROM products
		WHERE id = $1
	`
	err := p.db.Get(&product, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error().Err(err).Int("status", 404).Str("function", "get product by id").Msg("failed to get product by id")
			return nil, &response.Error{
				StatusCode: 404,
				Message:    "failed to fetch product",
				Error:      err,
			}
		}

		log.Error().Err(err).Int("status", 500).Str("function", "get product by id").Msg("failed to get product by id")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to fetch product",
			Error:      err,
		}
	}

	return &product, nil
}

// UpdateProduct implements ProductRepository.
func (p *productRepository) UpdateProduct(data entity.Product) *response.Error {
	tx, err := p.db.Beginx()
	if err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "update product").Msg("failed to update product")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to start transaction",
			Error:      err,
		}
	}

	defer tx.Rollback()

	query := `
		UPDATE products SET sku = $2, type = $3, name = $4, brand = $5, sugar_level = $6, production_year = $7, distributor = $8, updated_at = $9 WHERE id = $1
	`

	_, errExec := tx.Exec(query,
		data.ID,
		data.SKU,
		data.Type,
		data.Name,
		data.Brand,
		data.SugarLevel,
		data.ProductionYear,
		data.Distributor,
		data.UpdatedAt,
	)

	if errExec != nil {
		log.Error().Err(errExec).Int("status", 500).Str("function", "update product").Msg("failed to update product")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to update product",
			Error:      errExec,
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "update product").Msg("failed to update product")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to commit transaction",
			Error:      err,
		}
	}

	return nil
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &productRepository{db}
}

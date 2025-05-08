package product

import (
	dto "candyshop/internal/product/dto"
	entity "candyshop/internal/product/entity"
	repository "candyshop/internal/product/repository"
	"candyshop/pkg/response"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductService interface {
	GetAllProduct(offset, limit int) ([]entity.Product, *response.Error)
	GetProductByID(id uuid.UUID) (*entity.Product, *response.Error)
	CreateProduct(data dto.CreateProductRequest) (*entity.Product, *response.Error)
	UpdateProduct(data dto.UpdateProductRequest) *response.Error
	DeleteProduct(id uuid.UUID) *response.Error
}

type productService struct {
	repository repository.ProductRepository
}

// CreateProduct implements ProductService.
func (p *productService) CreateProduct(data dto.CreateProductRequest) (*entity.Product, *response.Error) {
	checkSKU, errSKU := p.repository.GetProductBySKU(data.SKU)
	if errSKU != nil && errSKU.StatusCode != 404 {
		return nil, errSKU
	}

	if checkSKU != nil {
		return nil, &response.Error{
			StatusCode: fiber.StatusConflict,
			Message:    fmt.Sprintf("sku %s already registered", data.SKU),
			Error:      nil,
		}
	}

	newUUID, _ := uuid.NewV7()

	dataProduct := &entity.Product{
		ID:             newUUID,
		SKU:            data.SKU,
		Type:           data.Type,
		Name:           data.Name,
		Brand:          data.Brand,
		SugarLevel:     data.SugarLevel,
		ProductionYear: data.ProductionYear,
		Distributor:    data.Distributor,
		Status:         true,
	}

	return p.repository.CreateProduct(*dataProduct)
}

// DeleteProduct implements ProductService.
func (p *productService) DeleteProduct(id uuid.UUID) *response.Error {
	checkProduct, errProduct := p.repository.GetProductByID(id)
	if errProduct != nil && errProduct.StatusCode != 404 {
		return errProduct
	}

	// check if product is deleted
	if checkProduct.DeletedAt != nil {
		return &response.Error{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "product is not active",
			Error:      nil,
		}
	}

	currentTime := time.Now()

	return p.repository.DeleteProduct(id, currentTime)
}

// GetAllProduct implements ProductService.
func (p *productService) GetAllProduct(offset int, limit int) ([]entity.Product, *response.Error) {
	return p.repository.GetAllProduct(offset, limit)
}

// GetProductByID implements ProductService.
func (p *productService) GetProductByID(id uuid.UUID) (*entity.Product, *response.Error) {
	return p.repository.GetProductByID(id)
}

// UpdateProduct implements ProductService.
func (p *productService) UpdateProduct(data dto.UpdateProductRequest) *response.Error {
	// check if product is exist
	product, errProduct := p.repository.GetProductByID(data.ID)
	if errProduct != nil {
		return errProduct
	}

	checkSKU, errSKU := p.repository.GetProductBySKU(data.SKU)
	if errSKU != nil && errSKU.StatusCode != 404 {
		return errSKU
	}

	// check if sku already registered
	if checkSKU != nil {
		return &response.Error{
			StatusCode: fiber.StatusConflict,
			Message:    fmt.Sprintf("sku %s already registered", data.SKU),
			Error:      nil,
		}
	}

	if data.SKU == "" {
		data.SKU = product.SKU
	}

	if data.Name == "" {
		data.Name = product.Name
	}

	if data.Type == "" {
		data.Type = product.Type
	}

	if data.Brand == "" {
		data.Brand = product.Brand
	}

	if data.SugarLevel == 0 {
		data.SugarLevel = product.SugarLevel
	}

	if data.ProductionYear == "" {
		data.ProductionYear = product.ProductionYear
	}

	if data.Distributor == "" {
		data.Distributor = product.Distributor
	}

	// assign the value from request to entity
	dataProduct := &entity.Product{
		ID:             data.ID,
		SKU:            data.SKU,
		Type:           data.Type,
		Name:           data.Name,
		Brand:          data.Brand,
		SugarLevel:     data.SugarLevel,
		ProductionYear: data.ProductionYear,
		Distributor:    data.Distributor,
	}

	return p.repository.UpdateProduct(*dataProduct)
}

func NewProductService(repository repository.ProductRepository) ProductService {
	return &productService{repository}
}

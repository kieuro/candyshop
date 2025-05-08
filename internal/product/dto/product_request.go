package product

import "github.com/google/uuid"

type CreateProductRequest struct {
	SKU            string `json:"sku"`
	Type           string `json:"type"`
	Name           string `json:"name"`
	Brand          string `json:"brand"`
	SugarLevel     int    `json:"sugar_level"`
	ProductionYear string `json:"production_year"`
	Distributor    string `json:"distributor"`
}

type UpdateProductRequest struct {
	ID             uuid.UUID `json:"id"`
	SKU            string    `json:"sku"`
	Type           string    `json:"type"`
	Name           string    `json:"name"`
	Brand          string    `json:"brand"`
	SugarLevel     int       `json:"sugar_level"`
	ProductionYear string    `json:"production_year"`
	Distributor    string    `json:"distributor"`
}

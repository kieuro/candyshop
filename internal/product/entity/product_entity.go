package product

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	SKU            string     `json:"sku" db:"sku"`
	Type           string     `json:"type" db:"type"`
	Name           string     `json:"name" db:"name"`
	Brand          string     `json:"brand" db:"brand"`
	SugarLevel     int        `json:"sugar_level" db:"sugar_level"`
	ProductionYear string     `json:"production_year" db:"production_year"`
	Distributor    string     `json:"distributor" db:"distributor"`
	Status         bool       `json:"status" db:"status"`
	CreatedAt      *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time `json:"-" db:"updated_at"`
	DeletedAt      *time.Time `json:"-" db:"deleted_at"`
}

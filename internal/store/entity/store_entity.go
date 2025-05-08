package store

import (
	"time"

	"github.com/google/uuid"
)

type Store struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Address   string     `json:"address" db:"address"`
	Status    bool       `json:"status" db:"status"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"-" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

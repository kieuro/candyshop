package customer

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	PhoneNumber string     `json:"phone_number" db:"phone_number"`
	Address     string     `json:"address" db:"address"`
	Status      bool       `json:"status" db:"status"`
	IsMember    bool       `json:"is_member" db:"is_member"`
	CreatedAt   *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time `json:"-" db:"updated_at"`
	DeletedAt   *time.Time `json:"-" db:"deleted_at"`
}

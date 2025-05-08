package store

import "github.com/google/uuid"

type CreateStoreRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}

type UpdateStoreRequest struct {
	ID      uuid.UUID `json:"id" binding:"required"`
	Name    string    `json:"name" binding:"required"`
	Address string    `json:"address" binding:"required"`
}

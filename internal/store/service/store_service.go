package store

import (
	dto "candyshop/internal/store/dto"
	entity "candyshop/internal/store/entity"
	repository "candyshop/internal/store/repository"
	"candyshop/pkg/response"
	"time"

	"github.com/google/uuid"
)

type StoreService interface {
	GetAllStore(offset, limit int) ([]entity.Store, *response.Error)
	GetStoreByID(id uuid.UUID) (*entity.Store, *response.Error)
	CreateStore(data dto.CreateStoreRequest) (*entity.Store, *response.Error)
	UpdateStore(data dto.UpdateStoreRequest) *response.Error
	DeleteStore(id uuid.UUID) *response.Error
}

type storeService struct {
	repository repository.StoreRepository
}

// CreateStore implements StoreService.
func (p *storeService) CreateStore(data dto.CreateStoreRequest) (*entity.Store, *response.Error) {
	newUUID, _ := uuid.NewV7()

	dataStore := &entity.Store{
		ID:      newUUID,
		Name:    data.Name,
		Address: data.Address,
		Status:  true,
	}

	return p.repository.CreateStore(*dataStore)
}

// DeleteStore implements storeService.
func (p *storeService) DeleteStore(id uuid.UUID) *response.Error {
	checkStore, errStore := p.repository.GetStoreByID(id)
	if errStore != nil {
		return errStore
	}

	if checkStore.DeletedAt != nil {
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to delete store because store already deleted",
			Error:      nil,
		}
	}

	currentTime := time.Now()

	return p.repository.DeleteStore(id, currentTime)
}

// GetAllStore implements storeService.
func (p *storeService) GetAllStore(offset int, limit int) ([]entity.Store, *response.Error) {
	return p.repository.GetAllStore(offset, limit)
}

// GetStoreByID implements storeService.
func (p *storeService) GetStoreByID(id uuid.UUID) (*entity.Store, *response.Error) {
	return p.repository.GetStoreByID(id)
}

// UpdateStore implements storeService.
func (p *storeService) UpdateStore(data dto.UpdateStoreRequest) *response.Error {
	checkStore, errStore := p.repository.GetStoreByID(data.ID)
	if errStore != nil {
		return errStore
	}

	if checkStore.DeletedAt != nil {
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to update store because store is deleted",
			Error:      nil,
		}
	}

	if data.Name == "" {
		data.Name = checkStore.Name
	}

	if data.Address == "" {
		data.Address = checkStore.Address
	}

	currentTime := time.Now()

	dataStore := &entity.Store{
		ID:      checkStore.ID,
		Name:    data.Name,
		Address: data.Address,
		Status:  true,
		UpdatedAt: &currentTime,
	}

	return p.repository.UpdateStore(*dataStore)
}

func NewStoreService(repository repository.StoreRepository) StoreService {
	return &storeService{repository}
}

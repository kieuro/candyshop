package store

import (
	entity "candyshop/internal/store/entity"
	"candyshop/pkg/response"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type StoreRepository interface {
	GetAllStore(offset, limit int) ([]entity.Store, *response.Error)
	GetStoreByID(id uuid.UUID) (*entity.Store, *response.Error)
	CreateStore(data entity.Store) (*entity.Store, *response.Error)
	UpdateStore(data entity.Store) *response.Error
	DeleteStore(id uuid.UUID, deletedAt time.Time) *response.Error
}

type storeRepository struct {
	db *sqlx.DB
}

// CreateStore implements StoreRepository.
func (s *storeRepository) CreateStore(data entity.Store) (*entity.Store, *response.Error) {
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "create store").Msg("failed to create store")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to start transaction",
			Error:      err,
		}
	}

	defer tx.Rollback()

	var model entity.Store

	query := `
		INSERT INTO stores (id, name, address, status) VALUES ($1, $2, $3, $4)
		RETURNING id, name, address, status, created_at
	`

	errInsert := tx.QueryRowx(query,
		data.ID,
		data.Name,
		data.Address,
		data.Status).Scan(&model.ID,
		&model.Name,
		&model.Address,
		&model.Status,
		&model.CreatedAt)

	if errInsert != nil {
		log.Error().Err(errInsert).Int("status", 500).Str("function", "create store").Msg("failed to create store")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to create store",
			Error:      errInsert,
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "create store").Msg("failed to create store")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to commit transaction",
			Error:      err,
		}
	}

	return &model, nil
}

// DeleteStore implements StoreRepository.
func (s *storeRepository) DeleteStore(id uuid.UUID, deletedAt time.Time) *response.Error {
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "delete store").Msg("failed to delete store")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to start transaction",
			Error:      err,
		}
	}

	defer tx.Rollback()

	query := `UPDATE stores SET deleted_at = $2, status = false WHERE id = $1`

	_, errExec := tx.Exec(query, id, deletedAt)
	if errExec != nil {
		log.Error().Err(errExec).Int("status", 500).Str("function", "delete store").Msg("failed to delete store")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to delete store",
			Error:      errExec,
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "delete store").Msg("failed to delete store")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to commit transaction",
			Error:      err,
		}
	}

	return nil
}

// GetAllStore implements StoreRepository.
func (s *storeRepository) GetAllStore(offset int, limit int) ([]entity.Store, *response.Error) {
	var stores []entity.Store

	query := `
		SELECT id, name, address, status, created_at, updated_at, deleted_at 
		FROM stores
		LIMIT $1 OFFSET $2
	`

	err := s.db.Select(&stores, query, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error().Err(err).Int("status", 404).Str("function", "get all store").Msg("failed to get all store")
			return nil, &response.Error{
				StatusCode: 404,
				Message:    "failed to fetch stores",
				Error:      err,
			}
		}

		log.Error().Err(err).Int("status", 500).Str("function", "get all store").Msg("failed to get all store")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to fetch stores",
			Error:      err,
		}
	}

	return stores, nil
}

// GetStoreByID implements StoreRepository.
func (s *storeRepository) GetStoreByID(id uuid.UUID) (*entity.Store, *response.Error) {
	var store entity.Store

	query := `
		SELECT id, name, address, status, created_at, updated_at, deleted_at 
		FROM stores
		WHERE id = $1
	`

	err := s.db.Get(&store, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error().Err(err).Int("status", 404).Str("function", "get store by id").Msg("failed to get store by id")
			return nil, &response.Error{
				StatusCode: 404,
				Message:    "failed to fetch store",
				Error:      err,
			}
		}

		log.Error().Err(err).Int("status", 500).Str("function", "get store by id").Msg("failed to get store by id")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to fetch store",
			Error:      err,
		}
	}

	return &store, nil
}

// UpdateStore implements StoreRepository.
func (s *storeRepository) UpdateStore(data entity.Store) *response.Error {
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "update store").Msg("failed to update store")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to start transaction",
			Error:      err,
		}
	}

	defer tx.Rollback()

	query := `
		UPDATE stores SET name = $2, address = $3, updated_at = $4 WHERE id = $1
	`

	_, errExec := tx.Exec(query, data.ID, data.Name, data.Address, data.UpdatedAt)
	if errExec != nil {
		log.Error().Err(errExec).Int("status", 500).Str("function", "update store").Msg("failed to update store")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to update store",
			Error:      errExec,
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "update store").Msg("failed to update store")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to commit transaction",
			Error:      err,
		}
	}

	return nil
}

func NewStoreRepository(db *sqlx.DB) StoreRepository {
	return &storeRepository{db}
}

package customer

import (
	entity "candyshop/internal/customer/entity"
	"candyshop/pkg/response"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type CustomerRepository interface {
	GetAllCustomer(offset, limit int) ([]entity.Customer, *response.Error)
	GetCustomerByID(id uuid.UUID) (*entity.Customer, *response.Error)
	CreateCustomer(data entity.Customer) (*entity.Customer, *response.Error)
	UpdateCustomer(data entity.Customer) *response.Error
	DeleteCustomer(id uuid.UUID, deletedAt time.Time) *response.Error
}

type customerRepository struct {
	db *sqlx.DB
}

// CreateCustomer implements CustomerRepository.
func (c *customerRepository) CreateCustomer(data entity.Customer) (*entity.Customer, *response.Error) {
	tx, err := c.db.Beginx()
	if err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "create customer").Msg("failed to create customer")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to start transaction",
			Error:      err,
		}
	}

	defer tx.Rollback()

	query := `
		INSERT INTO customers (id, name, phone_number, address, status, is_member) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, phone_number, address, status, is_member, created_at
	`

	var model entity.Customer

	errInsert := tx.QueryRowx(query,
		data.ID,
		data.Name,
		data.PhoneNumber,
		data.Address,
		data.Status,
		data.IsMember).Scan(&model.ID,
		&model.Name,
		&model.PhoneNumber,
		&model.Address,
		&model.Status,
		&model.IsMember,
		&model.CreatedAt)

	if errInsert != nil {
		log.Error().Err(errInsert).Int("status", 500).Str("function", "create customer").Msg("failed to create customer")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to create customer",
			Error:      errInsert,
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "create customer").Msg("failed to create customer")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to commit transaction",
			Error:      err,
		}
	}

	return &model, nil
}

// DeleteCustomer implements CustomerRepository.
func (c *customerRepository) DeleteCustomer(id uuid.UUID, deletedAt time.Time) *response.Error {
	tx, err := c.db.Beginx()
	if err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "delete customer").Msg("failed to delete customer")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to start transaction",
			Error:      err,
		}
	}

	defer tx.Rollback()

	query := `
		UPDATE customers SET deleted_at = $2, status = false WHERE id = $1
	`

	_, errExec := tx.Exec(query, id, deletedAt)
	if errExec != nil {
		log.Error().Err(errExec).Int("status", 500).Str("function", "delete customer").Msg("failed to delete customer")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to delete customer",
			Error:      errExec,
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "delete customer").Msg("failed to delete customer")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to commit transaction",
			Error:      err,
		}
	}

	return nil
}

// GetAllCustomer implements CustomerRepository.
func (c *customerRepository) GetAllCustomer(offset int, limit int) ([]entity.Customer, *response.Error) {
	var customers []entity.Customer

	query := `
		SELECT id, name, phone_number, address, status, created_at, updated_at, deleted_at FROM customers
		LIMIT $1 OFFSET $2
	`

	err := c.db.Select(&customers, query, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error().Err(err).Int("status", 404).Str("function", "get all customer").Msg("failed to get all customer")
			return nil, &response.Error{
				StatusCode: 404,
				Message:    "failed to fetch customers",
				Error:      err,
			}
		}

		log.Error().Err(err).Int("status", 500).Str("function", "get all customer").Msg("failed to get all customer")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to fetch customers",
			Error:      err,
		}
	}

	return customers, nil
}

// GetCustomerByID implements CustomerRepository.
func (c *customerRepository) GetCustomerByID(id uuid.UUID) (*entity.Customer, *response.Error) {
	var customer entity.Customer

	query := `
		SELECT id, name, phone_number, address, status, created_at, updated_at, deleted_at FROM customers
		WHERE id = $1
	`

	err := c.db.Get(&customer, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error().Err(err).Int("status", 404).Str("function", "get customer by id").Msg("failed to get customer by id")
			return nil, &response.Error{
				StatusCode: 404,
				Message:    "failed to fetch customer",
				Error:      err,
			}
		}

		log.Error().Err(err).Int("status", 500).Str("function", "get customer by id").Msg("failed to get customer by id")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to fetch customer",
			Error:      err,
		}
	}

	return &customer, nil
}

// UpdateCustomer implements CustomerRepository.
func (c *customerRepository) UpdateCustomer(data entity.Customer) *response.Error {
	tx, err := c.db.Beginx()
	if err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "update customer").Msg("failed to update customer")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to start transaction",
			Error:      err,
		}
	}

	defer tx.Rollback()

	query := `
		UPDATE customers SET name = $2, phone_number = $3, address = $4, updated_at = $5 WHERE id = $1
	`

	_, errExec := tx.Exec(query,
		data.ID,
		data.Name,
		data.PhoneNumber,
		data.Address,
		data.UpdatedAt)

	if errExec != nil {
		log.Error().Err(errExec).Int("status", 500).Str("function", "update customer").Msg("failed to update customer")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to update customer",
			Error:      errExec,
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "update customer").Msg("failed to update customer")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to commit transaction",
			Error:      err,
		}
	}

	return nil
}

func NewCustomerRepository(db *sqlx.DB) CustomerRepository {
	return &customerRepository{db}
}

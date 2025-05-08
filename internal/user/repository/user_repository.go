package user

import (
	entity "candyshop/internal/user/entity"
	"candyshop/pkg/response"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type UserRepository interface {
	GetAllUser(offset, limit int) ([]entity.User, *response.Error)
	CreateUser(data entity.User) (*entity.User, *response.Error)
	GetUserByID(id uuid.UUID) (*entity.User, *response.Error)
	GetUserByEmail(email string) (*entity.User, *response.Error)
	UpdateUser(data entity.User) *response.Error
	DeleteUser(id uuid.UUID, deletedAt time.Time) *response.Error
}

type userRepository struct {
	db *sqlx.DB
}

// GetUserByID implements UserRepository.
func (u *userRepository) GetUserByID(id uuid.UUID) (*entity.User, *response.Error) {
	var user entity.User

	query := `
		SELECT id, name, email, role, password, status, created_at, updated_at 
		FROM users
		WHERE id = $1
	`

	errData := u.db.Get(&user, query, id)
	if errData != nil {
		if errors.Is(errData, sql.ErrNoRows) {
			log.Error().Err(errData).Int("status", 404).Str("function", "get user by id").Msg("failed to get user by id")
			return nil, &response.Error{
				StatusCode: 404,
				Message:    "failed to fetch user",
				Error:      errData,
			}
		}

		log.Error().Err(errData).Int("status", 500).Str("function", "get user by id").Msg("failed to get user by id")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to fetch user",
			Error:      errData,
		}
	}

	return &user, nil
}

// DeleteUser implements UserRepository.
func (u *userRepository) DeleteUser(id uuid.UUID, deletedAt time.Time) *response.Error {
	tx, err := u.db.Beginx()
	if err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "delete user").Msg("failed to delete user")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to start transaction",
			Error:      err,
		}
	}

	defer tx.Rollback()

	query := `UPDATE users SET deleted_at = $2, status = false WHERE id = $1`

	_, errExec := tx.Exec(query, id, deletedAt)
	if errExec != nil {
		log.Error().Err(errExec).Int("status", 500).Str("function", "delete user").Msg("failed to delete user")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to delete user",
			Error:      errExec,
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "delete user").Msg("failed to delete user")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to commit transaction",
			Error:      err,
		}
	}

	return nil
}

// UpdateUser implements UserRepository.
func (u *userRepository) UpdateUser(data entity.User) *response.Error {
	tx, err := u.db.Beginx()
	if err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "update user").Msg("failed to update user")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to start transaction",
			Error:      err,
		}
	}

	defer tx.Rollback()

	query := `
		UPDATE users SET name = $1, email = $2, password = $3, role = $4, status = $5, updated_at = $6 WHERE id = $7
	`

	_, errExec := tx.Exec(query, data.Name, data.Email, data.Password, data.Role, data.Status, data.UpdatedAt, data.ID)
	if errExec != nil {
		log.Error().Err(errExec).Int("status", 500).Str("function", "update user").Msg("failed to update user")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to update user",
			Error:      errExec,
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "update user").Msg("failed to update user")
		return &response.Error{
			StatusCode: 500,
			Message:    "failed to commit transaction",
			Error:      err,
		}
	}

	return nil
}

// GetUserByEmail implements UserRepository.
func (u *userRepository) GetUserByEmail(email string) (*entity.User, *response.Error) {
	var user entity.User

	query := `
		SELECT id, name, email, role, password, status, created_at, updated_at 
		FROM users
		WHERE email = $1
	`

	errData := u.db.Get(&user, query, email)
	if errData != nil {
		if errors.Is(errData, sql.ErrNoRows) {
			log.Error().Err(errData).Int("status", 404).Str("function", "get user by email").Msg("failed to get user by email")
			return nil, &response.Error{
				StatusCode: 404,
				Message:    "failed to fetch user",
				Error:      errData,
			}
		}

		log.Error().Err(errData).Int("status", 500).Str("function", "get user by email").Msg("failed to get user by email")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to fetch user",
			Error:      errData,
		}
	}

	return &user, nil
}

// CreateUser implements UserRepository.
func (u *userRepository) CreateUser(data entity.User) (*entity.User, *response.Error) {
	tx, err := u.db.Beginx()
	if err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "create user").Msg("failed to create user")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to start transaction",
			Error:      err,
		}
	}

	defer tx.Rollback()

	query := `
		INSERT INTO users (id, name, email, role, status, password) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, email, role, status, created_at
	`

	var model entity.User

	errInsert := tx.QueryRowx(query,
		data.ID,
		data.Name,
		data.Email,
		data.Role,
		data.Status,
		data.Password).Scan(&model.ID,
		&model.Name,
		&model.Email,
		&model.Role,
		&model.Status,
		&model.CreatedAt)

	if errInsert != nil {
		log.Error().Err(errInsert).Int("status", 500).Str("function", "create user").Msg("failed to create user")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to create user",
			Error:      errInsert,
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Int("status", 500).Str("function", "create user").Msg("failed to create user")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to commit transaction",
			Error:      err,
		}
	}

	return &model, nil
}

// GetAllUser implements UserRepository.
func (u *userRepository) GetAllUser(offset int, limit int) ([]entity.User, *response.Error) {
	var users []entity.User

	query := `
		SELECT id, name, email, role, status, created_at, updated_at 
		FROM users
		LIMIT $1 OFFSET $2
	`

	errUser := u.db.Select(&users, query, limit, offset)
	if errUser != nil {
		if errors.Is(errUser, sql.ErrNoRows) {
			log.Error().Err(errUser).Int("status", 404).Str("function", "get all user").Msg("failed to get all user")
			return nil, &response.Error{
				StatusCode: 404,
				Message:    "failed to fetch users",
				Error:      errUser,
			}
		}

		log.Error().Err(errUser).Int("status", 500).Str("function", "get all user").Msg("failed to get all user")
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to fetch users",
			Error:      errUser,
		}
	}

	return users, nil
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db}
}

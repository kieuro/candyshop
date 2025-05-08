package user

import (
	dto "candyshop/internal/user/dto"
	entity "candyshop/internal/user/entity"
	repository "candyshop/internal/user/repository"
	"candyshop/pkg/response"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAllUser(offset, limit int) ([]entity.User, *response.Error)
	CreateUser(data dto.CreateUserRequest) (*entity.User, *response.Error)
	GetUserByEmail(email string) (*entity.User, *response.Error)
	UpdateUser(data dto.UpdateUserRequest) *response.Error
	DeleteUser(id uuid.UUID) *response.Error
}

type userService struct {
	repository repository.UserRepository
}

// DeleteUser implements UserService.
func (u *userService) DeleteUser(id uuid.UUID) *response.Error {
	// check if user is exist
	checkUser, err := u.repository.GetUserByID(id)
	if err != nil {
		return err
	}

	// check if user is active and not deleted
	if checkUser.DeletedAt != nil {
		return &response.Error{
			StatusCode: 500,
			Message:    "user is non active",
			Error:      errors.New("user is non active"),
		}
	}

	currentTime := time.Now()
	return u.repository.DeleteUser(id, currentTime)
}

// UpdateUser implements UserService.
func (u *userService) UpdateUser(data dto.UpdateUserRequest) *response.Error {
	name := data.Name
	email := data.Email
	role := data.Role
	password := data.Password
	status := data.Status

	checkEmail, errEmail := u.repository.GetUserByEmail(data.Email)
	if errEmail != nil && errEmail.StatusCode != 404 {
		return errEmail
	}

	// check email if new email is already registered
	if checkEmail != nil {
		return &response.Error{
			StatusCode: 409,
			Message:    "email already registered",
			Error:      nil,
		}
	}

	checkUser, errUser := u.repository.GetUserByID(data.ID)
	if errUser != nil {
		return errUser
	}

	// check is user is active and not deleted
	if checkUser.DeletedAt != nil {
		return &response.Error{
			StatusCode: 500,
			Message:    "user is non active",
			Error:      nil,
		}
	}

	if name == "" {
		name = checkUser.Name
	}

	if email == "" {
		email = checkUser.Email
	}

	if role == "" {
		role = checkUser.Role
	}

	if password != "" {
		// generate hash for new password if input password is exists
		newHashedPassword, errHashed := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if errHashed != nil {
			return &response.Error{
				StatusCode: 500,
				Message:    "failed to save password",
				Error:      errHashed,
			}
		}

		password = string(newHashedPassword)
	} else {
		password = checkUser.Password
	}

	currentTime := time.Now()

	// assign the request to the entity
	dataUser := &entity.User{
		ID:       data.ID,
		Name:     name,
		Email:    email,
		Password: password,
		Role:     role,
		Status:   status,
		UpdatedAt: &currentTime,
	}

	return u.repository.UpdateUser(*dataUser)
}

// GetUserByEmail implements UserService.
func (u *userService) GetUserByEmail(email string) (*entity.User, *response.Error) {
	return u.repository.GetUserByEmail(email)
}

// CreateUser implements UserService.
func (u *userService) CreateUser(data dto.CreateUserRequest) (*entity.User, *response.Error) {
	checkUser, err := u.repository.GetUserByEmail(data.Email)
	if err != nil && err.StatusCode != 404 {
		return nil, err
	}

	if checkUser != nil {
		return nil, &response.Error{
			StatusCode: 409,
			Message:    "email already registered",
			Error:      nil,
		}
	}

	hashedPassword, errHashed := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if errHashed != nil {
		return nil, &response.Error{
			StatusCode: 500,
			Message:    "failed to save password",
			Error:      errHashed,
		}
	}

	newUUID, _ := uuid.NewV7()

	user := &entity.User{
		ID:       newUUID,
		Name:     data.Name,
		Email:    data.Email,
		Password: string(hashedPassword),
		Role:     data.Role,
		Status:   true,
	}

	dataUser, errUser := u.repository.CreateUser(*user)
	if errUser != nil {
		return nil, errUser
	}

	return dataUser, nil
}

// GetAllUser implements UserService.
func (u *userService) GetAllUser(offset int, limit int) ([]entity.User, *response.Error) {
	return u.repository.GetAllUser(offset, limit)
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{repository}
}

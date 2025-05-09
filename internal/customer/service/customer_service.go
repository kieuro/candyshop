package customer

import (
	dto "candyshop/internal/customer/dto"
	entity "candyshop/internal/customer/entity"
	repository "candyshop/internal/customer/repository"
	"candyshop/pkg/response"
	"time"

	"github.com/google/uuid"
)

type CustomerService interface {
	GetAllCustomer(offset, limit int) ([]entity.Customer, *response.Error)
	GetCustomerByID(id uuid.UUID) (*entity.Customer, *response.Error)
	CreateCustomer(data dto.CreateCustomerRequest) (*entity.Customer, *response.Error)
	UpdateCustomer(data dto.UpdateCustomerRequest) *response.Error
	DeactiveCustomer(id uuid.UUID) *response.Error
}

type customerService struct {
	repository repository.CustomerRepository
}

// CreateCustomer implements CustomerService.
func (c *customerService) CreateCustomer(data dto.CreateCustomerRequest) (*entity.Customer, *response.Error) {
	newUUID, _ := uuid.NewV7()

	dataCustomer := &entity.Customer{
		ID:          newUUID,
		Name:        data.Name,
		PhoneNumber: data.PhoneNumber,
		Address:     data.Address,
		Status:      true,
		IsMember:    false,
	}

	return c.repository.CreateCustomer(*dataCustomer)
}

// DeactiveCustomer implements CustomerService.
func (c *customerService) DeactiveCustomer(id uuid.UUID) *response.Error {
	checkCustomer, errCust := c.repository.GetCustomerByID(id)
	if errCust != nil {
		return errCust
	}

	if checkCustomer.DeletedAt != nil && !checkCustomer.Status {
		return &response.Error{
			StatusCode: 409,
			Message:    "customer already deactived",
			Error:      nil,
		}
	}

	currentTime := time.Now()

	return c.repository.DeleteCustomer(checkCustomer.ID, currentTime)
}

// GetAllCustomer implements CustomerService.
func (c *customerService) GetAllCustomer(offset int, limit int) ([]entity.Customer, *response.Error) {
	return c.repository.GetAllCustomer(offset, limit)
}

// GetCustomerByID implements CustomerService.
func (c *customerService) GetCustomerByID(id uuid.UUID) (*entity.Customer, *response.Error) {
	return c.repository.GetCustomerByID(id)
}

// UpdateCustomer implements CustomerService.
func (c *customerService) UpdateCustomer(data dto.UpdateCustomerRequest) *response.Error {
	checkCustomer, errCust := c.repository.GetCustomerByID(data.ID)
	if errCust != nil {
		return errCust
	}

	if checkCustomer.DeletedAt != nil && !checkCustomer.Status {
		return &response.Error{
			StatusCode: 409,
			Message:    "customer is deactive",
			Error:      nil,
		}
	}

	currentTime := time.Now()

	if data.Name == "" {
		data.Name = checkCustomer.Name
	}

	if data.PhoneNumber == "" {
		data.PhoneNumber = checkCustomer.PhoneNumber
	}

	if data.Address == "" {
		data.Address = checkCustomer.Address
	}

	dataCustomer := &entity.Customer{
		ID:          checkCustomer.ID,
		Name:        data.Name,
		PhoneNumber: data.PhoneNumber,
		Address:     data.Address,
		UpdatedAt:   &currentTime,
	}

	return c.repository.UpdateCustomer(*dataCustomer)
}

func NewCustomerService(repository repository.CustomerRepository) CustomerService {
	return &customerService{repository}
}

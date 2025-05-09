package customer

import (
	dto "candyshop/internal/customer/dto"
	service "candyshop/internal/customer/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CustomerHandler struct {
	service service.CustomerService
}

func NewCustomerHandler(service service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service}
}

func (h *CustomerHandler) GetAllCustomer(c *fiber.Ctx) error {
	offset := c.QueryInt("offset")
	limit := c.QueryInt("limit")

	if offset < 0 || limit < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "offset or limit is invalid",
			"error":       nil,
		})
	}

	customers, errCust := h.service.GetAllCustomer(offset, limit)
	if errCust != nil {
		if errCust.StatusCode == fiber.StatusConflict {
			return c.Status(errCust.StatusCode).JSON(fiber.Map{
				"status_code": errCust.StatusCode,
				"message":     "failed to fetch customers",
				"error":       errCust.Message,
			})
		}

		return c.Status(errCust.StatusCode).JSON(fiber.Map{
			"status_code": errCust.StatusCode,
			"message":     "failed to fetch customers",
			"error":       errCust.Error.Error(),
		})
	}

	if len(customers) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status_code": fiber.StatusNotFound,
			"message":     "failed to fetch customers",
			"error":       "customer not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"message":     "success get data customers",
		"data":        customers,
	})
}

func (h *CustomerHandler) GetCustomerByID(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "id is invalid",
			"error":       nil,
		})
	}

	parseID, _ := uuid.Parse(id)

	customer, errCust := h.service.GetCustomerByID(parseID)
	if errCust != nil {
		return c.Status(errCust.StatusCode).JSON(fiber.Map{
			"status_code": errCust.StatusCode,
			"message":     "failed to fetch customer",
			"error":       errCust.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"message":     "success get data customer",
		"data":        customer,
	})
}

func (h *CustomerHandler) CreateCustomer(c *fiber.Ctx) error {
	var req dto.CreateCustomerRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "failed to input data customer",
			"error":       err.Error(),
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "name is invalid",
			"error":       nil,
		})
	}

	createCustomer, errCust := h.service.CreateCustomer(req)
	if errCust != nil {
		return c.Status(errCust.StatusCode).JSON(fiber.Map{
			"status_code": errCust.StatusCode,
			"message":     "failed to create customer",
			"error":       errCust.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status_code": fiber.StatusCreated,
		"message":     "success create customer",
		"data":        createCustomer,
	})
}

func (h *CustomerHandler) UpdateCustomer(c *fiber.Ctx) error {
	var req dto.UpdateCustomerRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "failed to input data customer",
			"error":       err.Error(),
		})
	}

	errUpdate := h.service.UpdateCustomer(req)
	if errUpdate != nil {
		if errUpdate.StatusCode == fiber.StatusConflict {
			return c.Status(errUpdate.StatusCode).JSON(fiber.Map{
				"status_code": errUpdate.StatusCode,
				"message":     "failed to update customer",
				"error":       errUpdate.Message,
			})
		}

		return c.Status(errUpdate.StatusCode).JSON(fiber.Map{
			"status_code": errUpdate.StatusCode,
			"message":     "failed to update customer",
			"error":       errUpdate.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"message":     "success update data customer",
		"data":        nil,
	})
}

func (h *CustomerHandler) DeactiveCustomer(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "id is invalid",
			"error":       nil,
		})
	}

	parseID, _ := uuid.Parse(id)

	errDeactive := h.service.DeactiveCustomer(parseID)
	if errDeactive != nil {
		if errDeactive.StatusCode == fiber.StatusConflict {
			return c.Status(errDeactive.StatusCode).JSON(fiber.Map{
				"status_code": errDeactive.StatusCode,
				"message":     "failed to deactive customer",
				"error":       errDeactive.Message,
			})
		}

		return c.Status(errDeactive.StatusCode).JSON(fiber.Map{
			"status_code": errDeactive.StatusCode,
			"message":     "failed to deactive customer",
			"error":       errDeactive.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"message":     "success deactive user",
		"data":        nil,
	})
}

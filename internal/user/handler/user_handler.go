package user

import (
	dto "candyshop/internal/user/dto"
	service "candyshop/internal/user/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) GetAllUser(c *fiber.Ctx) error {
	offset := c.QueryInt("offset")
	limit := c.QueryInt("limit")

	dataUsers, errUser := h.service.GetAllUser(offset, limit)
	if errUser != nil {
		return c.Status(errUser.StatusCode).JSON(fiber.Map{
			"status_code": errUser.StatusCode,
			"message":     errUser.Message,
			"error":       errUser.Error.Error(),
		})
	}

	if len(dataUsers) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status_code": fiber.StatusNotFound,
			"message":     "failed to fetch data user",
			"error":       "user not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"message":     "success get all user",
		"data":        dataUsers,
	})
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "failed to input data user",
			"error":       err.Error(),
		})
	}

	dataUser, errUser := h.service.CreateUser(req)
	if errUser != nil {
		if errUser.StatusCode == 409 {
			return c.Status(errUser.StatusCode).JSON(fiber.Map{
				"status_code": errUser.StatusCode,
				"message":     "failed to create user",
				"error":       errUser.Message,
			})
		}

		return c.Status(errUser.StatusCode).JSON(fiber.Map{
			"status_code": errUser.StatusCode,
			"message":     "failed to create user",
			"error":       errUser.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status_code": fiber.StatusCreated,
		"message":     "success create user",
		"data":        dataUser,
	})
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	var req dto.UpdateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "failed to input data user",
			"error":       err.Error(),
		})
	}

	errUpdate := h.service.UpdateUser(req)
	if errUpdate != nil {
		if errUpdate.StatusCode == fiber.StatusConflict {
			return c.Status(errUpdate.StatusCode).JSON(fiber.Map{
				"status_code": errUpdate.StatusCode,
				"message":     errUpdate.Message,
				"error":       nil,
			})
		}

		return c.Status(errUpdate.StatusCode).JSON(fiber.Map{
			"status_code": errUpdate.StatusCode,
			"message":     "failed to update user",
			"error":       errUpdate.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"message":     "success update user",
		"data":        nil,
	})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "id is invalid",
			"error":       nil,
		})
	}

	parseID, errParse := uuid.Parse(id)
	if errParse != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "failed to parsing id",
			"error":       errParse.Error(),
		})
	}

	errDelete := h.service.DeleteUser(parseID)
	if errDelete != nil {
		return c.Status(errDelete.StatusCode).JSON(fiber.Map{
			"status_code": errDelete.StatusCode,
			"message":     errDelete.Message,
			"error":       errDelete.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"message": "success delete user",
		"data": nil,
	})
}
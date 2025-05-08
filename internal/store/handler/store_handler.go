package store

import (
	dto "candyshop/internal/store/dto"
	service "candyshop/internal/store/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StoreHandler struct {
	service service.StoreService
}

func NewStoreHandler(service service.StoreService) *StoreHandler {
	return &StoreHandler{service}
}

func (h *StoreHandler) GetAllStore(c *fiber.Ctx) error {
	offset := c.QueryInt("offset")
	limit := c.QueryInt("limit")

	if offset < 0 || limit < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "offset or limit is invalid",
			"error":       nil,
		})
	}

	stores, errStore := h.service.GetAllStore(offset, limit)
	if errStore != nil {
		return c.Status(errStore.StatusCode).JSON(fiber.Map{
			"status_code": errStore.StatusCode,
			"message":     "failed to fetch stores",
			"error":       errStore.Error.Error(),
		})
	}

	if len(stores) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status_code": fiber.StatusNotFound,
			"message":     "failed to fetch stores",
			"error":       "store not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"message":     "success get data stores",
		"data":        stores,
	})
}

func (h *StoreHandler) CreateStore(c *fiber.Ctx) error {
	var req dto.CreateStoreRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "failed to input data store",
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

	if req.Address == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "address is invalid",
			"error":       nil,
		})
	}

	product, errProduct := h.service.CreateStore(req)
	if errProduct != nil {
		return c.Status(errProduct.StatusCode).JSON(fiber.Map{
			"status_code": errProduct.StatusCode,
			"message":     "failed to create store",
			"error":       errProduct.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status_code": fiber.StatusCreated,
		"message":     "success create store",
		"data":        product,
	})
}

func (h *StoreHandler) GetStoreByID(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "id is invalid",
			"error":       nil,
		})
	}

	parseID, _ := uuid.Parse(id)

	store, errStore := h.service.GetStoreByID(parseID)
	if errStore != nil {
		return c.Status(errStore.StatusCode).JSON(fiber.Map{
			"status_code": errStore.StatusCode,
			"message":     "failed to fetch store",
			"error":       errStore.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"message":     "success get store",
		"data":        store,
	})
}

func (h *StoreHandler) UpdateStore(c *fiber.Ctx) error {
	var req dto.UpdateStoreRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "failed to input data store",
			"error":       err.Error(),
		})
	}

	errUpdate := h.service.UpdateStore(req)
	if errUpdate != nil {
		return c.Status(errUpdate.StatusCode).JSON(fiber.Map{
			"status_code": errUpdate.StatusCode,
			"message":     "failed to update store",
			"error":       errUpdate.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"message":     "success update data store",
		"data":        nil,
	})
}

func (h *StoreHandler) DeleteStore(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "id is invalid",
			"error":       nil,
		})
	}

	parseID, _ := uuid.Parse(id)

	errDelete := h.service.DeleteStore(parseID)
	if errDelete != nil {
		return c.Status(errDelete.StatusCode).JSON(fiber.Map{
			"status_code": errDelete.StatusCode,
			"message":     "failed to delete store",
			"error":       errDelete.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"message":     "success delete data store",
		"data":        nil,
	})
}

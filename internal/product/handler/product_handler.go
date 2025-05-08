package product

import (
	dto "candyshop/internal/product/dto"
	service "candyshop/internal/product/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service}
}

func (h *ProductHandler) GetAllProduct(c *fiber.Ctx) error {
	offset := c.QueryInt("offset")
	limit := c.QueryInt("limit")

	if offset < 0 || limit < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "invalid offset or limit",
			"error":       nil,
		})
	}

	products, errProduct := h.service.GetAllProduct(offset, limit)
	if errProduct != nil {
		return c.Status(errProduct.StatusCode).JSON(fiber.Map{
			"status_code": errProduct.StatusCode,
			"message":     "failed to fetch data products",
			"error":       errProduct.Error.Error(),
		})
	}

	// check if product is exists
	if len(products) == 0 || products == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status_code": fiber.StatusNotFound,
			"message":     "failed to fetch data products",
			"error":       "product not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"message":     "success get data all product",
		"data":        products,
	})
}

func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "id is invalid",
			"error":       nil,
		})
	}

	// parse param id to uuid format
	parseID, _ := uuid.Parse(id)

	product, errProduct := h.service.GetProductByID(parseID)
	if errProduct != nil {
		return c.Status(errProduct.StatusCode).JSON(fiber.Map{
			"status_code": errProduct.StatusCode,
			"message":     "failed to fetch product",
			"error":       errProduct.Error.Error(),
		})
	}

	// check if product is exists
	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status_code": fiber.StatusNotFound,
			"message":     "failed to fetch product",
			"error":       "product not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"message":     "success get data product",
		"data":        product,
	})
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req dto.CreateProductRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "failed to input data product",
			"error":       err.Error(),
		})
	}

	product, errProduct := h.service.CreateProduct(req)
	if errProduct != nil {
		if errProduct.StatusCode == fiber.StatusConflict {
			return c.Status(errProduct.StatusCode).JSON(fiber.Map{
				"status_code": errProduct.StatusCode,
				"message":     "failed to create product",
				"error":       errProduct.Message,
			})
		}

		return c.Status(errProduct.StatusCode).JSON(fiber.Map{
			"status_code": errProduct.StatusCode,
			"message":     "failed to create product",
			"error":       errProduct.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status_code": fiber.StatusCreated,
		"message":     "success create data product",
		"data":        product,
	})
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	var req dto.UpdateProductRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "failed to input data product",
			"error":       err.Error(),
		})
	}

	errUpdateProduct := h.service.UpdateProduct(req)
	if errUpdateProduct != nil {
		if errUpdateProduct.StatusCode == fiber.StatusConflict {
			return c.Status(errUpdateProduct.StatusCode).JSON(fiber.Map{
				"status_code": errUpdateProduct.StatusCode,
				"message":     "failed to update product",
				"error":       errUpdateProduct.Message,
			})
		}

		return c.Status(errUpdateProduct.StatusCode).JSON(fiber.Map{
			"status_code": errUpdateProduct.StatusCode,
			"message":     "failed to update product",
			"error":       errUpdateProduct.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"message":     "success update data product",
		"data":        nil,
	})
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"message":     "id is invalid",
			"error":       nil,
		})
	}

	// parse param id to uuid format
	parseID, _ := uuid.Parse(id)

	errDeleteProduct := h.service.DeleteProduct(parseID)
	if errDeleteProduct != nil {
		return c.Status(errDeleteProduct.StatusCode).JSON(fiber.Map{
			"status_code": errDeleteProduct.StatusCode,
			"message":     "failed to delete product",
			"error":       errDeleteProduct.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"message":     "success delete data product",
		"data":        nil,
	})
}

package product

import (
	handler "candyshop/internal/product/handler"
	repository "candyshop/internal/product/repository"
	service "candyshop/internal/product/service"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func Init(router fiber.Router, db *sqlx.DB) {
	repo := repository.NewProductRepository(db)
	service := service.NewProductService(repo)
	handler := handler.NewProductHandler(service)

	productRoute := router.Group("api/v1/products")

	productRoute.Get("", handler.GetAllProduct)
	productRoute.Post("", handler.CreateProduct)
	productRoute.Get("/:id", handler.GetProductByID)
	productRoute.Patch("", handler.UpdateProduct)
	productRoute.Patch("/delete/:id", handler.DeleteProduct)
}

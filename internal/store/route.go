package store

import (
	handler "candyshop/internal/store/handler"
	repository "candyshop/internal/store/repository"
	service "candyshop/internal/store/service"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func Init(router fiber.Router, db *sqlx.DB) {
	repo := repository.NewStoreRepository(db)
	service := service.NewStoreService(repo)
	handler := handler.NewStoreHandler(service)

	storeRoute := router.Group("api/v1/stores")

	storeRoute.Get("", handler.GetAllStore)
	storeRoute.Get(":id", handler.GetStoreByID)
	storeRoute.Post("", handler.CreateStore)
	storeRoute.Patch("", handler.UpdateStore)
	storeRoute.Patch("/delete/:id", handler.DeleteStore)
}

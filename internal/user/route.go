package user

import (
	handler "candyshop/internal/user/handler"
	repository "candyshop/internal/user/repository"
	service "candyshop/internal/user/service"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func Init(router fiber.Router, db *sqlx.DB) {
	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo)
	handler := handler.NewUserHandler(service)

	userRoute := router.Group("api/v1/users")

	userRoute.Get("", handler.GetAllUser)
	userRoute.Post("", handler.CreateUser)
	userRoute.Patch("", handler.UpdateUser)
	userRoute.Patch("/delete/:id", handler.DeleteUser)
}

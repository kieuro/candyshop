package customer

import (
	handler "candyshop/internal/customer/handler"
	repository "candyshop/internal/customer/repository"
	service "candyshop/internal/customer/service"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func Init(router fiber.Router, db *sqlx.DB) {
	repo := repository.NewCustomerRepository(db)
	service := service.NewCustomerService(repo)
	handler := handler.NewCustomerHandler(service)

	customerRoute := router.Group("api/v1/customers")

	customerRoute.Get("", handler.GetAllCustomer)
	customerRoute.Get(":id", handler.GetCustomerByID)
	customerRoute.Post("", handler.CreateCustomer)
	customerRoute.Patch("", handler.UpdateCustomer)
	customerRoute.Patch("deactive/:id", handler.DeactiveCustomer)
}

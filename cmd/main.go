package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	r := fiber.New()

	r.Get("", func (c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status_code": fiber.StatusOK,
			"message": "Hello World",
		})
	})

	r.Listen(":5000")
}
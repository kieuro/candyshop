package main

import (
	product "candyshop/internal/product"
	store "candyshop/internal/store"
	user "candyshop/internal/user"
	"candyshop/pkg/db"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func initZeroLogger() {
	// Buka atau buat file `app.log` untuk menulis log
	logFile, err := os.OpenFile("./logs/app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("Error opening log file")
	}

	// Konfigurasikan `zerolog` untuk menulis ke file `app.log` dan `stdout`
	multi := zerolog.MultiLevelWriter(os.Stdout, logFile)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	// Set level log berdasarkan environment
	if os.Getenv("ENV") == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func main() {
	initZeroLogger()
	r := fiber.New()
	r.Use(loggerMiddleware)
	db := db.ConnectDBCandyShop()

	r.Get("", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status_code": fiber.StatusOK,
			"message":     "Hello World",
		})
	})

	user.Init(r, db)
	product.Init(r, db)
	store.Init(r, db)

	r.Listen(":5000")
}

func loggerMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	requestId, _ := uuid.NewV7()
	log.Info().
		Str("request_id", requestId.String()).
		Str("method", c.Method()).
		Str("path", c.Path()).
		Str("ip_address", c.IP()).
		Str("hostname", c.Hostname()).
		Int("status", c.Response().StatusCode()).
		Dur("latency", time.Since(start)).
		Msg("HTTP request")
	return err
}

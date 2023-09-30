package middleware

import (
	"myfinace/database"
	"myfinace/service"

	"github.com/gofiber/fiber/v2"
)

func PortfolioMiddleware(c *fiber.Ctx) error {
	if db, ok := c.Locals("DB").(database.DB); ok {
		c.Locals("Service", service.NewPortfolioService(db))
	} else {
		c.Locals("Service", service.NewPortfolioService(database.GetDB()))
	}
	return c.Next()
}

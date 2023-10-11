package middleware

import (
	"myfinance/database"
	"myfinance/service"

	"github.com/gofiber/fiber/v2"
)

func PortfolioMiddleware(c *fiber.Ctx) error {
	var db database.DB
	if is_db, ok := c.Locals("DB").(database.DB); ok {
		db = is_db
	} else {
		db = database.GetDB()
	}
	c.Locals("Service", service.NewPortfolioService(db))
	c.Locals("HoldingService", service.NewHoldingService(db))
	return c.Next()
}

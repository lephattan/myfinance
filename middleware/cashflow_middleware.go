package middleware

import (
	"myfinance/database"
	"myfinance/service"

	"github.com/gofiber/fiber/v2"
)

func CashflowMiddleware(c *fiber.Ctx) error {
	if db, ok := c.Locals("DB").(database.DB); ok {
		c.Locals("Service", service.NewCasflowService(db))
	} else {
		c.Locals("Service", service.NewCasflowService(database.GetDB()))
	}
	return c.Next()
}

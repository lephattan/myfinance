package htmx

import (
	"myfinance/middleware"
	"myfinance/model"
	"myfinance/service"

	"github.com/gofiber/fiber/v2"
)

func RegisterCashflowComponentController(router fiber.Router) {
	router.Use(middleware.CashflowMiddleware)

	router.Get("chart", HandleAllCashflow)

}

func HandleAllCashflow(c *fiber.Ctx) error {
	svc, _ := c.Locals("Service").(service.CashflowService)
	var cashflow model.Cashflow
	if err := svc.List(c.Context(), &cashflow.Days); err != nil {
		return err
	}
	data := fiber.Map{
		"Cashflow": &cashflow,
	}
	return c.Render("parts/cashflow/cashflow-chart", data)

}

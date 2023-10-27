package htmx

import (
	"myfinance/middleware"
	"myfinance/model"
	"myfinance/service"

	"github.com/gofiber/fiber/v2"
)

func RegisterCashflowComponentController(router fiber.Router) {
	router.Use(middleware.CashflowMiddleware)

	router.Get("chart", HandleCashflowChart)

}

func HandleCashflowChart(c *fiber.Ctx) error {
	svc, _ := c.Locals("Service").(service.CashflowService)

	var req model.CashflowListingOptions
	c.QueryParser(&req)

	var cashflow model.Cashflow
	if err := svc.List(c.Context(), req, &cashflow.Days); err != nil {
		return err
	}
	data := fiber.Map{
		"Cashflow": &cashflow,
	}
	return c.Render("parts/cashflow/cashflow-chart", data)

}

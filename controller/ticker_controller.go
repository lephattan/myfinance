package controller

import (
	"fmt"
	"myfinance/controller/htmx"
	"myfinance/database"
	"myfinance/model"
	"myfinance/service"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Register handlers to prefix "/ticker"
func RegisterTickerController(router fiber.Router) {
	router.Use(TickerMiddleware)
	router.Get("/", TickerListHanlde)
	router.Post("/", HanldeTickerCreate)
	router.Get("/:symbol", TickerHanlde)
	router.Put("/:symbol", HandleTickerUpdate)
}

func TickerMiddleware(c *fiber.Ctx) error {
	if db, ok := c.Locals("DB").(database.DB); ok {
		c.Locals("Service", service.NewTickerService(db))
	} else {
		c.Locals("Service", service.NewTickerService(database.GetDB()))
	}
	return c.Next()
}

// Handle ticker list request
func TickerListHanlde(c *fiber.Ctx) error {
	queryString := string(c.Request().URI().QueryString())
	data := fiber.Map{
		"Title":       "Tickers",
		"QueryString": queryString,
	}
	return c.Render("ticker/tickers", data, "layouts/main")
}

// Handle ticker detail request
func TickerHanlde(c *fiber.Ctx) error {
	symbol := c.Params("symbol")
	data := fiber.Map{
		"Title":  strings.ToUpper(symbol),
		"Symbol": symbol,
	}
	return c.Render("ticker/detail", data, "layouts/main")
}

// Handle ticker update request
func HandleTickerUpdate(c *fiber.Ctx) error {
	symbol := strings.TrimSpace(c.Params("symbol"))
	ticker := model.Ticker{
		Symbol: symbol,
		Name:   c.FormValue("ticker-name"),
	}
	svc, _ := c.Locals("Service").(service.TickerService)

	_, err := svc.Update(c.Context(), ticker)
	if err != nil {
		return c.SendString(fmt.Sprintf("<h3>%s</h3>", err.Error()))
	}
	return c.Redirect(fmt.Sprintf("/htmx/components/ticker/detail/%s", symbol), fiber.StatusSeeOther)
}

// Hanlde ticker create request
func HanldeTickerCreate(c *fiber.Ctx) error {
	symbol := c.FormValue("ticker-symbol")
	name := c.FormValue("ticker-name", "")
	ticker := model.Ticker{
		Symbol: symbol,
		Name:   name,
	}
	db := database.GetDB()
	service := service.NewTickerService(db)
	_, err := service.Create(c.Context(), ticker)
	if err != nil {
		return c.SendString(fmt.Sprintf("<h3>%s</h3>", err.Error()))
	}
	c.Set("HX-Trigger", "new-ticker")
	return htmx.HandleTickerAddForm(c)
}

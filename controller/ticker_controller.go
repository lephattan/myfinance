package controller

import (
	"fmt"
	"myfinace/controller/htmx"
	"myfinace/database"
	"myfinace/model"
	"myfinace/service"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kataras/iris/v12"
)

// Register handlers to prefix "/ticker"
func RegisterTickerController(router fiber.Router) {
	router.Get("/", TickerListHanlde)
	router.Post("/", HanldeTickerCreate)
	router.Get("/:symbol", TickerHanlde)
	router.Put("/:symbol", HandleTickerUpdate)
}

type TickerController struct {
	Service service.TickerService
	Ctx     iris.Context
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

	db := database.GetDB()
	service := service.NewTickerService(db)
	_, err := service.Update(c.Context(), ticker)
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

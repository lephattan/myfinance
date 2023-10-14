package htmx

import (
	"context"
	"log"
	"myfinance/database"
	"myfinance/model"
	"myfinance/service"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

// Register handlers to prefix "/htmx/components/ticker"
func RegisterTickerComponentController(router fiber.Router) {
	router.Use(TickerMiddleware)
	router.Get("/list", HandleTickerList)
	router.Get("/detail/:symbol", HandleTickerDetail)
	router.Get("/edit-form/:symbol", HandleTickerEditForm)
	router.Get("/add-form", HandleTickerAddForm).Name("CTickerAddForm")
}

func TickerMiddleware(c *fiber.Ctx) error {
	if db, ok := c.Locals("DB").(database.DB); ok {
		c.Locals("Service", service.NewTickerService(db))
	} else {
		c.Locals("Service", service.NewTickerService(database.GetDB()))
	}
	return c.Next()
}

// Get ticker by its symbol
func GetTicker(symbol string, db database.DB, ctx context.Context) (model.Ticker, error) {
	service := service.NewTickerService(db)
	var ticker model.Ticker
	err := service.Get(ctx, symbol, &ticker)
	return ticker, err
}

func HandleTickerList(c *fiber.Ctx) error {
	errors := []string{}
	queryString := string(c.Request().URI().QueryString())

	var tickers model.Tickers
	svc, _ := c.Locals("Service").(service.TickerService)
	url, err := url.ParseRequestURI(c.OriginalURL())
	if err != nil {
		return err
	}
	urlValues := url.Query()
	opt := tickers.ParseListOptions(&urlValues)
	if err := svc.List(c.Context(), opt, &tickers); err != nil {
		errors = append(errors, err.Error())
	}
	data := fiber.Map{
		"Tickers":     tickers,
		"Errors":      errors,
		"QueryString": queryString,
	}
	return c.Render("parts/ticker/list", data)
}

func HandleTickerDetail(c *fiber.Ctx) error {
	symbol := c.Params("symbol")
	svc, _ := c.Locals("Service").(service.TickerService)
	var ticker model.Ticker
	errors := []string{}
	err := svc.Get(c.Context(), symbol, &ticker)
	if err != nil {
		errors = append(errors, err.Error())
	}
	data := fiber.Map{
		"Ticker": &ticker,
		"Errors": &errors,
	}
	return c.Render("parts/ticker/detail", data)
}

func HandleTickerEditForm(c *fiber.Ctx) error {
	symbol := c.Params("symbol")
	log.Printf("Edit form for %s", symbol)
	svc, _ := c.Locals("Service").(service.TickerService)
	var ticker model.Ticker
	err := svc.Get(c.Context(), symbol, &ticker)
	if err != nil {
		return err
	}
	data := fiber.Map{
		"Ticker": &ticker,
	}
	return c.Render("parts/ticker/edit-form", data)
}

func HandleTickerAddForm(c *fiber.Ctx) error {
	return c.Render("parts/ticker/add-form", fiber.Map{})
}

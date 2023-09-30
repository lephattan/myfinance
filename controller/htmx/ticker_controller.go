package htmx

import (
	"context"
	"fmt"
	"log"
	"myfinace/database"
	"myfinace/env"
	"myfinace/model"
	"myfinace/service"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/kataras/iris/v12"
)

type HTMXTickerController struct {
	Service service.TickerService
	Ctx     iris.Context
}

// Register handlers to prefix "/htmx/components/ticker"
func RegisterTickerComponentController(router fiber.Router) {
	router.Get("/list", HandleTickerList)
	router.Get("/detail/:symbol", HandleTickerDetail)
	router.Get("/edit-form/:symbol", HandleTickerEditForm)
	router.Get("/add-form", HandleTickerAddForm).Name("CTickerAddForm")
}

// Get ticker by its symbol
func GetTicker(symbol string, db database.DB, ctx context.Context) (model.Ticker, error) {
	service := service.NewTickerService(db)
	var ticker model.Ticker
	err := service.Get(ctx, symbol, &ticker)
	return ticker, err
}

func HandleTickerList(c *fiber.Ctx) error {
	app_env := env.ReadEnv("APP_ENV", "production")
	db := database.NewDB(app_env)
	errors := []string{}
	queryString := string(c.Request().URI().QueryString())

	var tickers model.Tickers
	service := service.NewTickerService(db)
	url, err := url.ParseRequestURI(c.OriginalURL())
	if err != nil {
		return err
	}
	urlValues := url.Query()
	opt := tickers.ParseListOptions(&urlValues)
	if err := service.List(c.Context(), opt, &tickers); err != nil {
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
	app_env := env.ReadEnv("APP_ENV", "production")
	db := database.NewDB(app_env)
	service := service.NewTickerService(db)
	var ticker model.Ticker
	errors := []string{}
	err := service.Get(c.Context(), symbol, &ticker)
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
	db := database.GetDB()
	ticker, err := GetTicker(symbol, db, c.Context())
	if err != nil {
		return c.SendString(fmt.Sprintf("<h3>%s</h3>", err.Error()))
	}
	data := fiber.Map{
		"Ticker": &ticker,
	}
	return c.Render("parts/ticker/edit-form", data)
}

func HandleTickerAddForm(c *fiber.Ctx) error {
	return c.Render("parts/ticker/add-form", fiber.Map{})
}

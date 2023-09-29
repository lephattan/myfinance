package htmx

import (
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

func RegisterTickerComponentController(router fiber.Router) {
	router.Get("/list", HandleTickerList)
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

func (c *HTMXTickerController) GetList() {
	errors := []string{}
	var tickers model.Tickers
	ctx := c.Ctx.Request().Context()
	urlValues := c.Ctx.Request().URL.Query()
	opt := tickers.ParseListOptions(&urlValues)

	if err := c.Service.List(ctx, opt, &tickers); err != nil {
		errors = append(errors, err.Error())
	}
	data := iris.Map{
		"Tickers": tickers,
		"Errors":  errors,
	}
	if err := c.Ctx.View("parts/ticker/list", data); err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}

}

func (c *HTMXTickerController) GetAddnewform() {
	data := iris.Map{
		"ID": "add-ticker",
	}
	if err := c.Ctx.View("parts/ticker/add-new-form", data); err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

func (c *HTMXTickerController) PostAddnewform() {
	errors := []string{}
	ctx := c.Ctx
	symbol := ctx.FormValue("ticker-symbol")
	name := ctx.FormValueDefault("ticker-name", "")
	ticker := model.Ticker{
		Symbol: symbol,
		Name:   name,
	}

	_, err := c.Service.Create(c.Ctx.Request().Context(), ticker)
	if err != nil {
		errors = append(errors, err.Error())
	} else {
		ctx.Header("HX-Trigger", "new-ticker")
	}
	data := iris.Map{
		"Errors": errors,
	}

	if err := c.Ctx.View("parts/ticker/add-new-form", data); err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

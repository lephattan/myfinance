package controller

import (
	"fmt"
	"myfinace/model"
	"myfinace/service"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kataras/iris/v12"
)

func RegisterTickerController(router fiber.Router) {
	router.Get("/", TickerHanlde)
}

type TickerController struct {
	Service service.TickerService
	Ctx     iris.Context
}

func TickerHanlde(c *fiber.Ctx) error {
	data := fiber.Map{
		"Title": "Tickers",
	}
	return c.Render("ticker/tickers", data, "layouts/main")

}

func (c *TickerController) Get() {
	c.Ctx.ViewLayout("main")
	data := iris.Map{
		"Title": "Tickers",
	}
	if err := c.Ctx.View("ticker/tickers", data); err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

func (c *TickerController) GetBy(symbol string) {
	c.Ctx.ViewLayout("main")
	var ticker model.Ticker
	errors := []string{}
	err := c.Service.Get(c.Ctx.Request().Context(), symbol, &ticker)
	if err != nil {
		errors = append(errors, err.Error())
	}
	data := iris.Map{
		"Title":  strings.ToUpper(symbol),
		"Ticker": &ticker,
		"Errors": &errors,
	}
	if err := c.Ctx.View("ticker/detail", data); err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

func (c *TickerController) PostBy(symbol string) {
	symbol = strings.TrimSpace(symbol)
	ticker := model.Ticker{
		Symbol: symbol,
		Name:   c.Ctx.FormValue("ticker-name"),
	}
	_, err := c.Service.Update(c.Ctx.Request().Context(), ticker)
	if err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
	c.Ctx.Redirect(fmt.Sprintf("/ticker/%s", symbol), iris.StatusSeeOther)
}

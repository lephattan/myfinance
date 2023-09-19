package controller

import (
	"fmt"
	"myfinace/model"
	"myfinace/service"

	"github.com/kataras/iris/v12"
)

type TickerRequest struct {
}

type TickerController struct {
	Service service.TickerService
	Ctx     iris.Context
}

func (c *TickerController) Get() {
	c.Ctx.ViewLayout("main")
	var tickers model.Tickers
	ctx := c.Ctx.Request().Context()
	c.Service.List(ctx, &tickers)
	data := iris.Map{
		"Tickers": tickers,
	}
	if err := c.Ctx.View("ticker/tickers", data); err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

func (c *TickerController) Post() {
	symbol := c.Ctx.FormValue("ticker-symbol")
	name := c.Ctx.FormValueDefault("ticker-name", "")
	ticker := model.Ticker{
		Symbol: symbol,
		Name:   name,
	}
	c.Service.Create(c.Ctx.Request().Context(), ticker)
	c.Get()
}

func (c *TickerController) GetBy(symbol string) {
	c.Ctx.ViewLayout("main")
	var ticker model.Ticker
	errors := []string{}
	err := c.Service.Get(c.Ctx.Request().Context(), symbol, &ticker)
	if err != nil {
		errors = append(errors, err.Error())
	}
	fmt.Println(errors)
	data := iris.Map{
		"Ticker": &ticker,
		"Errors": &errors,
	}
	if err := c.Ctx.View("ticker/detail", data); err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

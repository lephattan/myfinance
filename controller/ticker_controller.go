package controller

import (
	"fmt"
	"myfinace/model"
	"myfinace/service"
	"strings"

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
		"Title":   "Tickers",
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
	_, err := c.Service.Create(c.Ctx.Request().Context(), ticker)
	if err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
	c.Ctx.Redirect("/ticker", iris.StatusSeeOther)
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

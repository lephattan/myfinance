package controller

import (
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
	/* tickers := []model.Ticker{
		{
			Symbol: "nvl",
			Name:   "Novaland",
		},
		{

			Symbol: "msb",
		},
	} */
	c.Ctx.ViewLayout("main")
	var tickers model.Tickers
	ctx := c.Ctx.Request().Context()
	c.Service.List(ctx, &tickers)
	data := iris.Map{
		"Tickers": tickers,
	}
	if err := c.Ctx.View("tickers", data); err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

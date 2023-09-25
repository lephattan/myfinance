package htmx

import (
	"myfinace/model"
	"myfinace/service"

	"github.com/kataras/iris/v12"
)

type HTMXTickerController struct {
	Service service.TickerService
	Ctx     iris.Context
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

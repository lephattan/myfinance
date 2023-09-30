package controller

import (
	"fmt"
	"log"
	"myfinace/model"
	"myfinace/service"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kataras/iris/v12"
)

func RegisterTickerController(router fiber.Router) {
	router.Get("/", TickerListHanlde)
	router.Get("/:symbol", TickerHanlde)
}

type TickerController struct {
	Service service.TickerService
	Ctx     iris.Context
}

func TickerListHanlde(c *fiber.Ctx) error {
	queryString := string(c.Request().URI().QueryString())
	data := fiber.Map{
		"Title":       "Tickers",
		"QueryString": queryString,
	}
	return c.Render("ticker/tickers", data, "layouts/main")
}

func TickerHanlde(c *fiber.Ctx) error {
	log.Print("Ticker detail handling")
	symbol := c.Params("symbol")
	data := fiber.Map{
		"Title":  strings.ToUpper(symbol),
		"Symbol": symbol,
	}
	return c.Render("ticker/detail", data, "layouts/main")
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

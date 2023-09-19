package controller

import (
	"myfinace/model"
	"myfinace/service"
	"strings"

	"github.com/kataras/iris/v12"
)

type PortfolioController struct {
	Service service.PortfolioService
	Ctx     iris.Context
	Errors  []string
}

func (c *PortfolioController) Error(err string) {
	c.Errors = append(c.Errors, err)
}

func (c *PortfolioController) Get() {
	c.Ctx.ViewLayout("main")
	var portfolios model.Portfolios
	ctx := c.Ctx.Request().Context()
	c.Service.List(ctx, &portfolios)
	data := iris.Map{
		"Title":      "Portfolios",
		"Portfolios": portfolios,
		"Errors":     c.Errors,
	}
	if err := c.Ctx.View("portfolio/portfolios", data); err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

func (c *PortfolioController) Post() {
	name := c.Ctx.FormValue("portfolio-name")
	des := c.Ctx.FormValueDefault("portfolio-des", "")
	portfolio := model.Portfolio{
		Name:        name,
		Description: des,
	}
	c.Service.Create(c.Ctx.Request().Context(), portfolio)
	c.Get()
}

func (c *PortfolioController) GetBy(id uint64) {
	c.Ctx.ViewLayout("main")
	var portfolio model.Portfolio
	err := c.Service.Get(c.Ctx.Request().Context(), id, &portfolio)
	if err != nil {
		c.Error(err.Error())
	}
	data := iris.Map{
		"Title":     strings.ToUpper(portfolio.Name),
		"Portfolio": &portfolio,
		"Errors":    c.Errors,
	}
	if err := c.Ctx.View("portfolio/detail", data); err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

func (c *PortfolioController) PostBy(id uint64) {
	portfolio := model.Portfolio{
		ID: id,
	}
	if name := c.Ctx.FormValue("portfolio-name"); len(name) > 0 {
		portfolio.Name = name
	}
	if des := c.Ctx.FormValue("portfolio-des"); len(des) > 0 {
		portfolio.Description = des
	}
	_, err := c.Service.Update(c.Ctx.Request().Context(), portfolio)
	if err != nil {
		c.Ctx.ViewLayout("main")
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
	c.GetBy(id)
}

func (c *PortfolioController) PostDeleteBy(id uint64) {
	_, err := c.Service.Delete(c.Ctx.Request().Context(), id)
	if err != nil {
		c.Error(err.Error())
		c.GetBy(id)
		return
	}
	c.Ctx.Redirect("/portfolio")
}

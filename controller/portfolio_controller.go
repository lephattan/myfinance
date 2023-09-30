package controller

import (
	"fmt"
	"myfinace/middleware"
	"myfinace/model"
	"myfinace/service"
	"strings"

	"github.com/gofiber/fiber/v2"
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

func RegisterPortfolioController(router fiber.Router) {
	router.Use(middleware.PortfolioMiddleware)
	router.Get("/", HandlePortfolioList)

}

func HandlePortfolioList(c *fiber.Ctx) error {
	queryString := string(c.Request().URI().QueryString())
	data := fiber.Map{
		"Title":       "Portfolio",
		"QueryString": queryString,
	}
	return c.Render("portfolio/portfolios", data, "layouts/main")
}

func (c *PortfolioController) Post() {
	name := c.Ctx.FormValue("portfolio-name")
	des := c.Ctx.FormValueDefault("portfolio-des", "")
	portfolio := model.Portfolio{
		Name:        name,
		Description: des,
	}
	_, err := c.Service.Create(c.Ctx.Request().Context(), portfolio)
	if err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
	c.Ctx.Redirect("/portfolio", iris.StatusSeeOther)
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
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
	c.Ctx.Redirect(fmt.Sprintf("/portfolio/%d", id), iris.StatusSeeOther)
}

func (c *PortfolioController) PostDeleteBy(id uint64) {
	_, err := c.Service.Delete(c.Ctx.Request().Context(), id)
	if err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
	c.Ctx.Redirect("/portfolio", iris.StatusSeeOther)
}

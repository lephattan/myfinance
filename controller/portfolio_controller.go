package controller

import (
	"fmt"
	"myfinace/controller/htmx"
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
	router.Post("/", HandlePortfolioCreate)
	router.Get("/:id", HandlePortfolioDetail)

}

// Get portfolios request
func HandlePortfolioList(c *fiber.Ctx) error {
	queryString := string(c.Request().URI().QueryString())
	data := fiber.Map{
		"Title":       "Portfolio",
		"QueryString": queryString,
	}
	return c.Render("portfolio/portfolios", data, "layouts/main")
}

// Create portfolio request handler
func HandlePortfolioCreate(c *fiber.Ctx) error {
	portfolio := new(model.Portfolio)
	if err := c.BodyParser(portfolio); err != nil {
		return err
	}
	svc, _ := c.Locals("Service").(service.PortfolioService)
	if _, err := svc.Create(c.Context(), *portfolio); err != nil {
		return err
	}
	c.Set("HX-Trigger", "new-portfolio")
	return htmx.HandlePortfolioAddForm(c)
}

// Ticker detail request handler
func HandlePortfolioDetail(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	data := fiber.Map{
		"Title":       "Portfolio",
		"PortfolioID": id,
	}
	return c.Render("portfolio/detail", data, "layouts/main")
	// portfolio := new(model.Portfolio)
	// svc, _ := c.Locals("Service").(service.PortfolioService)
	// err = svc.Get(c.Context(), uint64(id), &portfolio)

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

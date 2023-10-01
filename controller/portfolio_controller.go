package controller

import (
	"myfinace/controller/htmx"
	"myfinace/middleware"
	"myfinace/model"
	"myfinace/service"

	"github.com/gofiber/fiber/v2"
)

func RegisterPortfolioController(router fiber.Router) {
	router.Use(middleware.PortfolioMiddleware)
	router.Get("/", HandlePortfolioList).Name("VPortfolioList")
	router.Post("/", HandlePortfolioCreate)
	router.Get("/:id", HandlePortfolioDetail).Name("VPortfolioDetail")
	router.Put("/:id", HandlePortfolioUpdate)
	router.Delete("/:id", HandlePortfolioDelete)
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

// Portfolio detail request handler
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

}

// Portfolio update request handler
func HandlePortfolioUpdate(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	portfolio := new(model.Portfolio)
	if err := c.BodyParser(portfolio); err != nil {
		return err
	}
	portfolio.ID = uint64(id)
	svc, _ := c.Locals("Service").(service.PortfolioService)
	if _, err = svc.Update(c.Context(), *portfolio); err != nil {
		return err
	}
	return htmx.RenderPortfolioDetail(c, *portfolio)
}

// Portfolio delete request handler

func HandlePortfolioDelete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	svc, _ := c.Locals("Service").(service.PortfolioService)
	if _, err = svc.Delete(c.Context(), uint64(id)); err != nil {
		return err
	}
	c.Set("HX-Redirect", c.App().GetRoute("VPortfolioList").Path)
	return c.SendString("deleted")
}

// func (c *PortfolioController) PostDeleteBy(id uint64) {
// 	_, err := c.Service.Delete(c.Ctx.Request().Context(), id)
// 	if err != nil {
// 		c.Ctx.HTML("<h3>%s</h3>", err.Error())
// 		return
// 	}
// 	c.Ctx.Redirect("/portfolio", iris.StatusSeeOther)
// }

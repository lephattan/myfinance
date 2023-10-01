package htmx

import (
	"myfinace/middleware"
	"myfinace/model"
	"myfinace/service"

	"github.com/gofiber/fiber/v2"
)

// Register handlers to prefix "/htmx/components/portfolio"
func RegisterPortfolioComponentController(router fiber.Router) {
	router.Use(middleware.PortfolioMiddleware)
	router.Get("/list", HandlePortfolioList)
	router.Get("/add-form", HandlePortfolioAddForm)

}

func HandlePortfolioList(c *fiber.Ctx) error {
	errors := []string{}
	queryString := string(c.Request().URI().QueryString())

	var portfolios model.Portfolios
	svc, _ := c.Locals("Service").(service.PortfolioService)
	// url, err := url.ParseRequestURI(c.OriginalURL())
	// if err != nil {
	// 	return err
	// }
	// urlValues := url.Query()
	// opt := portfolios.ParseListOptions(&urlValues)
	err := svc.List(c.Context(), &portfolios)
	if err != nil {
		errors = append(errors, err.Error())
	}
	data := fiber.Map{
		"Title":       "Portfolios",
		"Errors":      errors,
		"QueryString": queryString,
		"Portfolios":  portfolios,
	}
	return c.Render("parts/portfolio/list", data)
}

func HandlePortfolioAddForm(c *fiber.Ctx) error {
	return c.Render("parts/portfolio/add-form", fiber.Map{})
}

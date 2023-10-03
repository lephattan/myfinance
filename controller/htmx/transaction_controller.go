package htmx

import (
	"myfinance/middleware"
	"myfinance/model"
	"myfinance/service"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func RegisterTransactionComponentController(router fiber.Router) {
	router.Use(middleware.TransactionMiddleware)

	router.Get("/list", HandleTransactionList)
	router.Get("/add-form", HandleTransactionAddForm)
}

func HandleTransactionList(c *fiber.Ctx) error {
	queryString := string(c.Request().URI().QueryString())

	var transactions model.Transactions
	svc, _ := c.Locals("Service").(service.TransactionService)
	url, err := url.ParseRequestURI(c.OriginalURL())
	if err != nil {
		return err
	}
	urlValues := url.Query()
	opt := transactions.ParseListOptions(&urlValues)
	if err = svc.List(c.Context(), opt, &transactions); err != nil {
		return err
	}
	data := fiber.Map{
		"QueryString":  queryString,
		"Transactions": transactions,
	}
	return c.Render("parts/transaction/list", data)

}

func HandleTransactionAddForm(c *fiber.Ctx) error {
	return c.Render("parts/transaction/add-form", fiber.Map{})
}

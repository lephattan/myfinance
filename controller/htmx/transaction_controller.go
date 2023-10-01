package htmx

import (
	// "database/sql"
	// "log"
	// "myfinace/database"
	// "myfinace/helper"
	"myfinace/middleware"
	"myfinace/model"
	"myfinace/service"
	"net/url"
	// "strings"
	// "time"

	"github.com/gofiber/fiber/v2"
	"github.com/kataras/iris/v12"
)

type HtmxTransactionController struct {
	Service service.TransactionService
	Ctx     iris.Context
}

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

func (c *HtmxTransactionController) GetAddnewform() {
	data := iris.Map{
		"ID": "add-transaction",
	}
	if err := c.Ctx.View("parts/transaction/add-new-form", data); err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

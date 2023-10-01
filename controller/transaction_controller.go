package controller

import (
	"myfinace/controller/htmx"
	"myfinace/middleware"
	"myfinace/model"
	"myfinace/names"
	"myfinace/service"

	"github.com/gofiber/fiber/v2"
)

func RegisterTransactionController(router fiber.Router) {
	router.Use(middleware.TransactionMiddleware)

	router.Get("/", HandleTransactionsView).Name(names.VTransactionList)
	router.Post("/", HandleTransactionCreate)
}

// Handle get transtions view req
func HandleTransactionsView(c *fiber.Ctx) error {
	queryString := string(c.Request().URI().QueryString())
	data := fiber.Map{
		"Title":       "Portfolio",
		"QueryString": queryString,
	}
	return c.Render("transaction/transactions", data, "layouts/main")
}

// Hanlde create transaction req
func HandleTransactionCreate(c *fiber.Ctx) error {
	var transaction model.Transaction
	if err := c.BodyParser(&transaction); err != nil {
		return err
	}
	svc, _ := c.Locals("Service").(service.TransactionService)
	if _, err := svc.Create(c.Context(), *&transaction); err != nil {
		return err
	}
	c.Set("HX-Trigger", "new-transaction")
	return htmx.HandleTransactionAddForm(c)
}

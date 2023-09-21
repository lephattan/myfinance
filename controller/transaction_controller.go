package controller

import (
	"database/sql"
	"log"
	"myfinace/database"
	"myfinace/model"
	"myfinace/service"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
)

type TransactionController struct {
	Service service.TransactionService
	Ctx     iris.Context
}

func (c *TransactionController) Get() {
	c.Ctx.ViewLayout("main")
	var transactions model.Transactions
	ctx := c.Ctx.Request().Context()
	opt := database.ListOptions{
		Table: "transactions",
	}
	c.Service.List(ctx, opt, &transactions)
	data := iris.Map{
		"Title":        "Transactions",
		"Transactions": transactions,
	}

	if err := c.Ctx.View("transaction/transactions", data); err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

func (c *TransactionController) Post() {
	date := c.Ctx.FormValue("date")
	ticker := c.Ctx.FormValue("ticker")
	transaction_type := c.Ctx.FormValue("type")
	volume := c.Ctx.FormValue("volume")
	commission := c.Ctx.FormValue("commission")
	note := c.Ctx.FormValue("note")
	portfolio_id := c.Ctx.FormValue("portfolio-id")
	ref_id := c.Ctx.FormValue("ref-id")

	errors := []string{}
	date_fmt := "2006-01-02"
	date_only, err := time.Parse(date_fmt, date)
	if err != nil {
		errors = append(errors, "Error parsing date "+err.Error())
	}

	i_volume, err := ParseUint64(volume)
	if err != nil {
		errors = append(errors, "Error parsing volume "+err.Error())
	}

	i_commission, err := ParseUint64(commission)
	if err != nil {
		errors = append(errors, "Error parsing commission "+err.Error())
	}

	i_portfolio, err := ParseUint64(portfolio_id)
	if err != nil {
		errors = append(errors, "Error parsing portfolio_id "+err.Error())
	}

	i_ref, err := ParseUint64(ref_id)
	if err != nil {
		errors = append(errors, "Error parsing ref_id "+err.Error())
	}

	log.Print("Errors: ", errors)
	transaction := model.Transaction{
		Date:            uint64(date_only.Unix()),
		TickerSymbol:    strings.TrimSpace(ticker),
		TransactionType: model.TransactionType(transaction_type),
		Volume:          uint64(i_volume),
		Commission:      uint64(i_commission),
		PortfolioID:     i_portfolio,
		Note:            sql.NullString{String: note, Valid: note == ""},
	}
	if i_ref > 0 {
		transaction.RefID = database.NullInt64{Int64: int64(i_ref), Valid: true}
	}
	log.Print(transaction.String())
	_, err = c.Service.Create(c.Ctx.Request().Context(), transaction)
	if err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
	c.Ctx.Redirect("/transaction", iris.StatusSeeOther)
}

// Parse string into uint64
// Return 0 on empty string
// Return 0 on error
func ParseUint64(s string) (uint64, error) {
	if s == "" {
		return 0, nil
	}
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		i = 0
	}
	return i, err
}

package htmx

import (
	"database/sql"
	"log"
	"myfinace/database"
	"myfinace/helper"
	"myfinace/model"
	"myfinace/service"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
)

type HtmxTransactionController struct {
	Service service.TransactionService
	Ctx     iris.Context
}

func (c *HtmxTransactionController) GetList() {
	errors := []string{}
	var transactions model.Transactions
	ctx := c.Ctx.Request().Context()
	opt := database.ListOptions{
		Table: "transactions",
	}
	if err := c.Service.List(ctx, opt, &transactions); err != nil {
		errors = append(errors, err.Error())

	}
	data := iris.Map{
		"Transactions": transactions,
		"Errors":       errors,
	}

	if err := c.Ctx.View("parts/transaction/list", data); err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
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

func (c *HtmxTransactionController) PostAddnewform() {
	ctx := c.Ctx
	date := ctx.FormValue("date")
	ticker := ctx.FormValue("ticker")
	transaction_type := ctx.FormValue("type")
	volume := ctx.FormValue("volume")
	price := ctx.FormValue("price")
	commission := ctx.FormValue("commission")
	note := ctx.FormValue("note")
	portfolio_id := ctx.FormValue("portfolio-id")
	ref_id := ctx.FormValue("ref-id")

	errors := []string{}
	date_fmt := "2006-01-02"
	date_only, err := time.Parse(date_fmt, date)
	if err != nil {
		errors = append(errors, "Error parsing date "+err.Error())
	}

	i_volume, err := helper.ParseUint64(volume)
	if err != nil {
		errors = append(errors, "Error parsing volume "+err.Error())
	}

	i_price, err := helper.ParseUint64(price)
	if err != nil {
		errors = append(errors, "Error parsing price "+err.Error())
	}

	i_commission, err := helper.ParseUint64(commission)
	if err != nil {
		errors = append(errors, "Error parsing commission "+err.Error())
	}

	i_portfolio, err := helper.ParseUint64(portfolio_id)
	if err != nil {
		errors = append(errors, "Error parsing portfolio_id "+err.Error())
	}

	i_ref, err := helper.ParseUint64(ref_id)
	if err != nil {
		errors = append(errors, "Error parsing ref_id "+err.Error())
	}

	log.Print("Errors: ", errors)
	transaction := model.Transaction{
		Date:            uint64(date_only.Unix()),
		TickerSymbol:    strings.TrimSpace(ticker),
		TransactionType: model.TransactionType(transaction_type),
		Volume:          uint64(i_volume),
		Price:           uint64(i_price),
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
		errors = append(errors, err.Error())
	} else {
		ctx.Header("HX-Trigger", "new-transaction")
	}
	data := iris.Map{
		"ID":     "add-transaction",
		"Errors": errors,
	}

	if err := c.Ctx.View("parts/transaction/add-new-form", data); err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

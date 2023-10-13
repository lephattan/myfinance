package htmx

import (
	"database/sql"
	"errors"
	"log"
	"myfinance/database"
	"myfinance/middleware"
	"myfinance/model"
	"myfinance/names"
	"myfinance/service"

	"github.com/gofiber/fiber/v2"
)

// Register handlers to prefix "/htmx/components/portfolio"
func RegisterPortfolioComponentController(router fiber.Router) {
	router.Use(middleware.PortfolioMiddleware)
	router.Get("/list", HandlePortfolioList)
	router.Get("/add-form", HandlePortfolioAddForm)
	router.Get("/detail/:id", HandlePortfolioDetail)
	router.Get("/edit-form/:id", HandlePortfolioEditForm)
	router.Get("/holding/:id", HandlePortfolioHolding).Name(names.PPortfolioHoldingList)
	router.Get("/holding/:portfolio_id/:symbol", HandlePortfolioSymbolHolding).Name(names.PPortfolioSymbolHolding)
	router.Get("/holding-value/:portfolio_id", HandlePortfolioHoldingValue)
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

func RenderPortfolioDetail(c *fiber.Ctx, portfolio model.Portfolio) error {
	data := fiber.Map{
		"Portfolio": &portfolio,
	}
	return c.Render("parts/portfolio/detail", data)
}

func HandlePortfolioDetail(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	var portfolio model.Portfolio
	svc, _ := c.Locals("Service").(service.PortfolioService)
	err = svc.Get(c.Context(), uint64(id), &portfolio)
	if err != nil {
		return err
	}
	return RenderPortfolioDetail(c, portfolio)
}

func HandlePortfolioEditForm(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	var portfolio model.Portfolio
	svc, _ := c.Locals("Service").(service.PortfolioService)
	err = svc.Get(c.Context(), uint64(id), &portfolio)
	if err != nil {
		return err
	}
	data := fiber.Map{
		"Portfolio": &portfolio,
	}
	return c.Render("parts/portfolio/edit-form", data)
}

func HandlePortfolioHolding(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	holding_svc, ok := c.Locals("HoldingService").(service.HoldingService)
	if !ok {
		return errors.New("Invalid PortfolioService")
	}

	listing_opt := database.ListOptions{
		WhereColumn: "portfolio_id",
		WhereValue:  id,
	}

	var holdings []*model.Holding

	err = holding_svc.List(c.Context(), listing_opt, &holdings)
	if err != nil {
		if err == sql.ErrNoRows {
		} else {
			return err
		}
	}

	data := fiber.Map{
		"Holdings": &holdings,
	}
	return c.Render("parts/portfolio/holding", data)
}

func HandlePortfolioSymbolHolding(c *fiber.Ctx) (err error) {
	id, err := c.ParamsInt("portfolio_id")
	if err != nil {
		log.Printf("Error reading portfolio_id")
		return err
	}

	symbol := c.Params("symbol", "")
	if symbol == "" {
		log.Printf("Error reading ticker symbol")
		return errors.New("Error reading ticker symbol")
	}

	holding_svc, ok := c.Locals("HoldingService").(service.HoldingService)
	if !ok {
		return errors.New("Invalid PortfolioService")
	}

	var holding model.Holding
	listing_opt := database.ListOptions{}

	cond := database.WhereGroup{Operator: "and"}
	cond.Where("portfolio_id", id, "and")
	cond.Where("symbol", symbol, "and")
	listing_opt.WhereGroup(&cond)

	err = holding_svc.Get(c.Context(), listing_opt, &holding)
	if err != nil {
		if err == sql.ErrNoRows {
		} else {
			return err
		}
	}
	return c.Render("parts/portfolio/symbol-holding", holding)
}

func HandlePortfolioHoldingValue(c *fiber.Ctx) (err error) {
	id, err := c.ParamsInt("portfolio_id")
	if err != nil {
		log.Printf("Error reading portfolio_id")
		return err
	}

	svc, ok := c.Locals("Service").(service.PortfolioService)
	if !ok {
		return errors.New("Invalid PortfolioService")
	}
	var holding_total int64
	svc.HoldingValue(c.Context(), uint64(id), &holding_total)
	return c.Render("parts/portfolio/holding-value", holding_total)
}

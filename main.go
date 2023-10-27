package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"

	"myfinance/controller"
	"myfinance/controller/htmx"
	"myfinance/database"
	"myfinance/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Set cutome type parser
	fiber.SetParserDecoder(fiber.ParserConfig{
		IgnoreUnknownKeys: true,
		ParserType: []fiber.ParserType{
			database.UnixDateParser,
			database.NullableInt64Parser,
		},
		ZeroEmpty: true,
	})

	views := MakeViews()

	app := fiber.New(fiber.Config{
		Views:        views,
		ErrorHandler: ErrorHandle,
	})
	MakeAccessLog(app)

	app.Use(AppMiddleWare)

	app.Static("/assets", "./assets")

	controller.RegisterRootController(app.Group("/"))

	controller.RegisterTickerController(app.Group("/ticker"))
	htmx.RegisterTickerComponentController(app.Group("/htmx/components/ticker"))

	controller.RegisterPortfolioController(app.Group("/portfolio"))
	htmx.RegisterPortfolioComponentController(app.Group("/htmx/components/portfolio"))

	controller.RegisterTransactionController(app.Group("/transaction"))
	htmx.RegisterTransactionComponentController(app.Group("/htmx/components/transaction"))

	htmx.RegisterCashflowComponentController(app.Group("/htmx/components/cashflow"))

	log.Fatal(app.Listen("0.0.0.0:8080"))
}

func AppMiddleWare(c *fiber.Ctx) error {
	db := database.GetDB()
	c.Locals("DB", db)
	return c.Next()
}

func MakeViews() *html.Engine {
	engine := html.New("./views", ".html")
	engine.AddFunc("UnixTimeFmt", helper.UnixTimeFmt)
	engine.AddFunc("format", message.NewPrinter(language.English).Sprintf)
	engine.AddFuncMap(map[string]interface{}{
		"json": func(in interface{}) (template.JS, error) {
			out, err := json.Marshal(in)
			return template.JS(out), err
		},
		"devideint64": helper.Devide[int64],
		"minus":       helper.Minus[int64],
		"seq":         helper.Sequence,
		"add":         helper.Add,
		"queryString": helper.QueryString,
	})
	return engine
}

func MakeAccessLog(app *fiber.App) {
	app.Use(logger.New())
}

// Recover and handle error
// Ref: https://docs.gofiber.io/guide/error-handling/
func ErrorHandle(ctx *fiber.Ctx, err error) error {
	err = ctx.SendString(fmt.Sprintf(`<h4 class="errors">%s</h4>`, err))
	return err
}

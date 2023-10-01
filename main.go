package main

import (
	"log"

	"myfinace/controller"
	"myfinace/controller/htmx"
	"myfinace/database"
	"myfinace/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Set cutome type parser
	fiber.SetParserDecoder(fiber.ParserConfig{
		IgnoreUnknownKeys: true,
		ParserType:        []fiber.ParserType{database.UnixDateParser},
		ZeroEmpty:         true,
	})

	views := MakeViews()

	app := fiber.New(fiber.Config{Views: views})
	MakeAccessLog(app)

	app.Use(AppMiddleWare)

	controller.RegisterRootController(app.Group("/"))

	controller.RegisterTickerController(app.Group("/ticker"))
	htmx.RegisterTickerComponentController(app.Group("/htmx/components/ticker"))

	controller.RegisterPortfolioController(app.Group("/portfolio"))
	htmx.RegisterPortfolioComponentController(app.Group("/htmx/components/portfolio"))

	controller.RegisterTransactionController(app.Group("/transaction"))
	htmx.RegisterTransactionComponentController(app.Group("/htmx/components/transaction"))

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
	return engine
}

func MakeAccessLog(app *fiber.App) {
	app.Use(logger.New())
}

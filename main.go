package main

import (
	"log"

	"myfinace/controller"
	"myfinace/controller/htmx"
	"myfinace/database"
	"myfinace/env"
	"myfinace/helper"
	"myfinace/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// app_env := env.ReadEnv("APP_ENV", "production")
	views := MakeViews()

	app := fiber.New(fiber.Config{Views: views})
	MakeAccessLog(app)

	app.Use(AppMiddleWare)

	controller.RegisterRootController(app.Group("/"))

	controller.RegisterTickerController(app.Group("/ticker"))
	htmx.RegisterTickerComponentController(app.Group("/htmx/components/ticker"))

	controller.RegisterPortfolioController(app.Group("/portfolio"))
	htmx.RegisterPortfolioComponentController(app.Group("/htmx/components/portfolio"))

	log.Fatal(app.Listen("0.0.0.0:8080"))
	// mvc.Configure(app.Party("portfolio"), portfolioSetup)
	// mvc.Configure(app.Party("transaction"), transactionSetup)
	//
	// // HTMX
	// mvc.Configure(app.Party("/htmx/components/transaction"), htmxComponentSetup)
	// mvc.Configure(app.Party("/htmx/components/ticker"), htmxTickerSetup)
	//
	// app.Get("/", getRoot)
	// app.Listen("0.0.0.0:8080")
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

func getRoot(ctx iris.Context) {
	data := iris.Map{
		"Title": "My Finance",
	}
	ctx.ViewLayout("main")
	if err := ctx.View("index", data); err != nil {
		ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

func setup(app *mvc.Application) {
	app_env := env.ReadEnv("APP_ENV", "production")
	app.Register(
		app_env,
		database.NewDB,
		service.NewGreetService,
	)
	app.Handle(new(controller.GreetController))
}

func portfolioSetup(app *mvc.Application) {
	app_env := env.ReadEnv("APP_ENV", "production")
	app.Register(
		app_env,
		database.NewDB,
		service.NewPortfolioService,
	)

	app.Handle(new(controller.PortfolioController))
}

func transactionSetup(app *mvc.Application) {
	app_env := env.ReadEnv("APP_ENV", "production")
	app.Register(
		app_env,
		database.NewDB,
		service.NewTransactionService,
	)

	app.Handle(new(controller.TransactionController))
}

func htmxComponentSetup(app *mvc.Application) {
	app_env := env.ReadEnv("APP_ENV", "production")
	app.Register(
		app_env,
		database.NewDB,
		service.NewTransactionService,
	)
	app.Handle(new(htmx.HtmxTransactionController))
}

package main

import (
	"log"
	"os"

	"myfinace/controller"
	"myfinace/database"
	"myfinace/env"
	"myfinace/service"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	ac := makeAccessLog()
	defer ac.Close()

	app := iris.New()
	app.Logger()

	app.UseRouter(ac.Handler)
	app.Get("/ping", pong).Describe("health check")

	app.RegisterView(iris.Blocks("./views", ".html").Reload(true))

	mvc.Configure(app.Party("greet"), setup)
	mvc.Configure(app.Party("ticker"), tickerSetup)
	mvc.Configure(app.Party("portfolio"), portfolioSetup)
	mvc.Configure(app.Party("transaction"), transactionSetup)

	app.Get("/", getRoot)
	app.Listen("0.0.0.0:8080")
}

func makeAccessLog() *accesslog.AccessLog {
	ac := accesslog.File("./access.log")
	ac.AddOutput(os.Stdout)
	// The default configuration:
	ac.Delim = '|'
	ac.TimeFormat = "2006-01-02 15:04:05"
	ac.Async = false
	ac.IP = true
	ac.BytesReceivedBody = true
	ac.BytesSentBody = true
	ac.BytesReceived = false
	ac.BytesSent = false
	ac.BodyMinify = true
	ac.RequestBody = true
	ac.ResponseBody = false
	ac.KeepMultiLineError = true
	ac.PanicLog = accesslog.LogHandler
	// Default line format if formatter is missing:
	// Time|Latency|Code|Method|Path|IP|Path Params Query Fields|Bytes Received|Bytes Sent|Request|Response|
	//
	// Set Custom Formatter:
	ac.SetFormatter(&accesslog.JSON{
		Indent:    "  ",
		HumanTime: true,
	})
	// ac.SetFormatter(&accesslog.CSV{})
	// ac.SetFormatter(&accesslog.Template{Text: "{{.Code}}"})
	return ac
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

func pong(ctx iris.Context) {
	ctx.WriteString("pong")
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

func tickerSetup(app *mvc.Application) {
	app_env := env.ReadEnv("APP_ENV", "production")
	app.Register(
		app_env,
		database.NewDB,
		service.NewTickerService,
	)

	app.Handle(new(controller.TickerController))
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

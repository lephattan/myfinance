package controller

import (
	"myfinace/service"

	"github.com/kataras/iris/v12"
)

type Request struct {
	Name string `url:"name"`
}

type GreetController struct {
	Service service.GreetSerivce
	Ctx     iris.Context
}

func (c *GreetController) Get(req Request) {
	message, err := c.Service.Say(req.Name)
	if err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}

	data := iris.Map{
		"Title":   "Greeting",
		"Message": message,
	}
	c.Ctx.ViewLayout("main")
	if err := c.Ctx.View("greet", data); err != nil {
		c.Ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

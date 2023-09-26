package controller

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRootController(group fiber.Router) {
	group.Get("/", IndexHandler)
	group.Get("/ping", PingHandler)
}

func PingHandler(c *fiber.Ctx) error {
	return c.SendString("pong")
}

func IndexHandler(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{}, "layouts/main")
}

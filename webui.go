package main

import (
	"log/slog"

	"github.com/a-h/templ"
	"github.com/alexng353/ihostproxy/views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

type GlobalState struct {
	Count int
}

var global GlobalState

func getHandler(c *fiber.Ctx) error {
	// component := views.Page()

	handler := adaptor.HTTPHandler(templ.Handler(views.Page()))

	return handler(c)
}

func startWebUI(ctx Env) {
	WebUIPort := ctx.WebUIPort

	app := fiber.New()

	app.Get("/", getHandler)

	app.Get("/goodbye", func(c *fiber.Ctx) error {
		return c.SendString("Goodbye, World!")
	})

	slog.Info("Starting web ui on port " + WebUIPort)
	app.Listen(":8080")
}

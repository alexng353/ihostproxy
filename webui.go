package main

import (
	"log/slog"
	"strconv"

	"github.com/a-h/templ"
	"github.com/alexng353/ihostproxy/views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

type GlobalState struct {
	Count int
}

var global GlobalState

func templAdapter(component templ.Component) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		handler := adaptor.HTTPHandler(templ.Handler(component))
		return handler(c)
	}
}

func startWebUI(ctx Env) {
	port := strconv.FormatInt(int64(validatePort(ctx.WebUIPort, 8080)), 10)

	app := fiber.New()

	app.Get("/", templAdapter(views.Page()))

	slog.Info("Starting WebUI", "port", port)
	app.Listen(":" + port)
}

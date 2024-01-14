package webui

import (
	"log/slog"
	"strconv"

	"github.com/a-h/templ"
	"github.com/alexng353/ihostproxy/context"
	"github.com/alexng353/ihostproxy/helpers"
	"github.com/alexng353/ihostproxy/views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

type GlobalState struct {
	Count int
}

var ctx = context.GetEnv()

func templAdapter(component templ.Component) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		handler := adaptor.HTTPHandler(templ.Handler(component))
		return handler(c)
	}
}

func StartWebUI(ctx context.Env) {
	port := strconv.FormatInt(int64(helpers.ValidatePort(ctx.WebUIPort, 8080)), 10)

	app := fiber.New()

	var StaticOptions = fiber.Static{
		CacheDuration: -1,
		MaxAge:        0,
	}
	app.Static("/static", "./static", StaticOptions)

	app.Get("/login", templAdapter(views.Login()))
	app.Post("/api/login", login)
	app.Get("/api/auth", auth)

	app.Use(authMiddleware)

	app.Get("/", templAdapter(views.Page()))
	app.Get("/users", templAdapter(views.Users()))
	app.Get("/logout", templAdapter(views.Logout()))

	app.Get("/api/getusers", getUsers)
	app.Post("/api/adduser", addUser)
	app.Post("/api/deleteuser", deleteUser)
	app.Post("/api/logout", logout)

	slog.Info("Starting WebUI", "port", port)
	app.Listen(":" + port)
}

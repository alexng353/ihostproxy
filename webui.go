package main

import (
	"log/slog"
	"strconv"

	"github.com/a-h/templ"
	"github.com/alexng353/ihostproxy/views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/golang-jwt/jwt/v5"
)

type GlobalState struct {
	Count int
}

var global GlobalState
var credentialStore *SQLiteCredentialStore = NewSQLiteCredentialStore()

func templAdapter(component templ.Component) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		handler := adaptor.HTTPHandler(templ.Handler(component))
		return handler(c)
	}
}

func login(c *fiber.Ctx) error {
	slog.Info("login", "method", c.Method())
	if c.Method() == "POST" {
		username := c.FormValue("username")
		password := c.FormValue("password")

		slog.Info("login post", "username", username, "password", password)

		id, err := credentialStore.GetEntry(username, password)
		if err != nil {
			slog.Error("login post", "err", err)
			return c.SendString("login post err")
		}

		slog.Info("login post", "id", id)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})

		tokenString, err := token.SignedString([]byte("secret"))
		if err != nil {
			slog.Error("login post", "err", err)
			return c.SendString("login post err")
		}

		slog.Info("login post", "token", tokenString)

		c.Cookie(&fiber.Cookie{
			Name:  "token",
			Value: tokenString,
		})

		return c.SendString("login post")
	}

	return c.SendString("login")
}

func startWebUI(ctx Env) {
	port := strconv.FormatInt(int64(validatePort(ctx.WebUIPort, 8080)), 10)

	app := fiber.New()

	var StaticOptions = fiber.Static{
		CacheDuration: -1,
		MaxAge:        0,
	}
	app.Static("/static", "./static", StaticOptions)

	app.Get("/", templAdapter(views.Page()))
	app.Get("/login", login)
	app.Post("/login", login)

	slog.Info("Starting WebUI", "port", port)
	app.Listen(":" + port)
}

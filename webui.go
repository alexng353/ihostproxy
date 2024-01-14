package main

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/alexng353/ihostproxy/credentials"
	"github.com/alexng353/ihostproxy/pika"
	"github.com/alexng353/ihostproxy/views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/golang-jwt/jwt/v5"
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

func login(c *fiber.Ctx) error {
	var creds = credentials.Get()

	slog.Info("login", "method", c.Method())
	if c.Method() == "POST" {
		username := c.FormValue("username")
		password := c.FormValue("password")

		user, err := creds.GetEntry(username, password)
		if err != nil {
			slog.Error("login post", "err", err)
			return c.SendString("login post err")
		}

		if !user.Admin {
			return c.SendString("not an admin")
		}

		// JWT Claims
		jti := pika.Gen("jti")
		now := time.Now()
		iat := now.Unix()
		// Expire after 14 days
		exp := now.Add(time.Hour * 24 * 14).Unix()

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.Id,
			"usr": username,
			"jti": jti,
			"iat": iat,
			"exp": exp,
		})

		tokenString, err := token.SignedString([]byte(ctx.JwtSecret))
		if err != nil {
			slog.Error("login post", "err", err)
			return c.SendString("login post err")
		}

		c.Cookie(&fiber.Cookie{
			Name:  "token",
			Value: tokenString,
		})

		return c.SendString("Welcome " + username)
	}

	return c.SendString("login")
}

func addUser(c *fiber.Ctx) error {
	slog.Info("addUsers", "method", c.Method())
	if c.Method() == "POST" {
		user, err := validateUser(c)
		if err != nil {
			slog.Error("addUsers post", "err", err)
			return c.SendString("addUsers post err")
		}

		if !user.Admin {
			return c.SendString("not an admin")
		}

		username := c.FormValue("username")
		password := c.FormValue("password")

		err = credentials.Get().AddEntry(username, password)
		if err != nil {
			slog.Error("addUsers post", "err", err)
			return c.SendString("Failed to add user " + username + " " + err.Error())
		}

		return c.SendString("Added " + username)
	}

	return c.SendString("addUsers")
}

func getUsers(c *fiber.Ctx) error {
	slog.Info("getUsers", "method", c.Method())
	if c.Method() == "GET" {

		user, err := validateUser(c)

		if err != nil {
			slog.Error("getUsers post", "err", err)
			return c.SendString("getUsers post err")
		}
		if !user.Admin {
			return c.SendString("not an admin")
		}

		users, err := credentials.Get().GetUsers()
		if err != nil {
			slog.Error("getUsers post", "err", err)
			return c.SendString("Failed to get users " + err.Error())
		}

		return templAdapter(views.UserList(users))(c)

	}

	return c.SendString("getUsers")
}

func deleteUser(c *fiber.Ctx) error {
	if c.Method() != "POST" {
		return c.SendString("deleteUser")
	}

	user, err := validateUser(c)
	if err != nil {
		return c.SendString("error validating user")
	}

	if !user.Admin {
		return c.SendString("not an admin")
	}

	id_to_remove := c.Query("id")
	if id_to_remove == "" {
		return c.SendString("no id provided")
	}

	if id_to_remove == user.Id {
		return c.SendString("cannot delete yourself")
	}

	err = credentials.Get().RemoveEntry(id_to_remove)
	if err != nil {
		return c.SendString("error deleting user")
	}

	return c.SendString("deleted user " + id_to_remove)

}

func auth(c *fiber.Ctx) error {
	cookie := c.Cookies("token")
	if cookie == "" {
		return c.SendString("Authenticate yourself!")
	}

	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(ctx.JwtSecret), nil
	})

	if err != nil {
		slog.Error("auth", "err", err)
		return err
	}

	sub := claims["sub"].(string)

	user, err := credentials.Get().GetUser(sub)
	if err != nil {
		slog.Error("auth", "err", err)
		return err
	}

	if !user.Admin {
		return c.SendString("not an admin")
	}

	return c.SendString("Welcome " + user.Username)
}

func validateUser(c *fiber.Ctx) (*credentials.User, error) {
	cookie := c.Cookies("token")

	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(ctx.JwtSecret), nil
	})

	if err != nil {
		slog.Error("validateUser", "err", err)
		return nil, err
	}

	sub := claims["sub"].(string)

	user, err := credentials.Get().GetUser(sub)
	if err != nil {
		slog.Error("validateUser", "err", err)
		return nil, err
	}

	return user, nil
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

	app.Post("/api/adduser", addUser)

	app.Get("/users", templAdapter(views.Users()))
	app.Get("/api/getusers", getUsers)

	app.Post("/api/deleteuser", deleteUser)

	app.Get("/api/auth", auth)

	slog.Info("Starting WebUI", "port", port)
	app.Listen(":" + port)
}

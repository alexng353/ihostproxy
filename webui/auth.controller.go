package webui

import (
	"log/slog"
	"os"
	"time"

	"github.com/alexng353/ihostproxy/credentials"
	"github.com/alexng353/ihostproxy/pika"
	"github.com/alexng353/ihostproxy/views"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

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
		exp := now.Add(time.Hour * 24 * 14)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.Id,
			"usr": username,
			"jti": jti,
			"iat": iat,
			"exp": exp.Unix(),
		})

		tokenString, err := token.SignedString([]byte(ctx.JwtSecret))
		if err != nil {
			slog.Error("login post", "err", err)
			return c.SendString("login post err")
		}

		var secure bool = false
		var isSecure = os.Getenv("SECURE")
		if isSecure == "true" {
			secure = true
		}

		c.Cookie(&fiber.Cookie{
			Name:     "token",
			Value:    tokenString,
			Expires:  exp,
			Secure:   secure,
			HTTPOnly: true,
			SameSite: "Strict",
		})

		return templAdapter(views.AuthRedirect(user.Username))(c)
	}

	return templAdapter(views.Login())(c)
}

func logout(c *fiber.Ctx) error {
	if c.Method() != "POST" {
		return c.SendString("logout")
	}

	c.Cookie(&fiber.Cookie{
		Name: "token",
	})

	c.ClearCookie("token")
	c.Set("HX-Redirect", "/login")

	return c.Redirect("/login")
}

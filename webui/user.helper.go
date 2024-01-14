package webui

import (
	"log/slog"

	"github.com/alexng353/ihostproxy/credentials"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

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

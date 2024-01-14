package webui

import "github.com/gofiber/fiber/v2"

func authMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("token")
	if cookie == "" {
		c.Set("HX-Redirect", "/login")
		return c.Redirect("/login")
	}

	user, err := validateUser(c)
	if err != nil {
		c.Set("HX-Redirect", "/login")
		return c.Redirect("/login")
	}

	if !user.Admin {
		c.Cookie(&fiber.Cookie{
			Name: "token",
		})
		c.ClearCookie("token")
		c.Set("HX-Redirect", "/login")
		return c.Redirect("/login")
	}

	return c.Next()
}

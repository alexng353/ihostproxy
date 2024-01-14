package webui

import (
	"log/slog"

	"github.com/alexng353/ihostproxy/credentials"
	"github.com/alexng353/ihostproxy/views"
	"github.com/gofiber/fiber/v2"
)

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

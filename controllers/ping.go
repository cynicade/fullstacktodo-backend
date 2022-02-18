package controllers

import "github.com/gofiber/fiber/v2"

func Ping(c *fiber.Ctx) error {
	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"message": "pong",
	})
}

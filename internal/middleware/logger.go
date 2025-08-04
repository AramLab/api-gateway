package middleware

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Printf("Request: %s %s", c.Method(), c.Path())
		return c.Next()
	}
}

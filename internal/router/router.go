package router

import (
	"github.com/AramLab/api-gateway/internal/handlers"
	"github.com/AramLab/api-gateway/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Routers struct {
	TaskHandler handlers.TaskHandler
	AuthHandler handlers.AuthHandler
}

func NewRouters(r *Routers, secret string) *fiber.App {
	app := fiber.New()

	// Настройка CORS (разрешенные методы, заголовки, авторизация)
	app.Use(cors.New(cors.Config{
		AllowMethods:  "GET, POST",
		AllowHeaders:  "Accept, Authorization, Content-Type, X-CSRF-Token, X-REQUEST-ID",
		ExposeHeaders: "Link",
		MaxAge:        300,
	}))

	appGroup := app.Group("/v1")

	appGroup.Post("/register", r.AuthHandler.Register())
	appGroup.Post("/login", r.AuthHandler.Login())

	appGroupProtect := appGroup.Group("/tasks", middleware.AuthMiddleware([]byte(secret)), middleware.LoggerMiddleware())

	appGroupProtect.Get("/", r.TaskHandler.GetTasks())
	appGroupProtect.Post("/", r.TaskHandler.CreateTask())

	return app
}

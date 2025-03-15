package routes

import (
	"agungmohmd/sematin-front-api/server/handlers"

	"github.com/gofiber/fiber/v2"
)

type AuthRoute struct {
	RouterGroup fiber.Router
	Handler     handlers.Handler
}

func (route AuthRoute) RegisterRoute() {
	handler := handlers.AuthHandler{Handler: route.Handler}
	r := route.RouterGroup.Group("/api/auth")

	r.Post("/", handler.Login)
	r.Post("/register", handler.Register)
}

package routes

import (
	"agungmohmd/intikm-test-api/server/handlers"
	"agungmohmd/intikm-test-api/server/middlewares"

	"github.com/gofiber/fiber/v2"
)

type UserRoute struct {
	RouterGroup fiber.Router
	Handler     handlers.Handler
}

func (route UserRoute) RegisterRoute() {
	handler := handlers.AuthHandler{Handler: route.Handler}
	jwtMiddlewareOutlet := middlewares.JwtMiddleware{
		ContractUC: handler.ContractUC,
	}
	r := route.RouterGroup.Group("/api/users")
	r.Use(jwtMiddlewareOutlet.New)
	r.Post("/update-profile", handler.UpdateProfile)
}

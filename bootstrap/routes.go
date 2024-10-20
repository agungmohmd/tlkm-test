package bootstrap

import (
	"agungmohmd/gateway-api/server/bootstrap/routes"
	"agungmohmd/gateway-api/server/handlers"
)

func (boot Bootstrap) RegisterRouters() {
	handler := handlers.Handler{
		FiberApp:   boot.App,
		ContractUC: &boot.ContractUC,
	}

	apiV1 := boot.App.Group("/V1")

	authRoutes := routes.AuthRoute{RouterGroup: apiV1, Handler: handler}
	authRoutes.RegisterRoute()
	userRoutes := routes.UserRoute{RouterGroup: apiV1, Handler: handler}
	userRoutes.RegisterRoute()

}

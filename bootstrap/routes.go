package bootstrap

import (
	"agungmohmd/intikm-test-api/server/bootstrap/routes"
	"agungmohmd/intikm-test-api/server/handlers"
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

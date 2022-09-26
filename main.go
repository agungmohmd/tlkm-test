package main

import (
	"agungmohmd/intikm-test-api/config"
	"agungmohmd/intikm-test-api/server/bootstrap"
	"agungmohmd/intikm-test-api/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/xid"
)

func main() {
	app := fiber.New()
	conf, _ := config.LoadConfig()
	defer conf.DB.Close()

	ContractUc := usecase.ContractUC{
		ReqID: xid.New().String(),
		DB:    conf.DB,
	}

	boot := bootstrap.Bootstrap{
		App:        app,
		ContractUC: ContractUc,
	}

	boot.App.Use(recover.New())
	boot.App.Use(requestid.New())
	boot.App.Use(cors.New())
	boot.RegisterRouters()
	boot.App.Listen(":3000")
}

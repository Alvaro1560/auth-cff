package api

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/api/rest/handler/auth"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/api/rest/handler/look_and_feel"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/api/rest/handler/users"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func routes(db *sqlx.DB, loggerHttp bool, allowedOrigins string) *fiber.App {
	app := fiber.New()

	prometheus := fiberprometheus.New("ecatchAuth")
	prometheus.RegisterAt(app, "/metrics")

	app.Use(recover.New())
	app.Use(prometheus.Middleware)
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowHeaders: "Origin, X-Requested-With, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST",
	}))
	if loggerHttp {
		app.Use(logger.New())
	}
	TxID := uuid.New().String()

	auth.AuthenticationRouter(app, db, TxID)
	register.UserRouter(app, db, TxID)
	look_and_feel.LookAndFeel(app, db, TxID)
	return app
}

package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func AuthenticationRouter(app *fiber.App, db *sqlx.DB, tx string) {

	ln := Handler{DB: db, TxID: tx}

	api := app.Group("/api")
	v3 := api.Group("/v3")
	v3.Post("/auth", ln.LoginV3)

	v4 := api.Group("/v4")
	v4.Post("/auth", ln.Login)
}

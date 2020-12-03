package register

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func UserRouter(app *fiber.App, db *sqlx.DB, tx string) {

	usr := Handler{DB: db, TxID: tx}
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/register", usr.CreateUser)

}

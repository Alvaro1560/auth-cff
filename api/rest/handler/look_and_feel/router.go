package look_and_feel

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func LookAndFeel(app *fiber.App, db *sqlx.DB, tx string) {

	laf := Handler{DB: db, TxID: tx}

	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/laf", laf.LockAndFeel)

}

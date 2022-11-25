package ciphers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func CipherRouter(app *fiber.App, db *sqlx.DB, tx string) {

	cp := Handler{DB: db, TxID: tx}

	api := app.Group("/api")
	v1 := api.Group("/v1/cipher")
	v1.Get("/encrypt/:text", cp.encrypt)
	v1.Post("/decrypt", cp.decrypt)
	v1.Get("/key-app", cp.getKeyApp)
}

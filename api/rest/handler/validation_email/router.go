package validation_email

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func ValidationEmailRouter(app *fiber.App, db *sqlx.DB, tx string) {

	ln := handlerValidationEmail{DB: db, TxID: tx}

	api := app.Group("/api")
	v1 := api.Group("/v1")
	validation := v1.Group("/email-verify")
	validation.Post("/generate", ln.sendCode)
	validation.Post("/generate-otp", ln.GenerateOtp)
	validation.Post("/validate", ln.verifyCode)

}

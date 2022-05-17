package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func AuthenticationRouter(app *fiber.App, db *sqlx.DB, tx string) {

	ln := Handler{DB: db, TxID: tx}

	api := app.Group("/api")
	v1 := api.Group("/v1/auth")
	v1.Post("/forgot-password", ln.ForgotPassword)
	v1.Post("/change-password", ln.ChangePassword)
	v1.Post("/password-policy", ln.PasswordPolicy)
	v1.Post("/autologin", ln.LoginGeneric)

	v3 := api.Group("/v3/auth")
	v3.Post("", ln.LoginV3)

	v4 := api.Group("/v4/auth")
	v4.Post("", ln.Login)
}

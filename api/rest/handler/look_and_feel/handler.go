package look_and_feel

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/response"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/msgs"

	"github.com/jmoiron/sqlx"
)

type Handler struct {
	DB   *sqlx.DB
	TxID string
}

//func Login(c echo.Context) error {
func (h *Handler) LockAndFeel(c *fiber.Ctx) error {
	res := response.Model{Error: true}
	var msg msgs.Model
	m := Model{
		LoginNameProject: "Ecatch",
		LoginVersion:     "6.1.0",
		LoginSlogan:      "Sistema para la administración de contenido empresarial",
		LoginLogo:        "assets/img/logo_ecapture_blanco_h.svg",
		MenuLogo:         "assets/img/logo_ecatch_new.svg",
		FooterLogo:       "",
		Primary:          "#353A48",
		Secondary:        "#039BE5",
		Tertiary:         "#262933",
		ID:               "****",
		Key:              "****",
	}
	res.Data = m
	res.Code, res.Type, res.Msg = msg.GetByCode(29)
	res.Error = false

	return c.Status(http.StatusOK).JSON(res)
}

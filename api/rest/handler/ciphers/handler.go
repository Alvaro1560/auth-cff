package ciphers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/ciphers"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/msgs"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/response"
	"net/http"
)

type Handler struct {
	DB   *sqlx.DB
	TxID string
}

//func Login(c echo.Context) error {
func (h *Handler) encrypt(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(CipherResponse{
		Text: ciphers.Encrypt(c.Params("text")),
	})
}

//func Login(c echo.Context) error {
func (h *Handler) decrypt(c *fiber.Ctx) error {
	res := response.Model{Error: true}
	var msg msgs.Model
	m := CipherRequest{}
	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf("no se pudo leer el Modelo CipherRequest: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	rsDecrypt , err := ciphers.Decrypt(m.TextDecrypt)
	cr := CipherResponse{
		Text: string(rsDecrypt),
	}
	if err != nil  {
		logger.Error.Println(err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	res.Data = cr
	return c.Status(http.StatusOK).JSON(res)
}

//func Login(c echo.Context) error {
func (h *Handler) getSecretKey(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(CipherResponse{
		SecretKey: []byte(ciphers.GetSecret()),
	})
}

package auth

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/response"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth/login"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/msgs"

	"github.com/jmoiron/sqlx"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
)

type Handler struct {
	DB   *sqlx.DB
	TxID string
}

//func Login(c echo.Context) error {
func (h *Handler) LoginV3(c *fiber.Ctx) error {
	res := response.Model{Error: true}
	var msg msgs.Model
	m := LoginRequest{}
	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf("no se pudo leer el Modelo User en login: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	m.RealIP = c.IP()
	serviceLogin := login.NewLoginService(h.DB, h.TxID)
	token, cod, err := serviceLogin.Login(m.ID, m.Username, m.Password, m.ClientID, m.HostName, m.RealIP)
	if err != nil {
		logger.Warning.Printf("no se pudo leer el Modelo User en login: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(cod)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	res.Data = token
	res.Code, res.Type, res.Msg = msg.GetByCode(cod)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	res := response.Model{Error: true}
	var msg msgs.Model
	m := LoginRequest{}
	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf("no se pudo leer el Modelo User en login: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	m.RealIP = c.IP()
	serviceLogin := login.NewLoginService(h.DB, h.TxID)
	token, cod, err := serviceLogin.Login(m.ID, m.Username, m.Password, m.ClientID, m.HostName, m.RealIP)
	if err != nil {
		logger.Warning.Printf("no se pudo leer el Modelo User en login: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(cod)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	mr := LoginResponse{
		AccessToken:  token,
		RefreshToken: token,
	}
	res.Data = mr
	res.Code, res.Type, res.Msg = msg.GetByCode(cod)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

func (h *Handler) ForgotPassword(c *fiber.Ctx) error {
	res := response.Model{Error: true}
	var msg msgs.Model
	m := ForgotPasswordRequest{}
	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf("no se pudo leer el forgot password: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

func (h *Handler) ChangePassword(c *fiber.Ctx) error {
	res := response.Model{Error: true}
	var msg msgs.Model
	m := ChangePasswordRequest{}
	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf("no se pudo leer el forgot password: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

func (h *Handler) PasswordPolicy(c *fiber.Ctx) error {
	res := response.Model{Error: true}
	var msg msgs.Model
	m := PasswordPolicyRequest{}
	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf("no se pudo leer el Modelo Password para validar politicas: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		res.Data = false
		return c.Status(http.StatusOK).JSON(res)
	}
	if len(m.Password) < 4 {
		res.Code, res.Type, res.Msg = msg.GetByCode(77)
		res.Data = false
		return c.Status(http.StatusOK).JSON(res)
	}
	res.Data = true
	res.Code, res.Type, res.Msg = msg.GetByCode(29)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

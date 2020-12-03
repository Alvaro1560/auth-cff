package register

import (
	"github.com/google/uuid"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/password"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth/users"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/response"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/msgs"

	"github.com/jmoiron/sqlx"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
)

type Handler struct {
	DB   *sqlx.DB
	TxID string
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	res := response.Model{Error: true}
	var msg msgs.Model
	var id string
	m := UserRequest{}
	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf("no se pudo leer el Modelo User en login: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	if m.ID != nil {
		id = *m.ID
	} else {
		id = uuid.New().String()
	}
	if m.Password != m.PasswordConfirm {
		logger.Error.Printf("password y passwordConfirm no coinciden: %v ", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(74)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	repositoryUsers := users.FactoryStorage(h.DB, nil, h.TxID)
	serviceUsers := users.NewUserService(repositoryUsers, nil, h.TxID)
	m.Password = password.Encrypt(m.Password)
	m.PasswordConfirm = ""
	user, cod, err := serviceUsers.CreateUser(id, m.Username, m.Name, m.LastName, m.Password,
		m.EmailNotifications, m.IdentificationNumber, m.IdentificationType)
	if err != nil {
		logger.Error.Println("couldn't create user")
		res.Code, res.Type, res.Msg = msg.GetByCode(cod)
		if cod == 15 {
			er := err.Error()
			res.Msg = er
		}
		return c.Status(http.StatusAccepted).JSON(res)
	}
	m.ID = &id
	m.Password = ""
	if user == nil {
		logger.Error.Println("couldn't create user")
		res.Code, res.Type, res.Msg = msg.GetByCode(15)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(cod)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

func (h *Handler) ExistUser(c *fiber.Ctx) error {
	res := response.Model{Error: true}
	var msg msgs.Model
	user := c.Query("user")
	if len(user) < 4 {
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusOK).JSON(res)
		res.Data = false
	}
	res.Data = true
	res.Code, res.Type, res.Msg = msg.GetByCode(29)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

func (h *Handler) ExistEmail(c *fiber.Ctx) error {
	res := response.Model{Error: true}
	var msg msgs.Model
	user := c.Query("email")
	if len(user) < 4 {
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		res.Data = false
		return c.Status(http.StatusOK).JSON(res)
	}
	res.Data = true
	res.Code, res.Type, res.Msg = msg.GetByCode(29)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

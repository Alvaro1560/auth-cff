package register

import (
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/ciphers"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/password"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth/roles_password_policy"
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
		logger.Error.Printf(h.TxID, "no se pudo leer el Modelo User en login: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	if m.ID != nil {
		id = *m.ID
	} else {
		id = uuid.New().String()
	}

	passByte := ciphers.Decrypt(m.Password)
	if passByte == "" {
		logger.Error.Printf(h.TxID, "no se pudo leer el Modelo User en login: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	m.Password = passByte

	passConfirmByte := ciphers.Decrypt(m.PasswordConfirm)
	if passConfirmByte == "" {
		logger.Error.Printf(h.TxID, "no se pudo leer el Modelo User en login: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	m.PasswordConfirm = passConfirmByte

	if m.Password != m.PasswordConfirm {
		logger.Error.Printf(h.TxID, "password y passwordConfirm no coinciden: %v ", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(74)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	resultValidStruct, err := govalidator.ValidateStruct(m)
	if err != nil {
		logger.Error.Printf(h.TxID, "Error en validación de datos : %v ", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(15)
		res.Msg = err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if !resultValidStruct {
		logger.Error.Printf(h.TxID, "No cumple la validación de datos : %v ", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(15)
		res.Msg = err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	repositoryRPasswordPolicy := roles_password_policy.FactoryStorage(h.DB, nil, h.TxID)
	servicesRoles := roles_password_policy.NewRolesPasswordPolicyService(repositoryRPasswordPolicy, nil, h.TxID)
	rs := []string{"50602690-B91F-4567-9A8D-A812B37A87BF"}
	pp, err := servicesRoles.GetAllRolesPasswordPolicyByRolesIDs(rs)
	if err != nil {
		logger.Error.Println("couldn't get role to validate passwordPolicy")
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	if pp == nil {
		logger.Error.Println("don't exists role to validate passwordPolicy")
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	repositoryUsers := users.FactoryStorage(h.DB, nil, h.TxID)
	serviceUsers := users.NewUserService(repositoryUsers, nil, h.TxID)
	var result bool
	for _, policy := range pp {
		valid, cod, err := serviceUsers.ValidatePasswordPolicy(m.Password, policy.MaxLength, policy.MinLength, policy.Alpha,
			policy.Digits, policy.Special, policy.UpperCase, policy.LowerCase, policy.Enable)
		if err != nil {
			logger.Error.Println("couldn't get password to validate")
			res.Code, res.Type, res.Msg = msg.GetByCode(cod)
			return c.Status(http.StatusAccepted).JSON(res)
		}
		result = valid
	}
	if !result {
		logger.Error.Println("Password no cumple politicas del rol")
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	m.Password = password.Encrypt(m.Password)
	m.PasswordConfirm = ""
	user, cod, err := serviceUsers.CreateUser(id, m.Username, m.Name, m.LastName, m.Password,
		m.EmailNotifications, m.IdentificationNumber, m.IdentificationType)
	if err != nil {
		logger.Error.Println("couldn't create user: ", err)
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

func (h *Handler) ValidatePassword(c *fiber.Ctx) error {
	res := response.Model{Error: true}
	var msg msgs.Model
	pass := c.Query("password")
	repositoryUsers := users.FactoryStorage(h.DB, nil, h.TxID)
	serviceUsers := users.NewUserService(repositoryUsers, nil, h.TxID)
	repositoryRPasswordPolicy := roles_password_policy.FactoryStorage(h.DB, nil, h.TxID)
	servicesRoles := roles_password_policy.NewRolesPasswordPolicyService(repositoryRPasswordPolicy, nil, h.TxID)
	rs := []string{"50602690-B91F-4567-9A8D-A812B37A87BF"}
	pp, err := servicesRoles.GetAllRolesPasswordPolicyByRolesIDs(rs)
	if err != nil {
		logger.Error.Println("couldn't get role to validate passwordPolicy")
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	if pp == nil {
		logger.Error.Println("don't exists role to validate passwordPolicy")
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	var result bool
	for _, policy := range pp {
		valid, cod, err := serviceUsers.ValidatePasswordPolicy(pass, policy.MaxLength, policy.MinLength, policy.Alpha,
			policy.Digits, policy.Special, policy.UpperCase, policy.LowerCase, policy.Enable)
		if err != nil {
			logger.Error.Println("couldn't get password to validate")
			res.Code, res.Type, res.Msg = msg.GetByCode(cod)
			return c.Status(http.StatusAccepted).JSON(res)
		}
		result = valid
	}

	res.Data = result
	res.Code, res.Type, res.Msg = msg.GetByCode(29)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)

}

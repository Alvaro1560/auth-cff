package auth

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/ciphers"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/env"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/response"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/sendmail"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth/roles"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth/roles_password_policy"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth/users"
	"net/http"
	"time"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth/login"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/msgs"

	"github.com/jmoiron/sqlx"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	genTemplate "gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/template"
)

type Handler struct {
	DB   *sqlx.DB
	TxID string
}

func (h *Handler) LoginV3(c *fiber.Ctx) error {
	res := response.Model{Error: true}
	var msg msgs.Model
	m := LoginRequest{}
	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf(h.TxID, "no se pudo leer el Modelo User en login: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	m.RealIP = c.IP()
	serviceLogin := login.NewLoginService(h.DB, h.TxID)
	token, cod, err := serviceLogin.Login(m.ID, m.Username, m.Password, m.ClientID, m.HostName, m.RealIP)
	if err != nil {
		logger.Warning.Printf(h.TxID, "no se pudo leer el Modelo User en login: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(cod)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	res.Data = token
	if cod == 1001 {
		res.Code, res.Type, res.Msg = 1001, "success", "Procesado correctamente"
	} else {
		res.Code, res.Type, res.Msg = msg.GetByCode(cod)
	}

	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	res := response.Model{Error: true}
	var msg msgs.Model
	m := LoginRequest{}
	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf(h.TxID, "no se pudo leer el Modelo User en login: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	m.RealIP = c.IP()
	serviceLogin := login.NewLoginService(h.DB, h.TxID)
	token, cod, err := serviceLogin.Login(m.ID, m.Username, m.Password, m.ClientID, m.HostName, m.RealIP)
	if err != nil {
		logger.Warning.Printf(h.TxID, "no se pudo leer el Modelo User en login: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(cod)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	mr := LoginResponse{
		AccessToken:  token,
		RefreshToken: token,
	}
	res.Data = mr
	if cod == 1001 {
		res.Code, res.Type, res.Msg = 1001, "success", "Procesado correctamente"
	} else {
		res.Code, res.Type, res.Msg = msg.GetByCode(cod)
	}
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

func (h *Handler) ForgotPassword(c *fiber.Ctx) error {
	res := response.Model{Error: true}
	e := env.NewConfiguration()
	var msg msgs.Model
	var parameters = make(map[string]string, 0)
	m := ForgotPasswordRequest{}
	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf(h.TxID, "no se pudo leer el forgot password: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	repositoryUsers := users.FactoryStorage(h.DB, nil, h.TxID)
	serviceUsers := users.NewUserService(repositoryUsers, nil, h.TxID)

	user, cod, err := serviceUsers.GetUserByUsername(m.Username)
	if err != nil {
		logger.Error.Printf(h.TxID, "couldn't get user by username : %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(cod)
		res.Msg = err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if user == nil {
		logger.Error.Printf(h.TxID, "couldn't user with username %s", m.Username)
		res.Code, res.Type, res.Msg = msg.GetByCode(cod)
		res.Msg = err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if m.Email != user.EmailNotifications {
		logger.Error.Printf(h.TxID, "El correo de confirmación no es correcto", m.Email)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		res.Msg = err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	token, cod, err := login.GenerateJWT(models.User(*user))
	if err != nil {
		logger.Error.Printf(h.TxID, "no se pudo obtener modulos del usuario : ", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(cod)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	parameters["@token"] = e.App.UrlPortal + "/recoverypwd?access-token=" + token
	parameters["USER-NAME"] = user.Name + " " + user.LastName
	var tos []string
	tos = append(tos, user.EmailNotifications)

	logger.Trace.Println(tos)

	bodyCode, err := genTemplate.GenerateTemplateMail(parameters, e.Template.Recovery)
	if err != nil {
		logger.Error.Printf(h.TxID, "couldn't generate body in notification email")
		return err
	}

	emailCode := &sendmail.Model{}

	emailCode.From = "no-reply@e-capture.co"
	emailCode.To = tos
	emailCode.Subject = "Recuperación de cuenta"
	emailCode.Body = bodyCode

	err = emailCode.SendMail()
	if err != nil {
		logger.Error.Println(h.TxID, "error when execute send email: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(45)
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
		logger.Error.Printf(h.TxID, "no se pudo leer el forgot password: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	repositoryUsers := users.FactoryStorage(h.DB, nil, h.TxID)
	serviceUsers := users.NewUserService(repositoryUsers, nil, h.TxID)

	code, err := serviceUsers.ChangePassword(m.ID, m.Password, m.PasswordConfirm)
	if err != nil {
		logger.Error.Printf(h.TxID, "no se actualizar la contraseña: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
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
		logger.Error.Printf(h.TxID, "no se pudo leer el Modelo Password para validar politicas: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		res.Data = false
		return c.Status(http.StatusOK).JSON(res)
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
	passByte := ciphers.Decrypt(m.Password)
	if passByte == "" {
		logger.Error.Println("couldn't get password to validate")
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	m.Password = passByte
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

func (h *Handler) PasswordPolicyV2(c *fiber.Ctx) error {
	res := response.Model{Error: true}
	var msg msgs.Model
	m := PasswordPolicyRequestV2{}
	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf(h.TxID, "no se pudo leer el Modelo Password para validar politicas: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		res.Data = false
		return c.Status(http.StatusOK).JSON(res)
	}
	repositoryRPasswordPolicy := roles_password_policy.FactoryStorage(h.DB, nil, h.TxID)
	servicesRoles := roles_password_policy.NewRolesPasswordPolicyService(repositoryRPasswordPolicy, nil, h.TxID)
	srvRoles := roles.FactoryStorage(h.DB, nil, h.TxID)
	serviceRoles := roles.NewRoleService(srvRoles, nil, h.TxID)

	var rs []string

	rolesUser, code, err := serviceRoles.GetRolesByUserID(m.UserID)
	if err != nil {
		logger.Error.Println("couldn't get roles assigned to user, erro: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if rolesUser != nil && len(rolesUser) > 0 {
		for _, role := range rolesUser {
			rs = append(rs, role.ID)
		}
	} else {
		rs = append(rs, "50602690-B91F-4567-9A8D-A812B37A87BF")
	}

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
	passByte := ciphers.Decrypt(m.Password)
	if passByte == "" {
		logger.Error.Println("couldn't get password to validate")
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	repositoryUsers := users.FactoryStorage(h.DB, nil, h.TxID)
	serviceUsers := users.NewUserService(repositoryUsers, nil, h.TxID)
	m.Password = passByte
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

func (h *Handler) PasswordLife(c *fiber.Ctx) error {
	res := response.Model{Error: true}
	var msg msgs.Model
	m := PasswordLife{}
	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf(h.TxID, "no se pudo leer el Modelo Password para validar politicas: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		res.Data = false
		return c.Status(http.StatusOK).JSON(res)
	}
	repositoryRPasswordPolicy := roles_password_policy.FactoryStorage(h.DB, nil, h.TxID)
	servicesRoles := roles_password_policy.NewRolesPasswordPolicyService(repositoryRPasswordPolicy, nil, h.TxID)
	srvRoles := roles.FactoryStorage(h.DB, nil, h.TxID)
	serviceRoles := roles.NewRoleService(srvRoles, nil, h.TxID)

	var rs []string

	rolesUser, code, err := serviceRoles.GetRolesByUserID(m.UserID)
	if err != nil {
		logger.Error.Println("couldn't get roles assigned to user, erro: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if rolesUser != nil && len(rolesUser) > 0 {
		for _, role := range rolesUser {
			rs = append(rs, role.ID)
		}
	} else {
		rs = append(rs, "50602690-B91F-4567-9A8D-A812B37A87BF")
	}

	pp, err := servicesRoles.GetAllRolesPasswordPolicyByRolesIDs(rs)
	if err != nil {
		logger.Error.Println("couldn't get role to validate passwordPolicy")
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if pp == nil || len(pp) <= 0 {
		logger.Error.Println("don't exists role to validate passwordPolicy")
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	life := pp[0].DaysPassValid

	for _, policy := range pp {
		if life > policy.DaysPassValid {
			life = policy.DaysPassValid
		}
	}

	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	lastHistoryPassword, code, err := srvAuth.SrvPassword.GetLastPasswordByUserId(m.UserID)
	if err != nil {
		logger.Error.Println("couldn't get last History password, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if lastHistoryPassword != nil {
		currentDate := time.Now()
		dueDate := lastHistoryPassword.CreatedAt.AddDate(0, 0, life)

		if dueDate.Sub(currentDate).Hours() <= 0 {
			res.Error = false
			res.Code, res.Type, res.Msg = 43, "success", "La contraseña ha expirado o esta pronto a expirar, se recomienda cambiarla"
			return c.Status(http.StatusAccepted).JSON(res)
		}
	}

	res.Code, res.Type, res.Msg = msg.GetByCode(29)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

func (h *Handler) LoginGeneric(c *fiber.Ctx) error {
	res := response.Model{Error: true}
	var msg msgs.Model
	e := env.NewConfiguration()
	if !e.App.Autologin {
		res.Code, res.Type, res.Msg = msg.GetByCode(29)
		res.Error = false
		return c.Status(http.StatusOK).JSON(res)
	}
	key := Autologin{}

	err := c.BodyParser(&key)
	if err != nil {
		logger.Error.Printf(h.TxID, "no se pudo leer el Modelo User en login: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	if e.App.KeywordAutologin != key.Keyword {
		res.Code, res.Type, res.Msg = msg.GetByCode(29)
		res.Error = false
		return c.Status(http.StatusOK).JSON(res)
	}

	m := LoginRequest{
		ID:       "",
		Username: e.App.User,
		Password: e.App.Password,
		ClientID: 2,
		HostName: "",
		RealIP:   "",
	}

	m.RealIP = c.IP()
	serviceLogin := login.NewLoginService(h.DB, h.TxID)
	token, cod, err := serviceLogin.Login(m.ID, m.Username, m.Password, m.ClientID, m.HostName, m.RealIP)
	if err != nil {
		logger.Warning.Printf(h.TxID, "no se pudo leer el Modelo User en login: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(cod)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	res.Data = token
	if cod == 1001 {
		res.Code, res.Type, res.Msg = 1001, "success", "Procesado correctamente"
	} else {
		res.Code, res.Type, res.Msg = msg.GetByCode(cod)
	}
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

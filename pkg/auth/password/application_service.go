package password

import (
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
)

type PortsServerPassword interface {
	GetLastPasswordByUserId(userId string) (*Password, int, error)
}

type Service struct {
	repository ServicesPasswordRepository
	user       *models.User
	txID       string
}

func NewPasswordService(repository ServicesPasswordRepository, user *models.User, TxID string) PortsServerPassword {
	return Service{repository: repository, user: user, txID: TxID}
}

func (s Service) GetLastPasswordByUserId(userId string) (*Password, int, error) {
	m, err := s.repository.GetByUserID(userId)
	if err != nil {
		logger.Error.Println("No se pudo obtner el ultimo historial de contraseña, err: %v", err)
		return nil, 22, err
	}

	return m, 29, nil
}

/*
func (s Service) PasswordManagementSelf(c echo.Context) error {
	rm := response.Model{}
	m := &Password{}
	u := c.Get("user").(Model)

	err := c.Bind(m)
	if err != nil {
		logger.Error.Printf("la estructura del objeto Password no es correcta: %v", err)
		rm.Set(true, nil, 1)
		return c.JSON(http.StatusAccepted, rm)
	}

	if m.New != m.ConfirmNew {
		rm.Set(true, nil, 74)
		return c.JSON(http.StatusAccepted, rm)
	}

	ru, err := u.GetByEmail(u.Email)
	if err == sql.ErrNoRows {
		rm.Set(true, nil, 10)
		return c.JSON(http.StatusAccepted, rm)
	}
	if err != nil {
		logger.Error.Printf("no se pudo consultar el usuario por el email para cambiar su propia contraseña: %s, %v", u.Email, err)
		rm.Set(true, nil, 70)
		return c.JSON(http.StatusAccepted, rm)
	}

	if u.Email == m.New {
		rm.Set(true, nil, 84)
		return c.JSON(http.StatusAccepted, rm)
	}

	if !ru.ComparePassword(m.Actual) {
		rm.Set(true, nil, 10)
		return c.JSON(http.StatusAccepted, rm)
	}

	// Valida si existe la clave en listas negras
	blp := black_list_password.Model{}
	_, err = blp.GetByPassword(m.New)
	if err != sql.ErrNoRows {
		rm.Set(true, nil, 75)
		return c.JSON(http.StatusAccepted, rm)
	}

	// Consulta la política de passwords del role
	pp := &password_policy.Model{}
	pp, err = pp.GetByRoleID(int64(u.RoleID))

	if pp.Enable {
		if passwordIsInHistory(pp, u.ID, m.New) {
			rm.Set(true, nil, 76)
			return c.JSON(http.StatusAccepted, rm)
		}

		if len(m.New) > pp.MaxLength || len(m.New) < pp.MinLength {
			rm.Set(true, nil, 77)
			return c.JSON(http.StatusAccepted, rm)
		}

		// Se evaluan las expresiones regulares de la política
		er := regexp.MustCompile(fmt.Sprintf("((.*[a-zA-Z]){%d})", pp.Alpha))
		if !er.Match([]byte(m.New)) {
			rm.Set(true, nil, 78)
			return c.JSON(http.StatusAccepted, rm)
		}

		er = regexp.MustCompile(fmt.Sprintf("((.*[0-9]){%d})", pp.Digits))
		if !er.Match([]byte(m.New)) {
			rm.Set(true, nil, 79)
			return c.JSON(http.StatusAccepted, rm)
		}

		er = regexp.MustCompile(fmt.Sprintf("((.*[a-z]){%d})", pp.LowerCase))
		if !er.Match([]byte(m.New)) {
			rm.Set(true, nil, 80)
			return c.JSON(http.StatusAccepted, rm)
		}

		er = regexp.MustCompile(fmt.Sprintf("((.*[A-Z]){%d})", pp.UpperCase))
		if !er.Match([]byte(m.New)) {
			rm.Set(true, nil, 81)
			return c.JSON(http.StatusAccepted, rm)
		}

		er = regexp.MustCompile("((.*(\\`|\\~|\\!|\\@|\\#|\\$|\\%|\\^|\\&|\\*|\\(|\\)|\\+|\\=|\\[|\\{|\\]|\\}|\\||\\'|\\<|\\,|\\.|\\>|\\?|\\/|\"|\\;|\\:))){" + strconv.Itoa(pp.Special) + "}")
		if !er.Match([]byte(m.New)) {
			rm.Set(true, nil, 82)
			return c.JSON(http.StatusAccepted, rm)
		}

	}

	ru.Password = m.New
	err = ru.ChangeNewPassword()
	if err != nil {
		logger.Error.Printf("no se pudo actualizar el password del usuario: %d, %s, %v", u.ID, u.Email, err)
		rm.Set(true, nil, 85)
		return c.JSON(http.StatusAccepted, rm)
	}

	rm.Set(false, nil, 84)
	return c.JSON(http.StatusOK, rm)
}

// passwordIsInHistory valida que la contraseña
// no se encuentre en el histórico configurado
// True: El password se encuentra en el histórico
// False: El password no se encuentra en el histórico o no se valida esta opción
func (s Service) passwordIsInHistory(pp *password_policy.Model, userID int64, pwd string) bool {
	// Si la política dice que la cantidad de veces que
	// una contraseña puede estar en el histórico es de 0 o menor
	// devuelve false
	if pp.Times <= 0 {
		return false
	}

	ph := password_history.Model{}
	ph.UserID = userID
	phs, err := ph.GetLastNByUserID(pp.Times)
	if err != nil {
		errStr := fmt.Errorf("no se pudo consultar los registros de password_history")
		logger.Error.Print(errStr, err)
		return false
	}

	for _, p := range phs {
		if p.Passwd == pwd {
			return true
		}
	}

	return false
}

func (s Service) PasswordManagementAdmin(c echo.Context) error {

	rm := response.Model{}
	u := &Model{}
	type request struct {
		Email string `json:"email"`
	}
	req := &request{}

	err := c.Bind(req)
	if err != nil {
		logger.Error.Printf("la estructura del objeto no es correcta: %v", err)
		rm.Set(true, nil, 1)
		return c.JSON(http.StatusAccepted, rm)
	}

	u, err = u.GetByEmail(req.Email)
	if err != nil {
		logger.Error.Printf("no se encontró usuario con el email: %s, %v", req.Email, err)
		rm.Set(true, nil, 22)
		return c.JSON(http.StatusAccepted, rm)
	}

	rp, err := uuid.NewV4()
	if err != nil {
		logger.Error.Printf("no se pudo generar el UUID de la clave: %s", err)
		rm.Set(true, nil, 70)
		return c.JSON(http.StatusAccepted, rm)
	}

	u.Password = rp.String()
	err = u.ChangeNewPassword()
	if err != nil {
		logger.Error.Printf("no se pudo asignar la contraseña al usuario: %d, %s: %v", u.ID, u.Email, err)
		rm.Set(true, nil, 70)
		return c.JSON(http.StatusAccepted, rm)
	}

	// to Es el correo a donde llegará el password
	cnfg := configuration.FromFile()
	var to string

	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(u.Email) {
		to = u.Email
	} else {
		to = cnfg.SupervisorEmail
	}

	sm := sendmail.Model{
		From:    cnfg.MailEmail,
		To:      []string{to},
		CC:      nil,
		Subject: "Reestablecer contraseña",
		Body:    fmt.Sprintf(`<h1>Reestablecer contraseña</h1><div><p>La contraseña para el usuario: %s</p><p style="font-weight:bold;">%s</p></div>`, u.Email, rp.String()),
	}

	err = sm.SendMail()
	if err != nil {
		logger.Error.Printf("no se pudo enviar el correo de reestablecer contraseña: %s, %v", u.Email, err)
		rm.Set(true, nil, 86)
		return c.JSON(http.StatusAccepted, rm)
	}

	err = u.SetChangePassword()
	if err != nil {
		logger.Error.Printf("no se pudo actualizar el campo change_password del usuario: %s, %v", u.Email, err)
	}

	rm.Set(false, nil, 29)
	return c.JSON(http.StatusAccepted, rm)
}

func (s Service) ChangePasswordAdmin(c echo.Context) error {

	rm := response.Model{}
	u := &Model{}
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	req := &request{}

	err := c.Bind(req)
	if err != nil {
		logger.Error.Printf("la estructura del objeto no es correcta: %v", err)
		rm.Set(true, nil, 1)
		return c.JSON(http.StatusAccepted, rm)
	}

	u, err = u.GetByEmail(req.Email)
	if err != nil {
		logger.Error.Printf("no se encontró usuario con el email: %s, %v", req.Email, err)
		rm.Set(true, nil, 22)
		return c.JSON(http.StatusAccepted, rm)
	}

	rp, err := uuid.NewV4()
	if err != nil {
		logger.Error.Printf("no se pudo generar el UUID de la clave: %s", err)
		rm.Set(true, nil, 70)
		return c.JSON(http.StatusAccepted, rm)
	}

	u.Password = rp.String()
	err = u.AdminChangePassword(req.Password)
	if err != nil {
		logger.Error.Printf("no se pudo asignar la contraseña al usuario: %d, %s: %v", u.ID, u.Email, err)
		rm.Set(true, nil, 70)
		return c.JSON(http.StatusAccepted, rm)
	}

	err = u.SetChangePassword()
	if err != nil {
		logger.Error.Printf("no se pudo actualizar el campo change_password del usuario: %s, %v", u.Email, err)
	}

	rm.Set(false, nil, 29)
	return c.JSON(http.StatusAccepted, rm)
}

func (s Service) UserDisableEnableAdmin(c echo.Context) error {
	rm := response.Model{}
	u := &Model{}
	type request struct {
		ID     int64 `json:"id"`
		Enable bool  `json:"enable"`
	}
	req := &request{}

	err := c.Bind(req)
	if err != nil {
		logger.Error.Printf("la estructura del objeto no es correcta: %v", err)
		rm.Set(true, nil, 1)
		return c.JSON(http.StatusAccepted, rm)
	}

	u, err = u.GetByID(req.ID)
	if err != nil {
		logger.Error.Printf("no se pudo consultar el usuario con el ID: %d, %v", req.ID, err)
		rm.Set(true, nil, 70)
		return c.JSON(http.StatusAccepted, rm)
	}

	if req.Enable {
		err = u.Enable()
		if err != nil {
			logger.Error.Printf("no se pudo habilitar el usuario con el ID: %d, %v", req.ID, err)
			rm.Set(true, nil, 70)
			return c.JSON(http.StatusAccepted, rm)
		}
	} else {
		err = u.Disable()
		if err != nil {
			logger.Error.Printf("no se pudo deshabilitar el usuario con el ID: %d, %v", req.ID, err)
			rm.Set(true, nil, 70)
			return c.JSON(http.StatusAccepted, rm)
		}
	}

	rm.Set(false, nil, 29)
	return c.JSON(http.StatusAccepted, rm)
}
*/

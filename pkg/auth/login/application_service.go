package login

import (
	"fmt"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/ciphers"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/ldap"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/parameters"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/password"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/transact"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth/roles"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth/roles_password_policy"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth/users"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/config/customers_projects"
)

type Service struct {
	DB   *sqlx.DB
	TxID string
}

func NewLoginService(db *sqlx.DB, TxID string) Service {
	return Service{DB: db, TxID: TxID}
}

func (s *Service) Login(id, Username, Password string, ClientID int, HostName, RealIP string) (string, int, error) {
	var token string
	m := NewLogin(id, Username, Password, ClientID, HostName, RealIP)
	if m.Username == "" {
		m.Username = m.ID
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.TxID, " - don't meet validations:", err)
		return token, 15, err
	}

	if m.ClientID == 2 {
		ciphers.Encrypt(Password)
		//ITC
	} else if m.ClientID == 9925 {
		infoByte := ciphers.Decrypt(Password)
		if infoByte == "" {
			logger.Error.Println(s.TxID, " - don't meet validations:")
			return token, 15, fmt.Errorf(s.TxID, " - don't meet validations:")
		}
		m.Password = infoByte

		userByte := ciphers.Decrypt(id)
		if userByte == "" {
			logger.Error.Println(s.TxID, " - don't meet validations:")
			return token, 15, fmt.Errorf(s.TxID, " - don't meet validations:")
		}
		m.ID = userByte
		m.Username = m.ID
		//WebClient WebConfig
	} else if m.ClientID == 9926 || m.ClientID == 9927 {
		infoByte := ciphers.DecryptKeyTemp(Password)
		if infoByte == "" {
			logger.Error.Println(s.TxID, " - don't meet validations:")
			return token, 15, fmt.Errorf(s.TxID, " - don't meet validations:")
		}
		m.Password = infoByte

		userByte := ciphers.DecryptKeyTemp(id)
		if userByte == "" {
			logger.Error.Println(s.TxID, " - don't meet validations:")
			return token, 15, fmt.Errorf(s.TxID, " - don't meet validations:")
		}
		m.ID = userByte
		m.Username = m.ID

	} else {
		logger.Error.Println(s.TxID, " - client not configured: ")
		return token, 70, fmt.Errorf(" - client not configured: ")
	}

	usr, cod, err := s.getUserByUsername(m.Username)
	if err != nil {
		transact.RegisterLogUsr("user-not-found", HostName, RealIP, RealIP, Username)
		logger.Error.Println(s.TxID, " - couldn't get user by id", err)
		return token, cod, err
	}

	if parameters.GetParameter("APP_VALIDATE_IP") == "t" {
		howManyUnAuthorizedIP, err := s.countUnAuthorizedIP(m.RealIP)
		if err != nil {
			logger.Error.Printf(s.TxID, " - No se pudo obtener listado de IPs a validar: %v ", err)
			return token, 70, err
		}
		if howManyUnAuthorizedIP == 0 {
			transact.RegisterLogUsr("ip-not-allowed", HostName, RealIP, RealIP, Username)
			logger.Warning.Printf(s.TxID, " - Intento de ingreso desde una IP no autorizada: Usuario: %s, IP: %s", m.ID, m.RealIP)
			return token, 91, err
		}
	}

	if parameters.GetParameter("LDAP") == "t" {
		var cod int
		var groupValid bool
		usr.Roles, groupValid, cod, err = s.ldapAuthentication(usr.ID, usr.Password, usr.Roles)
		if err != nil {
			logger.Error.Printf(s.TxID, " - Authentication Ldap role: %v", err)
			return token, cod, err
		}
		if !groupValid {
			transact.RegisterLogUsr("unauthorized-user", HostName, RealIP, RealIP, usr.ID)
			logger.Error.Printf(s.TxID, " - el usuario no pertenece a ningun role: %v", err)
			return token, cod, err
		}
	}
	usr.Roles, cod, err = s.getRolesByUserID(strings.ToLower(usr.ID))
	if err != nil {
		transact.RegisterLogUsr("user-has-not-role", HostName, RealIP, RealIP, usr.ID)
		logger.Warning.Printf(s.TxID, " - el usuario no tiene asignado roles: %s, IP: %s", m.ID, m.RealIP, err)
		return token, cod, nil
	}
	usr.Projects, cod, err = s.getProjectByRoles(usr.Roles)
	if err != nil {
		transact.RegisterLogUsr("user-has-not-project", HostName, RealIP, RealIP, usr.ID)
		logger.Warning.Printf(s.TxID, " - el usuario no tiene asignado proyecto: %s, IP: %s", m.ID, m.RealIP, err)
		return token, cod, nil
	}
	politics, cod, err := s.validatePasswordPolicies(usr.Roles)
	if err != nil {
		transact.RegisterLogUsr("user-has-not-password-politics", HostName, RealIP, RealIP, usr.ID)
		logger.Warning.Printf(s.TxID, " - No se pudo obtener las politicas de las contraseñas: %s, IP: %s", m.ID, m.RealIP)
		return token, cod, err
	}

	if politics == nil {
		transact.RegisterLogUsr("user-has-not-password-politics", HostName, RealIP, RealIP, usr.ID)
		logger.Warning.Printf(s.TxID, " - El usuario no tiene politicas de seguridad itegradas, contactese con su administrador: %s, IP: %s", m.ID, m.RealIP)
		return token, cod, fmt.Errorf("el usuario no tiene politicas de seguridad itegradas, contactese con su administrador")
	}

	if !password.Compare(usr.ID, usr.Password, m.Password) {
		transact.RegisterLogUsr("login-filed", HostName, RealIP, RealIP, usr.ID)

		failedAttempts := usr.FailedAttempts + 1
		_, err = s.registerFailedAttempts(usr.ID, failedAttempts)
		if err == nil {
			attempts := politics[0].FailedAttempts
			for _, politic := range politics {
				if attempts > politic.FailedAttempts {
					attempts = politic.FailedAttempts
				}
			}

			currentDate := time.Now()
			if failedAttempts >= attempts {
				usr.Status = 16
				usr.BlockDate = &currentDate
				transact.RegisterLogUsr("user-blocked", HostName, RealIP, RealIP, usr.ID)
				s.blockUser(usr.ID)
			}
		}

		return token, 10, fmt.Errorf("usuario o contraseña incorrecta")
	}
	/*
		if policies.IsDateDisallowed {
			transact.RegisterLogUsr("login-IsDateDisallowed")
			return token, 90, nil
		}

		if policies.LoginAllowed < len(user.LoggedUsers) {
			transact.RegisterLogUsr("allowed-connection-exceeded")
			return token, 88, nil
		}

		if policies.IsEnablePasswordPolicy {
			user.userUnblock(policies.TimeUnlock)
			user.mustChangePassword(policies.ValidityPassChange)
			user.changePasswordDaysLeft(policies.ValidityPassChange)

		}*/

	if usr.Status == 16 {
		timeUnlock := politics[0].TimeUnlock
		for _, politic := range politics {
			if timeUnlock < politic.TimeUnlock {
				timeUnlock = politic.TimeUnlock
			}
		}

		dateUnlock := usr.BlockDate.AddDate(0, 0, timeUnlock)
		currentDate := time.Now()

		if dateUnlock.Sub(currentDate).Hours() > 0 {
			return token, 69, nil
		}
		transact.RegisterLogUsr("user-unlocked", HostName, RealIP, RealIP, usr.ID)
		s.unBlockUser(usr.ID)
		s.registerFailedAttempts(usr.ID, 0)
	}

	cod, err = s.registerLoggedUser(usr.ID, m.RealIP, m.HostName)
	if err != nil {
		transact.RegisterLogUsr("register-login-filed", HostName, RealIP, RealIP, usr.ID)
		logger.Warning.Printf(s.TxID, " - no se pudo registrar la trazabilidad: %s, IP: %s", m.ID, m.RealIP)
		return token, cod, nil
	}
	transact.RegisterLogUsr("success-login", HostName, RealIP, RealIP, usr.ID)
	usr.SessionID = uuid.New().String()
	usr.ClientID = m.ClientID
	usr.RealIP = m.RealIP
	usr.Password = ""
	usr.Colors.Primary, usr.Colors.Secondary, usr.Colors.Tertiary = s.getColors("")

	token, cod, err = GenerateJWT(*usr)
	if err != nil {
		logger.Error.Printf(s.TxID, "no se pudo obtener modulos del usuario : ", err)
		return "", cod, err
	}

	code := 29

	for _, politic := range politics {
		if politic.Required2fa {
			code = 1001
		}
	}

	return token, code, nil
}

func (s *Service) getUserByUsername(username string) (*models.User, int, error) {

	repositoryUsers := users.FactoryStorage(s.DB, nil, s.TxID)
	serviceUser := users.NewUserService(repositoryUsers, nil, s.TxID)
	user, _, err := serviceUser.GetUserByUsername(username)
	if err != nil {
		logger.Error.Println("couldn't get user by id", err)
		return nil, 10, err
	}
	if user == nil {
		logger.Error.Println("couldn't get user by id", err)
		return nil, 10, fmt.Errorf("couldn't get user by id")
	}
	usr := models.User(*user)
	return &usr, 29, nil

}

// TODO implement countUnAuthorizedIP
func (s *Service) countUnAuthorizedIP(ip string) (int, error) {

	return 29, nil
}

func (s *Service) ldapAuthentication(id, password string, roles []*string) ([]*string, bool, int, error) {
	ldapSSO := parameters.GetParameter("LDAP_SSO")
	var groupValid bool
	var username, bindusername, bindpassword string
	if strings.ToLower(ldapSSO) == "t" {
		transact.RegisterLogUsr("ldap-single-sign-on", "HostName", "IpRequest", "IpRemote", username)
		username = strings.Split(id, "@")[0]
		bindusername = parameters.GetParameter("LDAP_USERNAME_SSO")
		bindpassword = parameters.GetParameter("LDAP_USERNAME_SSO")
	} else {
		transact.RegisterLogUsr("ldap-authentication", "HostName", "IpRequest", "IpRemote", username)
		username = strings.Split(id, "@")[0]
		bindusername = strings.Split(id, "@")[0]
		bindpassword = password
	}
	groups, err := ldap.Authentication(username, bindusername, bindpassword)
	if err != nil {
		transact.RegisterLogUsr("unauthorized-ldap", "HostName", "IpRequest", "IpRemote", username)
		logger.Error.Printf(s.TxID, "no fue posible consultar la información de grupos en LDAP: %v", err)
		return nil, false, 103, err
	}
	var rls []*string
	for _, g := range groups {
		for _, r := range roles {
			if strings.ToLower(g) == strings.ToLower(*r) {
				groupValid = true
				rls = append(rls, r)

			}
		}
	}
	return rls, groupValid, 29, nil
}

func (s *Service) validatePasswordPolicies(roles []*string) ([]*roles_password_policy.RolesPasswordPolicy, int, error) {
	var rs []string
	for _, r := range roles {
		rs = append(rs, *r)
	}
	repositoryPwd := roles_password_policy.FactoryStorage(s.DB, nil, s.TxID)
	servicePwd := roles_password_policy.NewRolesPasswordPolicyService(repositoryPwd, nil, s.TxID)

	politics, err := servicePwd.GetAllRolesPasswordPolicyByRolesIDs(rs)
	if err != nil {
		logger.Error.Println("couldn't get user by id", err)
		return nil, 22, err
	}

	return politics, 29, nil
}

// TODO implement registerLoggedUser
func (s *Service) registerLoggedUser(id, realIP, hostName string) (int, error) {

	return 25, nil
}

// TODO implement getColors
func (s *Service) getColors(project string) (string, string, string) {

	return "#353A48", "#039be5", "#262933"
}

func (s *Service) getProjectByRoles(roles []*string) ([]*string, int, error) {
	var rs []string
	for _, r := range roles {
		rs = append(rs, *r)
	}
	repositoryProjects := customers_projects.FactoryStorage(s.DB, nil, s.TxID)
	serviceRoles := customers_projects.NewProjectService(repositoryProjects, nil, s.TxID)
	projects, err := serviceRoles.GetProjectByRoles(rs)
	if err != nil {
		logger.Error.Println("couldn't get roles by user id")
		return nil, 104, err
	}

	return projects, 29, nil
}

func (s *Service) getRolesByUserID(id string) ([]*string, int, error) {
	var UserRoles []*string
	repositoryRoles := roles.FactoryStorage(s.DB, nil, s.TxID)
	serviceRoles := roles.NewRoleService(repositoryRoles, nil, s.TxID)
	roles, _, err := serviceRoles.GetRolesByUserID(id)
	if err != nil {
		logger.Error.Println("couldn't get roles by user id")
		return nil, 104, err
	}
	for _, r := range roles {
		UserRoles = append(UserRoles, &r.ID)
	}
	return UserRoles, 29, nil
}

func (s *Service) registerFailedAttempts(userId string, failedAttempts int) (int, error) {

	repositoryUsers := users.FactoryStorage(s.DB, nil, s.TxID)
	serviceUser := users.NewUserService(repositoryUsers, nil, s.TxID)
	err := serviceUser.UpdateFailedAttempts(userId, failedAttempts)
	if err != nil {
		logger.Error.Println("couldn't get user by id", err)
		return 22, err
	}

	return 29, nil
}

func (s *Service) blockUser(userId string) (int, error) {

	repositoryUsers := users.FactoryStorage(s.DB, nil, s.TxID)
	serviceUser := users.NewUserService(repositoryUsers, nil, s.TxID)
	err := serviceUser.BlockUser(userId)
	if err != nil {
		logger.Error.Println("couldn't block user by id", err)
		return 22, err
	}

	return 29, nil
}

func (s *Service) unBlockUser(userId string) (int, error) {

	repositoryUsers := users.FactoryStorage(s.DB, nil, s.TxID)
	serviceUser := users.NewUserService(repositoryUsers, nil, s.TxID)
	err := serviceUser.UnblockUser(userId)
	if err != nil {
		logger.Error.Println("couldn't un block user by id", err)
		return 22, err
	}

	return 29, nil
}

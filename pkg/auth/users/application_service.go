package users

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/password"

	"github.com/asaskevich/govalidator"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
)

type Service struct {
	repository ServicesUserRepository
	user       *models.User
	txID       string
}

func NewUserService(repository ServicesUserRepository, user *models.User, TxID string) Service {
	return Service{repository: repository, user: user, txID: TxID}
}

func (s Service) CreateUser(id string, Username string, Name string, LastName string, Password string, EmailNotifications string, IdentificationNumber string, IdentificationType string) (*User, int, error) {
	var changePass, isBlock, isDisabled bool
	m := NewUser(id, Username, Name, LastName, EmailNotifications, IdentificationNumber, IdentificationType)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	m.Password = password.Encrypt(Password)

	if m.ChangePassword == nil {
		m.ChangePassword = &changePass
	}
	m.IsBlock = &isBlock
	m.IsDisabled = &isDisabled
	if err := s.repository.Create(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't create User :", err)
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		return m, 3, err
	}
	return m, 29, nil
}
func (s Service) UpdateUser(id string, User string, Name string, LastName string, Password string, EmailNotifications string, IdentificationNumber string, IdentificationType string) (*User, int, error) {
	m := NewUser(id, User, Name, LastName, EmailNotifications, IdentificationNumber, IdentificationType)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.Update(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update User :", err)
		return m, 18, err
	}
	return m, 29, nil
}
func (s Service) DeleteUser(id string) (int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return 15, fmt.Errorf("id isn't uuid")
	}

	if err := s.repository.Delete(id); err != nil {
		if err.Error() == "ecatch:108" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update row:", err)
		return 20, err
	}
	return 28, nil
}
func (s Service) GetUserByID(id string) (*User, int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return nil, 15, fmt.Errorf("id isn't uuid")
	}
	m, err := s.repository.GetByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}
func (s Service) GetUserByUsername(username string) (*User, int, error) {
	m, err := s.repository.GetByUsername(username)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}
func (s Service) GetAllUser() ([]*User, error) {
	return s.repository.GetAll()
}
func (s Service) GetUsersByIDs(ids []string) ([]*User, error) {
	return s.repository.GetUsersByIDs(ids)
}
func (s Service) BlockUser(id string) error {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return fmt.Errorf("id isn't uuid")
	}
	err := s.repository.BlockUser(id)
	if err != nil {
		logger.Error.Printf(s.txID, "couldn't Block User: %v", err)
		return err
	}
	//myMail := &sendmail.Model{}
	//go myMail.SendMail("send_mail.gohtml", s.user.ID, s.user.EmailNotifications, fmt.Sprintf("Usuario %s Bloqueado ", s.user.Name))

	return nil
}
func (s Service) UnblockUser(id string) error {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return fmt.Errorf(s.txID, "id isn't uuid")
	}
	err := s.repository.UnblockUser(id)
	if err != nil {
		logger.Error.Printf(s.txID, "couldn't Unblock User: %v", err)
		return err
	}
	//myMail := &sendmail.Model{}
	//go myMail.SendMailNotification("send_mail.gohtml", s.user.ID, s.user.EmailNotifications, fmt.Sprintf("Usuario %s Desbloqueado", s.user.Name))

	return nil
}
func (s Service) LogoutUser(id string) (int, error) {
	if !govalidator.IsUUID(strings.ToLower(id)) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return 15, fmt.Errorf("id isn't uuid")
	}

	if err := s.repository.LogoutUser(id); err != nil {
		if err.Error() == "ecatch:31" {
			return 31, nil
		}
		logger.Error.Println(s.txID, " - couldn't LogoutUser row:", err)
		return 31, err
	}
	return 28, nil
}
func (s Service) ChangePassword(id string, pass string, passConfirm string) (int, error) {
	if !govalidator.IsUUID(strings.ToLower(id)) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return 15, fmt.Errorf("id isn't uuid")
	}

	if pass == passConfirm {
		passwordCipher := password.Encrypt(pass)
		if err := s.repository.ChangePassword(id, passwordCipher); err != nil {
			if err.Error() == "ecatch:31" {
				return 31, nil
			}
			logger.Error.Println(s.txID, " - couldn't ChangePassword row:", err)
			return 31, err
		}
	} else {
		logger.Error.Println(s.txID, " - password does not match")
		return 83, fmt.Errorf("password does not match")
	}

	return 28, nil
}
func (s Service) UpdatePasswordByUser(id string, pass string, passConfirm string, passwordOld string) (int, error) {
	if !govalidator.IsUUID(strings.ToLower(id)) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return 15, fmt.Errorf("id isn't uuid")
	}
	usr, err := s.repository.GetByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't UpdatePasswordByUser row:", err)
		return 10, nil
	}

	if pass == passConfirm {
		passwordCipher := password.Encrypt(pass)
		if !password.Compare(id, usr.Password, passwordOld) {
			return 10, fmt.Errorf("Incorrect user or password")
		}
		if err := s.repository.UpdatePasswordByUser(id, passwordCipher); err != nil {
			if err.Error() == "ecatch:71" {
				return 72, nil
			}
			logger.Error.Println(s.txID, " - couldn't UpdatePasswordByUser row:", err)
			return 72, err
		}
	} else {
		logger.Error.Println(s.txID, " - password does not match")
		return 83, fmt.Errorf("password does not match")
	}

	return 28, nil
}
func (s Service) UpdatePasswordByAdministrator(id string, pass string, passConfirm string) (int, error) {
	if !govalidator.IsUUID(strings.ToLower(id)) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return 15, fmt.Errorf("id isn't uuid")
	}
	if pass == passConfirm {
		passwordCipher := password.Encrypt(pass)
		if err := s.repository.UpdatePasswordByUser(id, passwordCipher); err != nil {
			if err.Error() == "ecatch:71" {
				return 72, nil
			}
			logger.Error.Println(s.txID, " - couldn't UpdatePasswordByUser row:", err)
			return 72, err
		}
	} else {
		logger.Error.Println(s.txID, " - password does not match")
		return 83, fmt.Errorf("password does not match")
	}

	return 28, nil
}
func (s Service) GetUserByUsernameAndIdentificationNumber(username string, identificationNumber string) (*User, int, error) {
	m, err := s.repository.GetByUsernameAndIdentificationNumber(username, identificationNumber)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't GetUserByUsernameAndIdentificationNumber row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s Service) ValidatePasswordPolicy(password string, maxLength, minLength, alpha, digits, special,
	upperCase, lowerCase int, enable bool) (bool, int, error) {
	if !enable {
		return true, 29, nil
	}
	var responseValidate bool

	if len(password) > maxLength || len(password) < minLength {
		return responseValidate, 77, fmt.Errorf("longitud")
	}
	er := regexp.MustCompile(fmt.Sprintf("((.*[a-zA-Z]){%d})", alpha))
	if !er.Match([]byte(password)) {
		return responseValidate, 78, fmt.Errorf("alpha")
	}
	er = regexp.MustCompile(fmt.Sprintf("((.*[0-9]){%d})", digits))
	if !er.Match([]byte(password)) {
		return responseValidate, 79, fmt.Errorf("digits")
	}
	er = regexp.MustCompile(fmt.Sprintf("((.*[a-z]){%d})", upperCase))
	if !er.Match([]byte(password)) {
		return responseValidate, 89, fmt.Errorf("lowercase")
	}
	er = regexp.MustCompile(fmt.Sprintf("((.*[A-Z]){%d})", lowerCase))
	if !er.Match([]byte(password)) {
		return responseValidate, 81, fmt.Errorf("uppercase")
	}
	er = regexp.MustCompile("((.*(\\-|\\_|\\`|\\~|\\!|\\@|\\#|\\$|\\%|\\^|\\&|\\*|\\(|\\)|\\+|\\=|\\[|\\{|\\]|\\}|\\||\\'|\\<|\\,|\\.|\\>|\\?|\\/|\"|\\;|\\:))){" + strconv.Itoa(special) + "}")
	if !er.Match([]byte(password)) {
		return responseValidate, 82, fmt.Errorf("special")
	}
	responseValidate = true
	return responseValidate, 29, nil

}

func (s Service) DeleteUserPasswordHistory(id string) (int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return 15, fmt.Errorf("id isn't uuid")
	}

	if err := s.repository.DeleteUserPasswordHistory(id); err != nil {
		if err.Error() == "ecatch:108" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't DeleteUserPasswordHistory row:", err)
		return 20, err
	}
	return 28, nil
}

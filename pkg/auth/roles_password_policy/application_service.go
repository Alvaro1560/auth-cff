package roles_password_policy

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
)

type Service struct {
	repository ServicesRolesPasswordPolicyRepository
	user       *models.User
	txID       string
}

func NewRolesPasswordPolicyService(repository ServicesRolesPasswordPolicyRepository, user *models.User, TxID string) Service {
	return Service{repository: repository, user: user, txID: TxID}
}

func (s Service) CreateRolesPasswordPolicy(id string, RoleId string, DaysPassValid int, MaxLength int, MinLength int, StorePassNotRepeated int, FailedAttempts int, TimeUnlock int, Alpha int, Digits int, Special int, UpperCase int, LowerCase int, Enable bool, InactivityTime int, Timeout int) (*RolesPasswordPolicy, int, error) {
	m := NewRolesPasswordPolicy(id, RoleId, DaysPassValid, MaxLength, MinLength, StorePassNotRepeated, FailedAttempts, TimeUnlock, Alpha, Digits, Special, UpperCase, LowerCase, Enable, InactivityTime, Timeout)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.Create(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't create RolesPasswordPolicy :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s Service) UpdateRolesPasswordPolicy(id string, RoleId string, DaysPassValid int, MaxLength int, MinLength int, StorePassNotRepeated int, FailedAttempts int, TimeUnlock int, Alpha int, Digits int, Special int, UpperCase int, LowerCase int, Enable bool, InactivityTime int, Timeout int) (*RolesPasswordPolicy, int, error) {
	m := NewRolesPasswordPolicy(id, RoleId, DaysPassValid, MaxLength, MinLength, StorePassNotRepeated, FailedAttempts, TimeUnlock, Alpha, Digits, Special, UpperCase, LowerCase, Enable, InactivityTime, Timeout)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.Update(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update RolesPasswordPolicy :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s Service) DeleteRolesPasswordPolicy(id string) (int, error) {
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

func (s Service) GetRolesPasswordPolicyByID(id string) (*RolesPasswordPolicy, int, error) {
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

func (s Service) GetAllRolesPasswordPolicy() ([]*RolesPasswordPolicy, error) {
	return s.repository.GetAll()
}
func (s Service) GetAllRolesPasswordPolicyByRolesIDs(RolesIDs []string) ([]*RolesPasswordPolicy, error) {
	return s.repository.GetAllRolesPasswordPolicyByRolesIDs(RolesIDs)
}

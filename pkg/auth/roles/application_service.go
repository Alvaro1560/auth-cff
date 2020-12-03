package roles

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
)

type Service struct {
	repository ServicesRoleRepository
	user       *models.User
	txID       string
}

func NewRoleService(repository ServicesRoleRepository, user *models.User, TxID string) Service {
	return Service{repository: repository, user: user, txID: TxID}
}
func (s Service) CreateRole(id string, Name string, Description string, SessionsAllowed int) (*Role, int, error) {
	m := NewRole(id, Name, Description, SessionsAllowed)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.Create(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't create Role :", err)
		return m, 3, err
	}
	return m, 29, nil
}
func (s Service) UpdateRole(id string, Name string, Description string, SessionsAllowed int) (*Role, int, error) {
	m := NewRole(id, Name, Description, SessionsAllowed)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.Update(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update Role :", err)
		return m, 18, err
	}
	return m, 29, nil
}
func (s Service) DeleteRole(id string) (int, error) {
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
func (s Service) GetRoleByID(id string) (*Role, int, error) {
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
func (s Service) GetAllRole() ([]*Role, error) {
	return s.repository.GetAll()
}
func (s Service) GetRolesByUserID(userID string) ([]*Role, int, error) {
	ms, err := s.repository.GetByUserID(userID)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't GetRolesByUserID row:", err)
		return nil, 22, err
	}
	return ms, 29, nil
}
func (s Service) GetRolesByProcessIDs(ProcessIDs []string) ([]*Role, error) {
	return s.repository.GetRolesByProcessIDs(ProcessIDs)
}
func (s Service) GetRolesByIDs(IDs []string) ([]*Role, error) {
	return s.repository.GetRolesByIDs(IDs)
}
func (s Service) GetRolesByUserIDs(userIDs []string) ([]*Role, int, error) {
	ms, err := s.repository.GetByUserIDs(userIDs)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't GetRolesByUserID row:", err)
		return nil, 22, err
	}
	return ms, 29, nil
}

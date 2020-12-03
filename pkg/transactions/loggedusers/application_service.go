package loggedusers

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
)

type Service struct {
	repository ServicesTxLoggedUserRepository
	user       *models.User
	txID       string
}

func NewTxLoggedUserService(repository ServicesTxLoggedUserRepository, user *models.User, TxID string) Service {
	return Service{repository: repository, user: user, txID: TxID}
}

func (s Service) CreateTxLoggedUser(Event string, HostName string, IpRequest string, IpRemote string, UserId string) (*TxLoggedUser, int, error) {
	m := NewCreateTxLoggedUser(Event, HostName, IpRequest, IpRemote, UserId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.Create(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't create TxLoggedUser :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s Service) UpdateTxLoggedUser(id int64, Event string, HostName string, IpRequest string, IpRemote string, UserId string) (*TxLoggedUser, int, error) {
	m := NewTxLoggedUser(id, Event, HostName, IpRequest, IpRemote, UserId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.Update(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update TxLoggedUser :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s Service) DeleteTxLoggedUser(id int64) (int, error) {
	if err := s.repository.Delete(id); err != nil {
		if err.Error() == "ecatch:108" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update row:", err)
		return 20, err
	}
	return 28, nil
}

func (s Service) GetTxLoggedUserByID(id int64) (*TxLoggedUser, int, error) {
	m, err := s.repository.GetByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s Service) GetAllTxLoggedUser() ([]*TxLoggedUser, error) {
	return s.repository.GetAll()
}

func (s Service) GetTxLoggedUserByUser(id string) ([]*TxLoggedUser, int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return nil, 15, fmt.Errorf("id isn't uuid")
	}
	m, err := s.repository.GetByUserID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

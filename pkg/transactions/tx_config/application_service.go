package tx_config

import (
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
)

type Service struct {
	repository ServicesTxConfigRepository
	user       *models.User
	txID       string
}

func NewTxConfigService(repository ServicesTxConfigRepository, user *models.User, TxID string) Service {
	return Service{repository: repository, user: user, txID: TxID}
}

func (s Service) CreateTxConfig(Action string, Description string, UserId string) (*TxConfig, int, error) {
	m := NewCreteTxConfig(Action, Description, UserId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.Create(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't create TxConfig :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s Service) UpdateTxConfig(id int64, Action string, Description string, UserId string) (*TxConfig, int, error) {
	m := NewTxConfig(id, Action, Description, UserId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.Update(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update TxConfig :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s Service) DeleteTxConfig(id int64) (int, error) {
	if err := s.repository.Delete(id); err != nil {
		if err.Error() == "ecatch:108" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update row:", err)
		return 20, err
	}
	return 28, nil
}

func (s Service) GetTxConfigByID(id int64) (*TxConfig, int, error) {
	m, err := s.repository.GetByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s Service) GetAllTxConfig() ([]*TxConfig, error) {
	return s.repository.GetAll()
}

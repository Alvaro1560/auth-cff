package parameters

import (
	"fmt"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
)

type Service struct {
	repository ServicesParameterRepository
	user       *models.User
	txID       string
}

func NewParameterService(repository ServicesParameterRepository, user *models.User, TxID string) Service {
	return Service{repository: repository, user: user, txID: TxID}
}

func (s Service) CreateParameter(Name string, Value string, Type string, Description string, ClientId int, IsCipher bool) (*Parameter, int, error) {
	m := NewParameter(0, Name, Value, Type, Description, ClientId, IsCipher)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.Create(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't create Parameter :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s Service) UpdateParameter(id int, Name string, Value string, Type string, Description string, ClientId int, IsCipher bool) (*Parameter, int, error) {
	m := NewParameter(id, Name, Value, Type, Description, ClientId, IsCipher)
	if id <= 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't int"))
		return m, 15, fmt.Errorf("id isn't int")
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.Update(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update Parameter :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s Service) DeleteParameter(id int) (int, error) {
	if id <= 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't int"))
		return 15, fmt.Errorf("id isn't int")
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

func (s Service) GetParameterByID(id int) (*Parameter, int, error) {
	if id <= 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't int"))
		return nil, 15, fmt.Errorf("id isn't int")
	}
	m, err := s.repository.GetByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s Service) GetParameterByName(name string) (*Parameter, int, error) {

	m, err := s.repository.GetByName(name)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't GetParameterByName row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s Service) GetAllParameter() ([]*Parameter, error) {
	return s.repository.GetAll()
}

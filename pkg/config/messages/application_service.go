package messages

import (
	"fmt"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
)

type Service struct {
	repository ServicesMessageRepository
	user       *models.User
	txID       string
}

func NewMessageService(repository ServicesMessageRepository, user *models.User, TxID string) Service {
	return Service{repository: repository, user: user, txID: TxID}
}

func (s Service) CreateMessage(id int, Spa string, Eng string, TypeMessage string) (*Message, int, error) {
	m := NewMessage(id, Spa, Eng, TypeMessage)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.Create(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't create Message :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s Service) UpdateMessage(id int, Spa string, Eng string, TypeMessage string) (*Message, int, error) {
	m := NewMessage(id, Spa, Eng, TypeMessage)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.Update(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update Message :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s Service) DeleteMessage(id int) (int, error) {
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

func (s Service) GetMessageByID(id int) (*Message, int, error) {
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

func (s Service) GetAllMessage() ([]*Message, error) {
	return s.repository.GetAll()
}

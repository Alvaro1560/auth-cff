package trace

import (
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
)

type Service struct {
	repository ServicesTxTraceRepository
	user       *models.User
	txID       string
}

func NewTxTraceService(repository ServicesTxTraceRepository, user *models.User, TxID string) Service {
	return Service{repository: repository, user: user, txID: TxID}
}

func (s Service) CreateTxTrace(TypeMessage string, FileName string, CodeLine int, Message string, TransactionId string) (*TxTrace, int, error) {
	m := NewCreateTxTrace(TypeMessage, FileName, CodeLine, Message, TransactionId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.Create(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't create TxTrace :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s Service) UpdateTxTrace(id int64, TypeMessage string, FileName string, CodeLine int, Message string, TransactionId string) (*TxTrace, int, error) {
	m := NewTxTrace(id, TypeMessage, FileName, CodeLine, Message, TransactionId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.Update(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update TxTrace :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s Service) DeleteTxTrace(id int64) (int, error) {

	if err := s.repository.Delete(id); err != nil {
		if err.Error() == "ecatch:108" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't Delete row:", err)
		return 20, err
	}
	return 28, nil
}

func (s Service) DeleteTxTraceByTypeMsgs(typeMsgs string) (int, error) {
	if err := s.repository.DeleteByTypeMsg(typeMsgs); err != nil {
		if err.Error() == "ecatch:108" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't Delete row:", err)
		return 20, err
	}
	return 28, nil
}

func (s Service) GetTxTraceByID(id int64) (*TxTrace, int, error) {
	m, err := s.repository.GetByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s Service) GetTxTraceByTypeMsgs(typeMsg string) ([]*TxTrace, int, error) {
	m, err := s.repository.GetByTypeMsg(typeMsg)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s Service) GetAllTxTrace() ([]*TxTrace, error) {
	return s.repository.GetAll()
}

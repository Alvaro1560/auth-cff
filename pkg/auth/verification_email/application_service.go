package verification_email

import (
	"fmt"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
	"time"
)

type PortsServerVerificationEmail interface {
	CreateVerificationEmail(email string, verificationCode string, identification string, verificationDate *time.Time) (*VerificationEmail, int, error)
	UpdateVerificationEmail(id int64, email string, verificationCode string, identification string, verificationDate *time.Time) (*VerificationEmail, int, error)
	DeleteVerificationEmail(id int64) (int, error)
	GetVerificationEmailByID(id int64) (*VerificationEmail, int, error)
	GetAllVerificationEmail() ([]*VerificationEmail, error)
}

type service struct {
	repository ServicesVerificationEmailRepository
	user       *models.User
	txID       string
}

func NewVerificationEmailService(repository ServicesVerificationEmailRepository, user *models.User, TxID string) PortsServerVerificationEmail {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateVerificationEmail(email string, verificationCode string, identification string, verificationDate *time.Time) (*VerificationEmail, int, error) {
	m := NewCreateVerificationEmail(email, verificationCode, identification, verificationDate)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create VerificationEmail :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateVerificationEmail(id int64, email string, verificationCode string, identification string, verificationDate *time.Time) (*VerificationEmail, int, error) {
	m := NewVerificationEmail(id, email, verificationCode, identification, verificationDate)
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update VerificationEmail :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteVerificationEmail(id int64) (int, error) {
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return 15, fmt.Errorf("id is required")
	}

	if err := s.repository.delete(id); err != nil {
		if err.Error() == "ecatch:108" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update row:", err)
		return 20, err
	}
	return 28, nil
}

func (s *service) GetVerificationEmailByID(id int64) (*VerificationEmail, int, error) {
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return nil, 15, fmt.Errorf("id is required")
	}
	m, err := s.repository.getByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllVerificationEmail() ([]*VerificationEmail, error) {
	return s.repository.getAll()
}

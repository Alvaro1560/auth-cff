package verification_email

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// VerificationEmail  Model struct VerificationEmail
type VerificationEmail struct {
	ID               int64      `json:"id" db:"id" valid:"-"`
	Email            string     `json:"email" db:"email" valid:"required"`
	VerificationCode string     `json:"verification_code" db:"verification_code" valid:"-"`
	Identification   string     `json:"identification" db:"identification" valid:"-"`
	VerificationDate *time.Time `json:"verification_date" db:"verification_date" valid:"-"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

func NewVerificationEmail(id int64, email string, verificationCode string, identification string, verificationDate *time.Time) *VerificationEmail {
	return &VerificationEmail{
		ID:               id,
		Email:            email,
		VerificationCode: verificationCode,
		Identification:   identification,
		VerificationDate: verificationDate,
	}
}

func NewCreateVerificationEmail(email string, verificationCode string, identification string, verificationDate *time.Time) *VerificationEmail {
	return &VerificationEmail{
		Email:            email,
		VerificationCode: verificationCode,
		Identification:   identification,
		VerificationDate: verificationDate,
	}
}

func (m *VerificationEmail) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

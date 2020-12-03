package users

import (
	"github.com/asaskevich/govalidator"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
)

type User models.User

func NewUser(id string, Username string, Name string, LastName string, EmailNotifications string, IdentificationNumber string, IdentificationType string) *User {
	return &User{
		ID:                   id,
		Username:             Username,
		Name:                 Name,
		LastName:             LastName,
		EmailNotifications:   EmailNotifications,
		IdentificationNumber: IdentificationNumber,
		IdentificationType:   IdentificationType,
	}
}

func (m *User) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

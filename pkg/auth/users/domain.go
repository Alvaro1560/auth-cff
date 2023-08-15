package users

import (
	"github.com/asaskevich/govalidator"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
)

type User models.User

func NewUser(id string, Username string, CodeStudent string, Dni string, Names string, LastnameFather string, LastnameMother string, Email string) *User {
	return &User{
		ID:             id,
		Username:       Username,
		CodeStudent:    CodeStudent,
		Dni:            Dni,
		Names:          Names,
		LastnameFather: LastnameFather,
		LastnameMother: LastnameMother,
		Email:          Email,
	}
}

func (m *User) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

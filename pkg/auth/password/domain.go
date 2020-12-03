package password

import (
	"github.com/asaskevich/govalidator"
)

// Model estructura de Role
type PasswordPolicy struct {
	ID       string `json:"id" valid:"required,uuid"`
	Password string `json:"password" valid:"required"`
}

/*
func NewPassword(id string, Password string) *PasswordPolicy {
	return *PasswordPolicy{
		ID:       id,
		Password: Password,
	}
}
*/
func (m *PasswordPolicy) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

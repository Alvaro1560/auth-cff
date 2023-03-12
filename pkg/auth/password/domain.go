package password

import (
	"github.com/asaskevich/govalidator"
	"time"
)

// Password estructura de Role
type Password struct {
	ID        string     `json:"id" db:"id" valid:"required,uuid"`
	Password  string     `json:"password" db:"password" valid:"required"`
	UserId    string     `json:"user_id" db:"user_id"`
	IdUser    string     `json:"id_user" db:"id_user" valid:"-"`
	IsDelete  bool       `json:"is_delete" db:"is_delete" valid:"-"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

func (m *Password) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

package roles

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de Role
type Role struct {
	ID              string    `json:"id" db:"id" valid:"required,uuid"`
	Name            string    `json:"name" db:"name" valid:"required"`
	Description     string    `json:"description" db:"description" valid:"required"`
	SessionsAllowed int       `json:"sessions_allowed" db:"sessions_allowed" valid:"required"`
	ProcessID       string    `json:"process_id,omitempty" db:"process_id" valid:"-"`
	UserID          string    `json:"user_id,omitempty" db:"user_id" valid:"-"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

func NewRole(id string, Name string, Description string, SessionsAllowed int) *Role {
	return &Role{
		ID:              id,
		Name:            Name,
		Description:     Description,
		SessionsAllowed: SessionsAllowed,
	}
}

func (m *Role) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

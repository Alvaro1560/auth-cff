package parameters

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de Parameter
type Parameter struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" valid:"required,stringlength(1|50),matches(^[A-Z0-9_]+$)"`
	Value       string    `json:"value" db:"value" valid:"required,stringlength(1|100)"`
	Type        string    `json:"type" db:"type" valid:"required,stringlength(1|15)"`
	Description string    `json:"description" db:"description" valid:"required,stringlength(1|500)"`
	ClientId    int       `json:"client_id" db:"client_id" valid:"required"`
	IsCipher    bool      `json:"is_cipher" db:"is_cipher" valid:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func NewParameter(id int, Name string, Value string, Type string, Description string, ClientId int, IsCipher bool) *Parameter {
	return &Parameter{
		ID:          id,
		Name:        Name,
		Value:       Value,
		Type:        Type,
		Description: Description,
		ClientId:    ClientId,
		IsCipher:    IsCipher,
	}
}

func (m *Parameter) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

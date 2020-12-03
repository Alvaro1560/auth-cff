package messages

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de Message
type Message struct {
	ID          int       `json:"id" db:"id" valid:"required"`
	Spa         string    `json:"spa" db:"spa" valid:"required,stringlength(1|255)"`
	Eng         string    `json:"eng" db:"eng" valid:"required,stringlength(1|255)"`
	TypeMessage string    `json:"type_message" db:"type_message" valid:"required,stringlength(1|10)"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func NewMessage(id int, Spa string, Eng string, TypeMessage string) *Message {
	return &Message{
		ID:          id,
		Spa:         Spa,
		Eng:         Eng,
		TypeMessage: TypeMessage,
	}
}

func (m *Message) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

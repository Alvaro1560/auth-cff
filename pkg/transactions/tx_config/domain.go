package tx_config

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de TxConfig
type TxConfig struct {
	ID          int64     `json:"id" db:"id" valid:"-"`
	Action      string    `json:"action" db:"action" valid:"required"`
	Description string    `json:"description" db:"description" valid:"required"`
	UserId      string    `json:"user_id" db:"user_id" valid:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func NewCreteTxConfig(Action string, Description string, UserId string) *TxConfig {
	return &TxConfig{
		Action:      Action,
		Description: Description,
		UserId:      UserId,
	}
}

func NewTxConfig(id int64, Action string, Description string, UserId string) *TxConfig {
	return &TxConfig{
		ID:          id,
		Action:      Action,
		Description: Description,
		UserId:      UserId,
	}
}

func (m *TxConfig) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

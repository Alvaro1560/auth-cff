package trace

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de TxTrace
type TxTrace struct {
	ID            int64     `json:"id" db:"id" valid:"-"`
	TypeMessage   string    `json:"type_message" db:"type_message" valid:"required"`
	FileName      string    `json:"file_name" db:"file_name" valid:"required"`
	CodeLine      int       `json:"code_line" db:"code_line" valid:"required"`
	Message       string    `json:"message" db:"message" valid:"required"`
	TransactionId string    `json:"transaction_id" db:"transaction_id" valid:"required"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

func NewCreateTxTrace(TypeMessage string, FileName string, CodeLine int, Message string, TransactionId string) *TxTrace {
	return &TxTrace{
		TypeMessage:   TypeMessage,
		FileName:      FileName,
		CodeLine:      CodeLine,
		Message:       Message,
		TransactionId: TransactionId,
	}
}
func NewTxTrace(id int64, TypeMessage string, FileName string, CodeLine int, Message string, TransactionId string) *TxTrace {
	return &TxTrace{
		ID:            id,
		TypeMessage:   TypeMessage,
		FileName:      FileName,
		CodeLine:      CodeLine,
		Message:       Message,
		TransactionId: TransactionId,
	}
}

func (m *TxTrace) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

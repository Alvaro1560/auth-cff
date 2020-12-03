package loggedusers

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de TxLoggedUser
type TxLoggedUser struct {
	ID        int64     `json:"id" db:"id" valid:"-"`
	Event     string    `json:"event" db:"event" valid:"required"`
	HostName  string    `json:"host_name" db:"host_name" valid:"required"`
	IpRequest string    `json:"ip_request" db:"ip_request" valid:"required"`
	IpRemote  string    `json:"ip_remote" db:"ip_remote" valid:"required"`
	UserId    string    `json:"user_id" db:"user_id" valid:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewCreateTxLoggedUser(Event string, HostName string, IpRequest string, IpRemote string, UserId string) *TxLoggedUser {
	return &TxLoggedUser{
		Event:     Event,
		HostName:  HostName,
		IpRequest: IpRequest,
		IpRemote:  IpRemote,
		UserId:    UserId,
	}
}

func NewTxLoggedUser(id int64, Event string, HostName string, IpRequest string, IpRemote string, UserId string) *TxLoggedUser {
	return &TxLoggedUser{
		ID:        id,
		Event:     Event,
		HostName:  HostName,
		IpRequest: IpRequest,
		IpRemote:  IpRemote,
		UserId:    UserId,
	}
}

func (m *TxLoggedUser) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

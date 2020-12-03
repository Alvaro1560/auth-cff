package login

import (
	"github.com/asaskevich/govalidator"
)

// Model estructura de Module
type Login struct {
	ID       string `json:"id" db:"id" valid:"-"`
	Username string `json:"username" db:"username" valid:"required"`
	Password string `json:"password" db:"password"`
	ClientID int    `json:"client_id" db:"client_id"`
	HostName string `json:"host_name" db:"host_name"`
	RealIP   string `json:"real_ip" db:"real_ip"`
}

func NewLogin(id, Username, Password string, ClientID int, HostName, RealIP string) *Login {
	return &Login{
		ID:       id,
		Username: Username,
		Password: Password,
		ClientID: ClientID,
		HostName: HostName,
		RealIP:   RealIP,
	}
}

func (m *Login) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

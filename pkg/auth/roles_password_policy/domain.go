package roles_password_policy

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// RolesPasswordPolicy estructura de RolesPasswordPolicy
type RolesPasswordPolicy struct {
	ID                   string    `json:"id" db:"id" valid:"required,uuid"`
	RoleId               string    `json:"role_id" db:"role_id" valid:"required"`
	DaysPassValid        int       `json:"days_pass_valid" db:"days_pass_valid" valid:"required"`
	MaxLength            int       `json:"max_length" db:"max_length" valid:"required"`
	MinLength            int       `json:"min_length" db:"min_length" valid:"required"`
	StorePassNotRepeated int       `json:"store_pass_not_repeated" db:"store_pass_not_repeated" valid:"required"`
	FailedAttempts       int       `json:"failed_attempts" db:"failed_attempts" valid:"required"`
	TimeUnlock           int       `json:"time_unlock" db:"time_unlock" valid:"required"`
	Alpha                int       `json:"alpha" db:"alpha" valid:"required"`
	Digits               int       `json:"digits" db:"digits" valid:"required"`
	Special              int       `json:"special" db:"special" valid:"required"`
	UpperCase            int       `json:"upper_case" db:"upper_case" valid:"required"`
	LowerCase            int       `json:"lower_case" db:"lower_case" valid:"required"`
	Enable               bool      `json:"enable" db:"enable" valid:"required"`
	InactivityTime       int       `json:"inactivity_time" db:"inactivity_time" valid:"required"`
	Timeout              int       `json:"timeout" db:"timeout" valid:"required"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
	IdUser               string    `json:"id_user" db:"id_user"`
	IsDelete             bool      `json:"is_delete" db:"is_delete"`
	Required2fa          bool      `json:"required_2fa" db:"required_2fa"`
}

func NewRolesPasswordPolicy(id string, RoleId string, DaysPassValid int, MaxLength int, MinLength int, StorePassNotRepeated int, FailedAttempts int, TimeUnlock int, Alpha int, Digits int, Special int, UpperCase int, LowerCase int, Enable bool, InactivityTime int, Timeout int, required2fa bool) *RolesPasswordPolicy {
	return &RolesPasswordPolicy{
		ID:                   id,
		RoleId:               RoleId,
		DaysPassValid:        DaysPassValid,
		MaxLength:            MaxLength,
		MinLength:            MinLength,
		StorePassNotRepeated: StorePassNotRepeated,
		FailedAttempts:       FailedAttempts,
		TimeUnlock:           TimeUnlock,
		Alpha:                Alpha,
		Digits:               Digits,
		Special:              Special,
		UpperCase:            UpperCase,
		LowerCase:            LowerCase,
		Enable:               Enable,
		InactivityTime:       InactivityTime,
		Timeout:              Timeout,
		Required2fa:          required2fa,
	}
}

func (m *RolesPasswordPolicy) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

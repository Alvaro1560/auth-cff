package customers_projects

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de Project
type Project struct {
	ID           string    `json:"id" db:"id" valid:"required,uuid"`
	Name         string    `json:"name" db:"name" valid:"required"`
	Description  string    `json:"description" db:"description" valid:"required"`
	Department   string    `json:"department" db:"department" valid:"required"`
	Email        string    `json:"email" db:"email" valid:"required"`
	Phone        string    `json:"phone" db:"phone" valid:"required"`
	ProductOwner string    `json:"product_owner" db:"product_owner" valid:"required"`
	CustomersId  string    `json:"customers_id" db:"customers_id" valid:"required"`
	RoleId       string    `json:"role_id,omitempty" db:"role_id" valid:"required"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func NewProject(id string, Name string, Description string, Department string, Email string, Phone string, ProductOwner string, CustomersId string) *Project {
	return &Project{
		ID:           id,
		Name:         Name,
		Description:  Description,
		Department:   Department,
		Email:        Email,
		Phone:        Phone,
		ProductOwner: ProductOwner,
		CustomersId:  CustomersId,
	}
}

func (m *Project) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

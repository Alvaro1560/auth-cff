package customers_projects

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
)

type Service struct {
	repository ServicesProjectRepository
	user       *models.User
	txID       string
}

func NewProjectService(repository ServicesProjectRepository, user *models.User, TxID string) Service {
	return Service{repository: repository, user: user, txID: TxID}
}

func (s Service) CreateProject(id string, Name string, Description string, department string, Email string, Phone string, ProductOwner string, CustomersId string) (*Project, int, error) {
	m := NewProject(id, Name, Description, department, Email, Phone, ProductOwner, CustomersId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.Create(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't create Project :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s Service) UpdateProject(id string, Name string, Description string, department string, Email string, Phone string, ProductOwner string, CustomersId string) (*Project, int, error) {
	m := NewProject(id, Name, Description, department, Email, Phone, ProductOwner, CustomersId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.Update(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update Project :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s Service) DeleteProject(id string) (int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return 15, fmt.Errorf("id isn't uuid")
	}

	if err := s.repository.Delete(id); err != nil {
		if err.Error() == "ecatch:108" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update row:", err)
		return 20, err
	}
	return 28, nil
}

func (s Service) GetProjectByID(id string) (*Project, int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return nil, 15, fmt.Errorf("id isn't uuid")
	}
	m, err := s.repository.GetByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s Service) GetAllProject() ([]*Project, error) {
	return s.repository.GetAll()
}

func (s Service) GetProjectByRoles(roleIDs []string) ([]*string, error) {
	return s.repository.getProjectByRoles(roleIDs)
}

func (s Service) GetProjectByRolesAndProjectID(roleIDs, projectID []string) ([]*Project, error) {
	return s.repository.getProjectByRolesAndProjectID(roleIDs, projectID)
}

func (s Service) GetProjectsByIds(projectID []string) ([]*Project, error) {
	return s.repository.getProjectsByIds(projectID)
}
func (s Service) GetProjectByRoleIDs(RoleIDs []string) ([]*Project, error) {
	return s.repository.GetProjectByRoleIDs(RoleIDs)
}

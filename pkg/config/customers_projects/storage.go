package customers_projects

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
)

const (
	Postgresql = "postgres"
	SqlServer  = "sqlserver"
	Oracle     = "oci8"
)

type ServicesProjectRepository interface {
	Create(m *Project) error
	Update(m *Project) error
	Delete(id string) error
	GetByID(id string) (*Project, error)
	GetAll() ([]*Project, error)
	getProjectByRoles(roleIDs []string) ([]*string, error)
	getProjectByRolesAndProjectID(roleIDs, projectID []string) ([]*Project, error)
	getProjectsByIds(projectID []string) ([]*Project, error)
	GetProjectByRoleIDs(RoleIDs []string) ([]*Project, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesProjectRepository {
	var s ServicesProjectRepository
	engine := db.DriverName()
	switch engine {
	case SqlServer:
		return NewProjectSqlServerRepository(db, user, txID)
	case Postgresql:
		return NewProjectPsqlRepository(db, user, txID)
	case Oracle:
		return NewProjectOrclRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}

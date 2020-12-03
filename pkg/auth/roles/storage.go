package roles

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

type ServicesRoleRepository interface {
	Create(m *Role) error
	Update(m *Role) error
	Delete(id string) error
	GetByID(id string) (*Role, error)
	GetAll() ([]*Role, error)
	GetByUserID(userID string) ([]*Role, error)
	GetRolesByProcessIDs(ProcessIDs []string) ([]*Role, error)
	GetRolesByIDs(IDs []string) ([]*Role, error)
	GetByUserIDs(userIDs []string) ([]*Role, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesRoleRepository {
	var s ServicesRoleRepository
	engine := db.DriverName()
	switch engine {
	case SqlServer:
		return NewRoleSqlServerRepository(db, user, txID)
	case Postgresql:
		return NewRolePsqlRepository(db, user, txID)
	case Oracle:
		return NewRoleOrclRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}

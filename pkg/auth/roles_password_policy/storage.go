package roles_password_policy

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

type ServicesRolesPasswordPolicyRepository interface {
	Create(m *RolesPasswordPolicy) error
	Update(m *RolesPasswordPolicy) error
	Delete(id string) error
	GetByID(id string) (*RolesPasswordPolicy, error)
	GetAll() ([]*RolesPasswordPolicy, error)
	GetAllRolesPasswordPolicyByRolesIDs(RolesIDs []string) ([]*RolesPasswordPolicy, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesRolesPasswordPolicyRepository {
	var s ServicesRolesPasswordPolicyRepository
	engine := db.DriverName()
	switch engine {
	case SqlServer:
		return NewRolesPasswordPolicySqlServerRepository(db, user, txID)
	case Postgresql:
		return NewRolesPasswordPolicyPsqlRepository(db, user, txID)
	case Oracle:
		return NewRolesPasswordPolicyOrclRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}

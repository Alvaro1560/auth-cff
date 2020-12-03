package parameters

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

type ServicesParameterRepository interface {
	Create(m *Parameter) error
	Update(m *Parameter) error
	Delete(id int) error
	GetByID(id int) (*Parameter, error)
	GetByName(name string) (*Parameter, error)
	GetAll() ([]*Parameter, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesParameterRepository {
	var s ServicesParameterRepository
	engine := db.DriverName()
	switch engine {
	case SqlServer:
		return NewParameterSqlServerRepository(db, user, txID)
	case Postgresql:
		return NewParameterPsqlRepository(db, user, txID)
	case Oracle:
		return NewParameterOrclRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}

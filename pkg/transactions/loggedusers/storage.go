package loggedusers

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

type ServicesTxLoggedUserRepository interface {
	Create(m *TxLoggedUser) error
	Update(m *TxLoggedUser) error
	Delete(id int64) error
	GetByID(id int64) (*TxLoggedUser, error)
	GetAll() ([]*TxLoggedUser, error)
	GetByUserID(id string) ([]*TxLoggedUser, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesTxLoggedUserRepository {
	var s ServicesTxLoggedUserRepository
	engine := db.DriverName()
	switch engine {
	case SqlServer:
		return NewTxLoggedUserSqlServerRepository(db, user, txID)
	case Postgresql:
		return NewTxLoggedUserPsqlRepository(db, user, txID)
	case Oracle:
		return NewTxLoggedUserOrclRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}

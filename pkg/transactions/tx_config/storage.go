package tx_config

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

type ServicesTxConfigRepository interface {
	Create(m *TxConfig) error
	Update(m *TxConfig) error
	Delete(id int64) error
	GetByID(id int64) (*TxConfig, error)
	GetAll() ([]*TxConfig, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesTxConfigRepository {
	var s ServicesTxConfigRepository
	engine := db.DriverName()
	switch engine {
	case SqlServer:
		return NewTxConfigSqlServerRepository(db, user, txID)
	case Postgresql:
		return NewTxConfigPsqlRepository(db, user, txID)
	case Oracle:
		return NewTxConfigOrclRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}

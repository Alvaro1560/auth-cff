package trace

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

type ServicesTxTraceRepository interface {
	Create(m *TxTrace) error
	Update(m *TxTrace) error
	Delete(id int64) error
	DeleteByTypeMsg(typeMsg string) error
	GetByID(id int64) (*TxTrace, error)
	GetByTypeMsg(typeMsg string) ([]*TxTrace, error)
	GetAll() ([]*TxTrace, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesTxTraceRepository {
	var s ServicesTxTraceRepository
	engine := db.DriverName()
	switch engine {
	case SqlServer:
		return NewTxTraceSqlServerRepository(db, user, txID)
	case Postgresql:
		return NewTxTracePsqlRepository(db, user, txID)
	case Oracle:
		return NewTxTraceOrclRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}

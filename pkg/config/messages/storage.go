package messages

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

type ServicesMessageRepository interface {
	Create(m *Message) error
	Update(m *Message) error
	Delete(id int) error
	GetByID(id int) (*Message, error)
	GetAll() ([]*Message, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesMessageRepository {
	var s ServicesMessageRepository
	engine := db.DriverName()
	switch engine {
	case SqlServer:
		return NewMessageSqlServerRepository(db, user, txID)
	case Postgresql:
		return NewMessagePsqlRepository(db, user, txID)
	case Oracle:
		return NewMessageOrclRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}

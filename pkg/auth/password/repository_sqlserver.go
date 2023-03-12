package password

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
)

// sqlServer estructura de conexión a la BD de mssql
type sqlserver struct {
	DB   *sqlx.DB
	user *models.User
	TxID string
}

func NewRolesPasswordSqlServerRepository(db *sqlx.DB, user *models.User, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// GetByID consulta un registro por su ID
func (s *sqlserver) GetByUserID(userId string) (*Password, error) {
	const psqlGetByID = `SELECT top(1) id, user_id, "password", created_at, id_user, is_delete, deleted_at FROM auth.users_password_history where user_id = @user_id order by created_at desc`
	mdl := Password{}
	err := s.DB.Get(&mdl, psqlGetByID, sql.Named("user_id", userId))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByID Password: %v", err)
		return &mdl, err
	}
	return &mdl, nil
}

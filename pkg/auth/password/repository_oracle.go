package password

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
)

// Orcl estructura de conexión a la BD de Oracle
type orcl struct {
	DB   *sqlx.DB
	user *models.User
	TxID string
}

func NewRolesPasswordOrclRepository(db *sqlx.DB, user *models.User, txID string) *orcl {
	return &orcl{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// GetByUserID consulta un registro por su ID
func (s *orcl) GetByUserID(userId string) (*Password, error) {
	const osqlGetByID = `SELECT id, user_id, "password", created_at, id_user, is_delete, deleted_at FROM auth.users_password_history where user_id = :1 order by created_at desc limit 1`
	mdl := Password{}
	err := s.DB.Get(&mdl, osqlGetByID, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByUserID Password: %v", err)
		return &mdl, err
	}
	return &mdl, nil
}

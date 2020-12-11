package tx_config

import (
	"database/sql"
	"fmt"

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

func NewTxConfigSqlServerRepository(db *sqlx.DB, user *models.User, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *sqlserver) Create(m *TxConfig) error {
	const sqlInsert = `INSERT INTO tx.config ([action], description, user_id) VALUES (@action, @description, @user_id) SELECT ID = convert(bigint, SCOPE_IDENTITY()) `
	stmt, err := s.DB.Prepare(sqlInsert)
	if err != nil {
		return nil
	}
	err = stmt.QueryRow(
		sql.Named("action", m.Action),
		sql.Named("description", m.Description),
		sql.Named("user_id", m.UserId),
	).Scan(&m.ID)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert tx config: %v", err)
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *sqlserver) Update(m *TxConfig) error {
	const sqlUpdate = `UPDATE tx.config SET action = :action, description = :description, user_id = :user_id, updated_at = GetDate() WHERE id = :id `
	rs, err := s.DB.NamedExec(sqlUpdate, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't update TxConfig: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Delete elimina un registro de la BD
func (s *sqlserver) Delete(id int64) error {
	const sqlDelete = `DELETE FROM tx.config WHERE id = :id `
	m := TxConfig{ID: id}
	rs, err := s.DB.NamedExec(sqlDelete, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't delete TxConfig: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *sqlserver) GetByID(id int64) (*TxConfig, error) {
	const sqlGetByID = `SELECT  id , action, description, CONVERT(NVARCHAR(50), user_id) user_id, created_at, updated_at FROM tx.config  WITH (NOLOCK)  WHERE id = @id `
	mdl := TxConfig{}
	err := s.DB.Get(&mdl, sqlGetByID, sql.Named("id", id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByID TxConfig: %v", err)
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *sqlserver) GetAll() ([]*TxConfig, error) {
	var ms []*TxConfig
	const sqlGetAll = `SELECT  id , action, description, CONVERT(NVARCHAR(50), user_id) user_id, created_at, updated_at FROM tx.config  WITH (NOLOCK) `

	err := s.DB.Select(&ms, sqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetAll config: %v", err)
		return ms, err
	}
	return ms, nil
}

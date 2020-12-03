package tx_config

import (
	"database/sql"
	"fmt"

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

func NewTxConfigOrclRepository(db *sqlx.DB, user *models.User, txID string) *orcl {
	return &orcl{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *orcl) Create(m *TxConfig) error {
	const orclInsert = `INSERT INTO config (action, description, user_id) VALUES ($1, $2, $3,) RETURNING id `
	stmt, err := s.DB.Prepare(orclInsert)
	if err != nil {
		return nil
	}
	err = stmt.QueryRow(
		m.Action,
		m.Description,
		m.UserId,
	).Scan(&m.ID)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert tx config: %v", err)
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *orcl) Update(m *TxConfig) error {
	const orclUpdate = `UPDATE config SET action = :action, description = :description, user_id = :user_id, updated_at = now() WHERE id = :id `
	rs, err := s.DB.NamedExec(orclUpdate, &m)
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
func (s *orcl) Delete(id int64) error {
	const orclDelete = `DELETE FROM config WHERE id = :id `
	m := TxConfig{ID: id}
	rs, err := s.DB.NamedExec(orclDelete, &m)
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
func (s *orcl) GetByID(id int64) (*TxConfig, error) {
	const orclGetByID = `SELECT id , action, description, user_id,  created_at, updated_at FROM config WHERE id = $1 `
	mdl := TxConfig{}
	err := s.DB.Get(&mdl, orclGetByID, id)
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
func (s *orcl) GetAll() ([]*TxConfig, error) {
	var ms []*TxConfig
	const orclGetAll = ` SELECT id , action, description, user_id, created_at, updated_at FROM config `

	err := s.DB.Select(&ms, orclGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetAll TxConfig: %v", err)
		return ms, err
	}
	return ms, nil
}

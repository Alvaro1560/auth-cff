package tx_config

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
)

// psql estructura de conexión a la BD de postgresql
type psql struct {
	DB   *sqlx.DB
	user *models.User
	TxID string
}

func NewTxConfigPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) Create(m *TxConfig) error {
	const psqlInsert = `INSERT INTO config ("action", description, user_id) VALUES ($1, $2, $3,) RETURNING id `
	stmt, err := s.DB.Prepare(psqlInsert)
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
func (s *psql) Update(m *TxConfig) error {
	const psqlUpdate = `UPDATE config SET action = :action, description = :description, user_id = :user_id, updated_at = now() WHERE id = :id `
	rs, err := s.DB.NamedExec(psqlUpdate, &m)
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
func (s *psql) Delete(id int64) error {
	const psqlDelete = `DELETE FROM config WHERE id = :id `
	m := TxConfig{ID: id}
	rs, err := s.DB.NamedExec(psqlDelete, &m)
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
func (s *psql) GetByID(id int64) (*TxConfig, error) {
	const psqlGetByID = `SELECT id , action, description, user_id,  created_at, updated_at FROM config WHERE id = $1 `
	mdl := TxConfig{}
	err := s.DB.Get(&mdl, psqlGetByID, id)
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
func (s *psql) GetAll() ([]*TxConfig, error) {
	var ms []*TxConfig
	const psqlGetAll = ` SELECT id , action, description, user_id, created_at, updated_at FROM config `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetAll TxConfig: %v", err)
		return ms, err
	}
	return ms, nil
}

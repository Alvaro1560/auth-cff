package parameters

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

func NewParameterPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) Create(m *Parameter) error {
	var id int64
	const sqlInsert = `INSERT INTO cfg.parameters (name, value, type, description, client_id, is_cipher,created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6,Now(), Now())   RETURNING id `
	stmt, err := s.DB.Prepare(sqlInsert)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert Parameters: %v", err)
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(m.Name, m.Value, m.Type, m.Description, m.ClientId, m.IsCipher).Scan(&id)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert Parameters: %v", err)
		return err
	}
	m.ID = int(id)
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) Update(m *Parameter) error {
	const sqlUpdate = `UPDATE cfg.parameters SET name = :name, value = :value, type = :type, description = :description, client_id = :client_id, is_cipher = :is_cipher, updated_at = Now() WHERE id = :id `
	rs, err := s.DB.NamedExec(sqlUpdate, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't update Parameter: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Delete elimina un registro de la BD
func (s *psql) Delete(id int) error {
	const sqlDelete = `DELETE FROM cfg.parameters WHERE id = :id `
	m := Parameter{ID: id}
	rs, err := s.DB.NamedExec(sqlDelete, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't delete Parameter: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *psql) GetByID(id int) (*Parameter, error) {
	const sqlGetByID = `SELECT  id ,name ,value , type, description, client_id, is_cipher, created_at, updated_at FROM cfg.parameters    WHERE id = $1 `
	mdl := Parameter{}
	err := s.DB.Get(&mdl, sqlGetByID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByID Parameter: %v", err)
		return &mdl, err
	}
	return &mdl, nil
}

// GetByName consulta un registro por su name
func (s *psql) GetByName(name string) (*Parameter, error) {
	const osqlGetByName = `SELECT id ,name ,value , type, description, client_id, is_cipher, created_at, updated_at FROM cfg.parameters WHERE name = $1  `
	mdl := Parameter{}
	err := s.DB.Get(&mdl, osqlGetByName, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByName Parameter: %v", err)
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *psql) GetAll() ([]*Parameter, error) {
	var ms []*Parameter
	const sqlGetAll = `SELECT  id ,name ,value , type, description, client_id, is_cipher, created_at, updated_at FROM cfg.parameters   `

	err := s.DB.Select(&ms, sqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetAll cfg.parameters: %v", err)
		return ms, err
	}
	return ms, nil
}

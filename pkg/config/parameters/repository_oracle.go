package parameters

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

func NewParameterOrclRepository(db *sqlx.DB, user *models.User, txID string) *orcl {
	return &orcl{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *orcl) Create(m *Parameter) error {
	const osqlInsert = `INSERT INTO cfg.parameters (id ,name,value, type, description, client_id, is_cipher)  VALUES (:id ,:name,:value, :type, :description, :client_id, :is_cipher)`
	_, err := s.DB.NamedExec(osqlInsert, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert Parameter: %v", err)
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *orcl) Update(m *Parameter) error {
	const osqlUpdate = `UPDATE cfg.parameters SET name = :name,value = :value, type = :type, description = :description, client_id = :client_id, is_cipher = :is_cipher, updated_at = sysdate WHERE id = :id  `
	rs, err := s.DB.NamedExec(osqlUpdate, &m)
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
func (s *orcl) Delete(id int) error {
	const osqlDelete = `DELETE FROM cfg.parameters WHERE id = :id `
	m := Parameter{ID: id}
	rs, err := s.DB.NamedExec(osqlDelete, &m)
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
func (s *orcl) GetByID(id int) (*Parameter, error) {
	const osqlGetByID = `SELECT id , name,value, type, description, client_id, is_cipher, created_at, updated_at FROM cfg.parameters WHERE id = :1 `
	mdl := Parameter{}
	err := s.DB.Get(&mdl, osqlGetByID, id)
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
func (s *orcl) GetByName(name string) (*Parameter, error) {
	const osqlGetByID = `SELECT id , name,value, type, description, client_id, is_cipher, created_at, updated_at FROM cfg.parameters WHERE name = :1 `
	mdl := Parameter{}
	err := s.DB.Get(&mdl, osqlGetByID, name)
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
func (s *orcl) GetAll() ([]*Parameter, error) {
	var ms []*Parameter
	const osqlGetAll = ` SELECT id , name,value, type, description, client_id, is_cipher, created_at, updated_at FROM cfg.parameters `

	err := s.DB.Select(&ms, osqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetAll Parameter: %v", err)
		return ms, err
	}
	return ms, nil
}

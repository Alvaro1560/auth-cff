package messages

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

func NewMessageOrclRepository(db *sqlx.DB, user *models.User, txID string) *orcl {
	return &orcl{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *orcl) Create(m *Message) error {
	const osqlInsert = `INSERT INTO cfg.messages (id ,spa, eng, type_message)  VALUES (:id ,:spa, :eng, :type_message)`
	_, err := s.DB.NamedExec(osqlInsert, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert Message: %v", err)
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *orcl) Update(m *Message) error {
	const osqlUpdate = `UPDATE cfg.messages SET spa = :spa, eng = :eng, type_message = :type_message, updated_at = sysdate WHERE id = :id  `
	rs, err := s.DB.NamedExec(osqlUpdate, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't update Message: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Delete elimina un registro de la BD
func (s *orcl) Delete(id int) error {
	const osqlDelete = `DELETE FROM cfg.messages WHERE id = :id `
	m := Message{ID: id}
	rs, err := s.DB.NamedExec(osqlDelete, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't delete Message: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *orcl) GetByID(id int) (*Message, error) {
	const osqlGetByID = `SELECT id , spa, eng, type_message, created_at, updated_at FROM cfg.messages WHERE id = :1 `
	mdl := Message{}
	err := s.DB.Get(&mdl, osqlGetByID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByID Message: %v", err)
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *orcl) GetAll() ([]*Message, error) {
	var ms []*Message
	const osqlGetAll = ` SELECT id , spa, eng, type_message, created_at, updated_at FROM cfg.messages `

	err := s.DB.Select(&ms, osqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetAll Message: %v", err)
		return ms, err
	}
	return ms, nil
}

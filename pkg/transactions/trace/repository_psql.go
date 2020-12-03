package trace

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

func NewTxTracePsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) Create(m *TxTrace) error {
	const psqlInsert = `INSERT INTO trace (type_message, file_name, code_line, message, transaction_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	stmt, err := s.DB.Prepare(psqlInsert)
	if err != nil {
		return nil
	}
	err = stmt.QueryRow(
		m.TypeMessage,
		m.FileName,
		m.CodeLine,
		m.Message,
		m.TransactionId,
	).Scan(&m.ID)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert tx trace: %v", err)
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) Update(m *TxTrace) error {
	const psqlUpdate = `UPDATE trace SET type_message = :type_message, file_name = :file_name, code_line = :code_line, message = :message, transaction_id = :transaction_id, updated_at = now() WHERE id = :id `
	rs, err := s.DB.NamedExec(psqlUpdate, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't update TxTrace: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Delete elimina un registro de la BD
func (s *psql) Delete(id int64) error {
	const psqlDelete = `DELETE FROM trace WHERE id = :id `
	m := TxTrace{ID: id}
	rs, err := s.DB.NamedExec(psqlDelete, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't delete TxTrace: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *psql) GetByID(id int64) (*TxTrace, error) {
	const psqlGetByID = `SELECT id , type_message, file_name, code_line, message, transaction_id, created_at, updated_at FROM trace WHERE id = $1 `
	mdl := TxTrace{}
	err := s.DB.Get(&mdl, psqlGetByID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByID TxTrace: %v", err)
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *psql) GetAll() ([]*TxTrace, error) {
	var ms []*TxTrace
	const psqlGetAll = ` SELECT id , type_message, file_name, code_line, message, transaction_id, created_at, updated_at FROM trace `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetAll TxTrace: %v", err)
		return ms, err
	}
	return ms, nil
}

func (s *psql) DeleteByTypeMsg(typeMsg string) error {
	const sqlDelete = `DELETE FROM tx.trace WHERE type_message = :type_message `
	m := TxTrace{TypeMessage: typeMsg}
	rs, err := s.DB.NamedExec(sqlDelete, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't delete TxTrace By Type_Message: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

func (s *psql) GetByTypeMsg(typeMsg string) ([]*TxTrace, error) {
	const sqlGetByTypeMsg = `SELECT  id , type_message, file_name, code_line, message,  transaction_id, created_at FROM tx.trace   WHERE type_message = $1 `
	var ms []*TxTrace
	err := s.DB.Select(&ms, sqlGetByTypeMsg, typeMsg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByTypeMsg TxTrace: %v", err)
		return ms, err
	}
	return ms, nil
}

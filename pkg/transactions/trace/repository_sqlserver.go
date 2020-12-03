package trace

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

func NewTxTraceSqlServerRepository(db *sqlx.DB, user *models.User, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *sqlserver) Create(m *TxTrace) error {
	const sqlInsert = `INSERT INTO tx.trace (type_message, file_name, code_line, message, transaction_id ) VALUES (@type_message, @file_name, @code_line, @message, @transaction_id ) SELECT ID = convert(bigint, SCOPE_IDENTITY())`
	stmt, err := s.DB.Prepare(sqlInsert)
	if err != nil {
		return nil
	}
	err = stmt.QueryRow(
		sql.Named("type_message", m.TypeMessage),
		sql.Named("file_name", m.FileName),
		sql.Named("code_line", m.CodeLine),
		sql.Named("message", m.Message),
		sql.Named("transaction_id", m.TransactionId),
	).Scan(&m.ID)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert tx trace: %v", err)
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *sqlserver) Update(m *TxTrace) error {
	const sqlUpdate = `UPDATE trace SET type_message = :type_message, file_name = :file_name, code_line = :code_line, message = :message, transaction_id = :transaction_id, updated_at = GetDate() WHERE id = :id `
	rs, err := s.DB.NamedExec(sqlUpdate, &m)
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
func (s *sqlserver) Delete(id int64) error {
	const sqlDelete = `DELETE FROM tx.trace WHERE id = :id `
	m := TxTrace{ID: id}
	rs, err := s.DB.NamedExec(sqlDelete, &m)
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
func (s *sqlserver) GetByID(id int64) (*TxTrace, error) {
	const sqlGetByID = `SELECT  id , type_message, file_name, code_line, message, convert(nvarchar(50), transaction_id) transaction_id, created_at FROM tx.trace  WITH (NOLOCK)  WHERE id = @id `
	mdl := TxTrace{}
	err := s.DB.Get(&mdl, sqlGetByID, sql.Named("id", id))
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
func (s *sqlserver) GetAll() ([]*TxTrace, error) {
	var ms []*TxTrace
	const sqlGetAll = `SELECT id , type_message, file_name, code_line, message, convert(nvarchar(50), transaction_id) transaction_id, created_at FROM tx.trace  WITH (NOLOCK) `

	err := s.DB.Select(&ms, sqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetAll trace: %v", err)
		return ms, err
	}
	return ms, nil
}

func (s *sqlserver) DeleteByTypeMsg(typeMsg string) error {
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

func (s *sqlserver) GetByTypeMsg(typeMsg string) ([]*TxTrace, error) {
	const sqlGetByTypeMsg = `SELECT id , type_message, file_name, code_line, message, convert(nvarchar(50), transaction_id) transaction_id, created_at FROM tx.trace  WITH (NOLOCK)  WHERE type_message = @typeMsg `
	var ms []*TxTrace
	err := s.DB.Select(&ms, sqlGetByTypeMsg, sql.Named("typeMsg", typeMsg))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByTypeMsg TxTrace: %v", err)
		return ms, err
	}
	return ms, nil
}

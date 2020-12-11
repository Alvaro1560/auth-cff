package loggedusers

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

func NewTxLoggedUserSqlServerRepository(db *sqlx.DB, user *models.User, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *sqlserver) Create(m *TxLoggedUser) error {
	const sqlInsert = `INSERT INTO tx.loggedusers ([event], host_name, ip_request, ip_remote, user_id) VALUES (@event, @host_name, @ip_request, @ip_remote, @user_id)  SELECT ID = convert(bigint, SCOPE_IDENTITY()) `

	stmt, err := s.DB.Prepare(sqlInsert)
	if err != nil {
		return nil
	}
	err = stmt.QueryRow(
		sql.Named("event", m.Event),
		sql.Named("host_name", m.HostName),
		sql.Named("ip_request", m.IpRequest),
		sql.Named("ip_remote", m.IpRemote),
		sql.Named("user_id", m.UserId),
	).Scan(&m.ID)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert tx loggedusers: %v", err)
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *sqlserver) Update(m *TxLoggedUser) error {

	const sqlUpdate = `UPDATE tx.loggedusers SET [event] = :event, host_name = :host_name, ip_request = :ip_request, ip_remote = :ip_remote, user_id = :user_id WHERE id = :id `

	rs, err := s.DB.NamedExec(sqlUpdate, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't update TxLoggedUser: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Delete elimina un registro de la BD
func (s *sqlserver) Delete(id int64) error {
	const sqlDelete = `DELETE FROM tx.loggedusers WHERE id = :id `
	m := TxLoggedUser{ID: id}
	rs, err := s.DB.NamedExec(sqlDelete, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't delete TxLoggedUser: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *sqlserver) GetByID(id int64) (*TxLoggedUser, error) {
	const sqlGetByID = `SELECT id , [event], host_name, ip_request, ip_remote, user_id, created_at FROM tx.loggedusers  WITH (NOLOCK)  WHERE id = @id `
	mdl := TxLoggedUser{}
	err := s.DB.Get(&mdl, sqlGetByID, sql.Named("id", id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByID TxLoggedUser: %v", err)
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *sqlserver) GetAll() ([]*TxLoggedUser, error) {
	var ms []*TxLoggedUser
	const sqlGetAll = `SELECT id , [event] , host_name, ip_request, ip_remote,convert(nvarchar(50), user_id) user_id, created_at FROM tx.loggedusers  WITH (NOLOCK) `

	err := s.DB.Select(&ms, sqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetAll loggedusers: %v", err)
		return ms, err
	}
	return ms, nil
}

func (s *sqlserver) GetByUserID(id string) ([]*TxLoggedUser, error) {
	var ms []*TxLoggedUser
	const sqlGetByUserID = `SELECT  id , [event] , host_name, ip_request, ip_remote,convert(nvarchar(50), user_id) user_id, created_at FROM tx.loggedusers  WITH (NOLOCK)  WHERE user_id = @id  `

	err := s.DB.Select(&ms, sqlGetByUserID, sql.Named("id", id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByUserID loggedusers: %v", err)
		return ms, err
	}
	return ms, nil
}

package loggedusers

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

func NewTxLoggedUserPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) Create(m *TxLoggedUser) error {
	const psqlInsert = `INSERT INTO tx.loggedusers (event, host_name, ip_request, ip_remote, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id `
	stmt, err := s.DB.Prepare(psqlInsert)
	if err != nil {
		return nil
	}
	err = stmt.QueryRow(
		sql.Named("id_execution", m.Event),
		sql.Named("id_execution", m.HostName),
		sql.Named("id_execution", m.IpRequest),
		sql.Named("id_execution", m.IpRemote),
		sql.Named("id_execution", m.UserId),
	).Scan(&m.ID)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert tx loggeduser: %v", err)
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) Update(m *TxLoggedUser) error {
	const psqlUpdate = `UPDATE tx.loggedusers SET event = :event, host_name = :host_name, ip_request = :ip_request, ip_remote = :ip_remote, user_id = :user_id, updated_at = now() WHERE id = :id `
	rs, err := s.DB.NamedExec(psqlUpdate, &m)
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
func (s *psql) Delete(id int64) error {
	const psqlDelete = `DELETE FROM tx.loggedusers WHERE id = :id `
	m := TxLoggedUser{ID: id}
	rs, err := s.DB.NamedExec(psqlDelete, &m)
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
func (s *psql) GetByID(id int64) (*TxLoggedUser, error) {
	const psqlGetByID = `SELECT id , event, host_name, ip_request, ip_remote, user_id, created_at, updated_at FROM tx.loggedusers WHERE id = $1 `
	mdl := TxLoggedUser{}
	err := s.DB.Get(&mdl, psqlGetByID, id)
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
func (s *psql) GetAll() ([]*TxLoggedUser, error) {
	var ms []*TxLoggedUser
	const psqlGetAll = ` SELECT id , event, host_name, ip_request, ip_remote, user_id, created_at, updated_at FROM tx.loggedusers `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetAll TxLoggedUser: %v", err)
		return ms, err
	}
	return ms, nil
}

func (s *psql) GetByUserID(id string) ([]*TxLoggedUser, error) {
	var ms []*TxLoggedUser
	const sqlGetByUserID = `SELECT id , event, host_name, ip_request, ip_remote, user_id, created_at, updated_at FROM tx.loggedusers  WITH (NOLOCK) WHERE user_id = $1 `

	err := s.DB.Select(&ms, sqlGetByUserID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByUserID loggedusers: %v", err)
		return ms, err
	}
	return ms, nil
}

package roles

import (
	"database/sql"
	"fmt"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/helper"

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

func NewRolePsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) Create(m *Role) error {
	m.UserID = s.user.ID
	const sqlInsert = `INSERT INTO auth.roles (id ,name, description, sessions_allowed,created_at, updated_at, id_user, is_delete) VALUES (:id ,:name, :description, :sessions_allowed, Now(), Now(), :id_user, false) `
	_, err := s.DB.NamedExec(sqlInsert, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert Role: %v", err)
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) Update(m *Role) error {
	const sqlUpdate = `UPDATE auth.roles SET name = :name, description = :description, sessions_allowed = :sessions_allowed, updated_at = Now() WHERE id = :id `
	rs, err := s.DB.NamedExec(sqlUpdate, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't update Role: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Delete elimina un registro de la BD
func (s *psql) Delete(id string) error {
	const sqlDelete = `DELETE FROM auth.roles WHERE id = :id `
	m := Role{ID: id}
	rs, err := s.DB.NamedExec(sqlDelete, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't delete Role: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *psql) GetByID(id string) (*Role, error) {
	const sqlGetByID = `SELECT id , name, description, sessions_allowed, created_at, updated_at FROM auth.roles  WHERE id = $1 `
	mdl := Role{}
	err := s.DB.Get(&mdl, sqlGetByID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByID Role: %v", err)
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *psql) GetAll() ([]*Role, error) {
	var ms []*Role
	const sqlGetAll = `SELECT id , name, description, sessions_allowed, created_at, updated_at FROM auth.roles `

	err := s.DB.Select(&ms, sqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetAll auth.roles: %v", err)
		return ms, err
	}
	return ms, nil
}

// GetByUserID consulta un registro por su ID
func (s *psql) GetByUserID(userID string) ([]*Role, error) {
	var ms []*Role
	const sqlGetByID = `SELECT r.id , r.name, r.description, r.sessions_allowed, r.created_at, r.updated_at FROM auth.roles r 
			    		JOIN auth.users_roles ur ON r.id = ur.role_id WHERE ur.user_id = $1 `
	err := s.DB.Select(&ms, sqlGetByID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByUserID auth.roles: %v", err)
		return ms, err
	}
	return ms, nil
}

func (s *psql) GetRolesByProcessIDs(ProcessIDs []string) ([]*Role, error) {
	var ms []*Role
	const sqlGetRolesByProcessIDs = `SELECT a.id, a.name, a.description, a.sessions_allowed, b.process_id, a.created_at, a.updated_at FROM auth.roles a JOIN cfg.process_roles b ON  a.id = b.role_id WHERE b.process_id in (%s)`

	err := s.DB.Select(&ms, fmt.Sprintf(sqlGetRolesByProcessIDs, helper.SliceToString(ProcessIDs)))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetRolesByProcessIDs auth.roles: %v", err)
		return ms, err
	}
	return ms, nil
}

func (s *psql) GetRolesByIDs(IDs []string) ([]*Role, error) {
	var ms []*Role
	const sqlGetRolesByIDs = `SELECT a.id, a.name, a.description, a.sessions_allowed, a.created_at, a.updated_at FROM auth.roles a WHERE a.id in (%s)`

	err := s.DB.Select(&ms, fmt.Sprintf(sqlGetRolesByIDs, helper.SliceToString(IDs)))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetRolesByIDs auth.roles: %v", err)
		return ms, err
	}
	return ms, nil
}

func (s *psql) GetByUserIDs(userIDs []string) ([]*Role, error) {
	var ms []*Role
	const sqlGetByID = `SELECT r.id, r.name, r.description, r.sessions_allowed, ur.id_user, r.created_at, r.updated_at  FROM auth.roles r JOIN auth.users_roles ur ON r.id = ur.role_id WHERE ur.id_user IN (%s) `
	query := fmt.Sprintf(sqlGetByID, helper.SliceToString(userIDs))
	err := s.DB.Select(&ms, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByUserIDs auth.roles: %v", err)
		return ms, err
	}
	return ms, nil
}

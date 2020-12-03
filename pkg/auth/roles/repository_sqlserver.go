package roles

import (
	"database/sql"
	"fmt"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/helper"

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

func NewRoleSqlServerRepository(db *sqlx.DB, user *models.User, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *sqlserver) Create(m *Role) error {
	const sqlInsert = `INSERT INTO auth.roles (id ,name, description, sessions_allowed,created_at, updated_at) VALUES (:id ,:name, :description, :sessions_allowed,GetDate(), GetDate()) `
	_, err := s.DB.NamedExec(sqlInsert, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert Role: %v", err)
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *sqlserver) Update(m *Role) error {
	const sqlUpdate = `UPDATE auth.roles SET name = :name, description = :description, sessions_allowed = :sessions_allowed, updated_at = GetDate() WHERE id = :id `
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
func (s *sqlserver) Delete(id string) error {
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
func (s *sqlserver) GetByID(id string) (*Role, error) {
	const sqlGetByID = `SELECT convert(nvarchar(50), id) id , name, description, sessions_allowed, created_at, updated_at FROM auth.roles  WITH (NOLOCK)  WHERE id = @id `
	mdl := Role{}
	err := s.DB.Get(&mdl, sqlGetByID, sql.Named("id", id))
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
func (s *sqlserver) GetAll() ([]*Role, error) {
	var ms []*Role
	const sqlGetAll = `SELECT convert(nvarchar(50), id) id , name, description, sessions_allowed, created_at, updated_at FROM auth.roles  WITH (NOLOCK) `

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
func (s *sqlserver) GetByUserID(userID string) ([]*Role, error) {
	var ms []*Role
	const sqlGetByID = `SELECT convert(nvarchar(50), r.id) id , r.name, r.description, r.sessions_allowed, r.created_at, r.updated_at FROM auth.roles r  WITH (NOLOCK) 
			    		JOIN auth.users_roles ur ON r.id = ur.role_id WHERE ur.user_id = @user_id `
	err := s.DB.Select(&ms, sqlGetByID, sql.Named("user_id", userID))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByUserID auth.roles: %v", err)
		return ms, err
	}
	return ms, nil
}

func (s *sqlserver) GetRolesByProcessIDs(ProcessIDs []string) ([]*Role, error) {
	var ms []*Role
	const sqlGetRolesByProcessIDs = `SELECT CONVERT(NVARCHAR(50), a.id) id , a.name, a.description, a.sessions_allowed, CONVERT(NVARCHAR(50), b.process_id) process_id, a.created_at, a.updated_at FROM auth.roles a WITH (NOLOCK) JOIN cfg.process_roles b WITH(NOLOCK) ON  a.id = b.role_id WHERE b.process_id in (%s)`

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

func (s *sqlserver) GetRolesByIDs(IDs []string) ([]*Role, error) {
	var ms []*Role
	const sqlGetRolesByIDs = `SELECT CONVERT(NVARCHAR(50), a.id) id , a.name, a.description, a.sessions_allowed, a.created_at, a.updated_at FROM auth.roles a WITH (NOLOCK)  WHERE a.id in (%s)`

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

func (s *sqlserver) GetByUserIDs(userIDs []string) ([]*Role, error) {
	var ms []*Role
	const sqlGetByID = `SELECT CONVERT(NVARCHAR(50), r.id) id, r.name, r.description, r.sessions_allowed, CONVERT(NVARCHAR(50), ur.[user_id]) [user_id], r.created_at, r.updated_at  FROM auth.roles r  WITH (NOLOCK) JOIN auth.users_roles ur ON r.id = ur.role_id WHERE ur.user_id IN (%s) `
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

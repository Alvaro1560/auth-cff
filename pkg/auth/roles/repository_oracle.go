package roles

import (
	"database/sql"
	"fmt"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/helper"

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

func NewRoleOrclRepository(db *sqlx.DB, user *models.User, txID string) *orcl {
	return &orcl{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *orcl) Create(m *Role) error {
	const osqlInsert = `INSERT INTO auth.roles (id ,name, description, sessions_allowed)  VALUES (:id ,:name, :description, :sessions_allowed)`
	_, err := s.DB.NamedExec(osqlInsert, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert Role: %v", err)
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *orcl) Update(m *Role) error {
	const osqlUpdate = `UPDATE auth.roles SET name = :name, description = :description, sessions_allowed = :sessions_allowed, updated_at = sysdate WHERE id = :id  `
	rs, err := s.DB.NamedExec(osqlUpdate, &m)
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
func (s *orcl) Delete(id string) error {
	const osqlDelete = `DELETE FROM auth.roles WHERE id = :id `
	m := Role{ID: id}
	rs, err := s.DB.NamedExec(osqlDelete, &m)
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
func (s *orcl) GetByID(id string) (*Role, error) {
	const osqlGetByID = `SELECT id , name, description, sessions_allowed, created_at, updated_at FROM auth.roles WHERE id = :1 `
	mdl := Role{}
	err := s.DB.Get(&mdl, osqlGetByID, id)
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
func (s *orcl) GetAll() ([]*Role, error) {
	var ms []*Role
	const osqlGetAll = ` SELECT id , name, description, sessions_allowed, created_at, updated_at FROM auth.roles `

	err := s.DB.Select(&ms, osqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetAll Role: %v", err)
		return ms, err
	}
	return ms, nil
}

// GetByUserID consulta un registro por su ID
func (s *orcl) GetByUserID(userID string) ([]*Role, error) {
	var ms []*Role
	const osqlGetByID = `SELECT r.id, r.name, r.description, r.sessions_allowed, r.created_at, r.updated_at FROM auth.roles r 
			    		JOIN auth.users_roles ur ON r.id = ur.role_id WHERE ur.user_id = :1 `
	err := s.DB.Select(&ms, osqlGetByID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByUserID auth.roles: %v", err)
		return ms, err
	}
	return ms, nil
}

func (s *orcl) GetRolesByProcessIDs(ProcessIDs []string) ([]*Role, error) {
	var ms []*Role
	const osqlGetRolesByProcessIDs = `SELECT CONVERT(NVARCHAR(50), a.id) id , a.name, a.description, a.sessions_allowed, CONVERT(NVARCHAR(50), b.process_id) process_id, a.created_at, a.updated_at FROM auth.roles a WITH (NOLOCK) JOIN cfg.process_roles b WITH(NOLOCK) ON  a.id = b.role_id WHERE b.process_id in (%s)`

	err := s.DB.Select(&ms, fmt.Sprintf(osqlGetRolesByProcessIDs, helper.SliceToString(ProcessIDs)))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetRolesByProcessIDs auth.roles: %v", err)
		return ms, err
	}
	return ms, nil
}

func (s *orcl) GetRolesByIDs(IDs []string) ([]*Role, error) {
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

func (s *orcl) GetByUserIDs(userIDs []string) ([]*Role, error) {
	var ms []*Role
	const osqlGetByID = `SELECT convert(nvarchar(50), r.id) id , r.name, r.description, r.sessions_allowed, r.created_at, r.updated_at FROM auth.roles r  WITH (NOLOCK) 
			    		JOIN auth.users_roles ur ON r.id = ur.role_id WHERE ur.user_id IN (%s) `
	err := s.DB.Select(&ms, osqlGetByID, helper.SliceToString(userIDs))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByUserIDs auth.roles: %v", err)
		return ms, err
	}
	return ms, nil
}

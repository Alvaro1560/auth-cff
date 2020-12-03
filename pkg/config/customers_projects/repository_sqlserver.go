package customers_projects

import (
	"database/sql"
	"fmt"
	"strings"

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

func NewProjectSqlServerRepository(db *sqlx.DB, user *models.User, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *sqlserver) Create(m *Project) error {
	const sqlInsert = `INSERT INTO cfg.customers_projects (id, name, description, department, email, phone, product_owner, customers_id,created_at, updated_at) VALUES (:id, :name, :description, :department, :email, :phone, :product_owner, :customers_id,GetDate(), GetDate()) `
	_, err := s.DB.NamedExec(sqlInsert, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert Project: %v", err)
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *sqlserver) Update(m *Project) error {
	const sqlUpdate = `UPDATE cfg.customers_projects SET name = :name, description = :description, department = :department, email = :email, phone = :phone, product_owner = :product_owner, customers_id = :customers_id, updated_at = GetDate() WHERE id = :id `
	rs, err := s.DB.NamedExec(sqlUpdate, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't update Project: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Delete elimina un registro de la BD
func (s *sqlserver) Delete(id string) error {
	const sqlDelete = `DELETE FROM cfg.customers_projects WHERE id = :id `
	m := Project{ID: id}
	rs, err := s.DB.NamedExec(sqlDelete, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't delete Project: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *sqlserver) GetByID(id string) (*Project, error) {
	const sqlGetByID = `SELECT convert(nvarchar(50), id) id, name, description, department, email, phone, product_owner, customers_id, created_at, updated_at FROM cfg.customers_projects  WITH (NOLOCK)  WHERE id = @id `
	mdl := Project{}
	err := s.DB.Get(&mdl, sqlGetByID, sql.Named("id", id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByID Project: %v", err)
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *sqlserver) GetAll() ([]*Project, error) {
	var ms []*Project
	const sqlGetAll = `SELECT convert(nvarchar(50), id) id, name, description, department, email, phone, product_owner, customers_id, created_at, updated_at FROM cfg.customers_projects  WITH (NOLOCK) `

	err := s.DB.Select(&ms, sqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetAll cfg.customers_projects: %v", err)
		return ms, err
	}
	return ms, nil
}

func (s *sqlserver) getProjectByRoles(roleIDs []string) ([]*string, error) {
	var ms []*Project
	var projects []*string
	const sqlGetAll = `SELECT convert(nvarchar(50), p.id) id  FROM cfg.customers_projects p WITH(NOLOCK) JOIN auth.roles_projects rp ON p.id = rp.project 
			    WHERE rp.role_id IN (@roles) `

	err := s.DB.Select(&ms, sqlGetAll, sql.Named("roles", strings.Join(roleIDs, ",")))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetProjectByRoles %v", err)
		return projects, err
	}
	for _, p := range ms {
		projects = append(projects, &p.ID)
	}
	return projects, nil
}

func (s *sqlserver) getProjectByRolesAndProjectID(roleIDs, pjts []string) ([]*Project, error) {
	var ms []*Project
	const sqlGetProjectByRolesAndProjectID = `SELECT convert(nvarchar(50), p.id) id, p.name, p.description, p.department, p.email, p.phone, p.product_owner, convert(nvarchar(50), p.customers_id) customers_id, p.created_at, p.updated_at  FROM cfg.customers_projects p WITH(NOLOCK) JOIN auth.roles_projects rp  WITH(NOLOCK) ON p.id = rp.project 
			    WHERE rp.role_id IN (@roles) AND p.customers_id IN (@projects) `

	err := s.DB.Select(&ms, sqlGetProjectByRolesAndProjectID, sql.Named("roles", strings.Join(roleIDs, ",")), sql.Named("projects", strings.Join(pjts, ",")))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute getProjectByRolesAndProjectID %v", err)
		return ms, err
	}

	return ms, nil
}

func (s *sqlserver) getProjectsByIds(projectIDs []string) ([]*Project, error) {
	var ms []*Project
	const sqlGetProjectsByIds = `SELECT convert(nvarchar(50), p.id) id, p.name, p.description, p.department, p.email, p.phone, p.product_owner, convert(nvarchar(50), p.customers_id) customers_id, p.created_at, p.updated_at  FROM cfg.customers_projects p WITH(NOLOCK)
			    WHERE  p.id IN (%s) `
	queryGetProjectsByIds := fmt.Sprintf(sqlGetProjectsByIds, helper.SliceToString(projectIDs))
	err := s.DB.Select(&ms, queryGetProjectsByIds)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute getProjectsByIds %v", err)
		return ms, err
	}

	return ms, nil
}
func (s *sqlserver) GetProjectByRoleIDs(RoleIDs []string) ([]*Project, error) {
	var ms []*Project
	const sqlGetProjectsByRoleIDs = `SELECT CONVERT(NVARCHAR(50),b.id) id, CONVERT(NVARCHAR(50),b.customers_id) customers_id, b.name, b.description, b.department, b.email, b.phone, b.product_owner, CONVERT(NVARCHAR(50),a.role_id) role_id,  b.created_at, b.updated_at FROM auth.roles_projects a WITH (NOLOCK) JOIN [cfg].[customers_projects] b WITH(NOLOCK) ON a.project = b.id WHERE a.role_id IN (%s)`

	err := s.DB.Select(&ms, fmt.Sprintf(sqlGetProjectsByRoleIDs, helper.SliceToString(RoleIDs)))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetProjectByRoleIDs auth.roles_projects: %v", err)
		return ms, err
	}
	return ms, nil
}

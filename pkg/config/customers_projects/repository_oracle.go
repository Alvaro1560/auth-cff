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

// Orcl estructura de conexión a la BD de Oracle
type orcl struct {
	DB   *sqlx.DB
	user *models.User
	TxID string
}

func NewProjectOrclRepository(db *sqlx.DB, user *models.User, txID string) *orcl {
	return &orcl{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *orcl) Create(m *Project) error {
	const osqlInsert = `INSERT INTO cfg.customers_projects (id , name, description, department, email, phone, product_owner, customers_id)  VALUES (:id, :name, :description, :department, :email, :phone, :product_owner, :customers_id)`
	_, err := s.DB.NamedExec(osqlInsert, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert Project: %v", err)
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *orcl) Update(m *Project) error {
	const osqlUpdate = `UPDATE cfg.customers_projects SET name = :name, description = :description, department = :department, email = :email, phone = :phone, product_owner = :product_owner, customers_id = :customers_id, updated_at = sysdate WHERE id = :id  `
	rs, err := s.DB.NamedExec(osqlUpdate, &m)
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
func (s *orcl) Delete(id string) error {
	const osqlDelete = `DELETE FROM cfg.customers_projects WHERE id = :id `
	m := Project{ID: id}
	rs, err := s.DB.NamedExec(osqlDelete, &m)
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
func (s *orcl) GetByID(id string) (*Project, error) {
	const osqlGetByID = `SELECT id, name, description, department, email, phone, product_owner, customers_id, created_at, updated_at FROM cfg.customers_projects WHERE id = :1 `
	mdl := Project{}
	err := s.DB.Get(&mdl, osqlGetByID, id)
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
func (s *orcl) GetAll() ([]*Project, error) {
	var ms []*Project
	const osqlGetAll = ` SELECT id, name, description, department, email, phone, product_owner, customers_id, created_at, updated_at FROM cfg.customers_projects `

	err := s.DB.Select(&ms, osqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetAll Project: %v", err)
		return ms, err
	}
	return ms, nil
}

func (s *orcl) getProjectByRoles(roleIDs []string) ([]*string, error) {
	var ms []*Project
	var projects []*string
	const sqlGetProjectByRoles = `SELECT p.id  FROM cfg.customers_projects p JOIN auth.roles_projects rp ON p.id = rp.project 
			    WHERE rp.role_id IN (%s) `
	query := fmt.Sprintf(sqlGetProjectByRoles, helper.SliceToString(roleIDs))
	err := s.DB.Select(&ms, query)
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

func (s *orcl) getProjectByRolesAndProjectID(roleIDs, pjts []string) ([]*Project, error) {
	var ms []*Project
	const sqlGetProjectByRolesAndProjectID = `SELECT p.id id, p.name, p.description, p.department, p.email, p.phone, p.product_owner,  p.customers_id customers_id, p.created_at, p.updated_at  FROM cfg.customers_projects p JOIN auth.roles_projects rp  ON p.id = rp.project 
			    WHERE rp.role_id IN (%s) AND p.customers_id IN (%s) `
	query := fmt.Sprintf(sqlGetProjectByRolesAndProjectID, helper.SliceToString(roleIDs), helper.SliceToString(pjts))
	err := s.DB.Select(&ms, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute getProjectByRolesAndProjectID %v", err)
		return ms, err
	}

	return ms, nil
}

func (s *orcl) getProjectsByIds(projectIDs []string) ([]*Project, error) {
	var ms []*Project
	const osqlGetProjectsByIds = `SELECT p.id , p.name, p.description, p.department, p.email, p.phone, p.product_owner, p.customers_id, p.created_at, p.updated_at  FROM cfg.customers_projects p 
			    WHERE  p.customers_id IN (:1) `

	err := s.DB.Select(&ms, osqlGetProjectsByIds, strings.Join(projectIDs, ","))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute getProjectsByIds %v", err)
		return ms, err
	}

	return ms, nil
}
func (s *orcl) GetProjectByRoleIDs(RoleIDs []string) ([]*Project, error) {
	var ms []*Project
	const osqlGetProjectsByRoleIDs = `SELECT CONVERT(NVARCHAR(50),b.id) id, CONVERT(NVARCHAR(50),b.customers_id) customers_id, b.name, b.description, b.department, b.email, b.phone, b.product_owner, CONVERT(NVARCHAR(50),a.role_id) role_id,  b.created_at, b.updated_at FROM auth.roles_projects a WITH (NOLOCK) JOIN [cfg].[customers_projects] b WITH(NOLOCK) ON a.project = b.id WHERE a.role_id IN (%s)`

	err := s.DB.Select(&ms, fmt.Sprintf(osqlGetProjectsByRoleIDs, helper.SliceToString(RoleIDs)))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetProjectByRoleIDs auth.roles_projects: %v", err)
		return ms, err
	}
	return ms, nil
}

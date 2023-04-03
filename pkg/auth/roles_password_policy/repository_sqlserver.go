package roles_password_policy

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

func NewRolesPasswordPolicySqlServerRepository(db *sqlx.DB, user *models.User, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *sqlserver) Create(m *RolesPasswordPolicy) error {
	const sqlInsert = `INSERT INTO auth.roles_password_policy (id ,role_id, days_pass_valid, max_length, min_length, store_pass_not_repeated, failed_attempts, time_unlock, alpha, digits, special, upper_case, lower_case, enable, inactivity_time, timeout,created_at, updated_at, required_2fa) VALUES (:id ,:role_id, :days_pass_valid, :max_length, :min_length, :store_pass_not_repeated, :failed_attempts, :time_unlock, :alpha, :digits, :special, :upper_case, :lower_case, :enable, :inactivity_time, :timeout,GetDate(), GetDate(), :required_2fa) `
	_, err := s.DB.NamedExec(sqlInsert, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert RolesPasswordPolicy: %v", err)
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *sqlserver) Update(m *RolesPasswordPolicy) error {
	const sqlUpdate = `UPDATE auth.roles_password_policy SET role_id = :role_id, days_pass_valid = :days_pass_valid, max_length = :max_length, min_length = :min_length, store_pass_not_repeated = :store_pass_not_repeated, failed_attempts = :failed_attempts, time_unlock = :time_unlock, alpha = :alpha, digits = :digits, special = :special, upper_case = :upper_case, lower_case = :lower_case, enable = :enable, inactivity_time = :inactivity_time, timeout = :timeout, updated_at = GetDate(), required_2fa = :required_2fa WHERE id = :id `
	rs, err := s.DB.NamedExec(sqlUpdate, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't update RolesPasswordPolicy: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Delete elimina un registro de la BD
func (s *sqlserver) Delete(id string) error {
	const sqlDelete = `DELETE FROM auth.roles_password_policy WHERE id = :id `
	m := RolesPasswordPolicy{ID: id}
	rs, err := s.DB.NamedExec(sqlDelete, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't delete RolesPasswordPolicy: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *sqlserver) GetByID(id string) (*RolesPasswordPolicy, error) {
	const sqlGetByID = `SELECT convert(nvarchar(50), id) id , role_id, days_pass_valid, max_length, min_length, store_pass_not_repeated, failed_attempts, time_unlock, alpha, digits, special, upper_case, lower_case, enable, inactivity_time, timeout, required_2fa, created_at, updated_at FROM auth.roles_password_policy  WITH (NOLOCK)  WHERE id = @id `
	mdl := RolesPasswordPolicy{}
	err := s.DB.Get(&mdl, sqlGetByID, sql.Named("id", id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByID RolesPasswordPolicy: %v", err)
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *sqlserver) GetAll() ([]*RolesPasswordPolicy, error) {
	var ms []*RolesPasswordPolicy
	const sqlGetAll = `SELECT convert(nvarchar(50), id) id , role_id, days_pass_valid, max_length, min_length, store_pass_not_repeated, failed_attempts, time_unlock, alpha, digits, special, upper_case, lower_case, enable, inactivity_time, timeout, required_2fa, created_at, updated_at FROM auth.roles_password_policy  WITH (NOLOCK) `

	err := s.DB.Select(&ms, sqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetAll auth.roles_password_policy: %v", err)
		return ms, err
	}
	return ms, nil
}

func (s *sqlserver) GetAllRolesPasswordPolicyByRolesIDs(RolesIDs []string) ([]*RolesPasswordPolicy, error) {
	var ms []*RolesPasswordPolicy
	const sqlGetPasswordPolicyByRolesIDs = `SELECT convert(nvarchar(50), id) id , convert(nvarchar(50), role_id) role_id, days_pass_valid, max_length, min_length, store_pass_not_repeated, failed_attempts, time_unlock, alpha, digits, special, upper_case, lower_case, enable, inactivity_time, timeout, required_2fa, created_at, updated_at FROM auth.roles_password_policy  WITH (NOLOCK) WHERE role_id in (%s)`

	err := s.DB.Select(&ms, fmt.Sprintf(sqlGetPasswordPolicyByRolesIDs, helper.SliceToString(RolesIDs)))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetAllRolesPasswordPolicyByRolesIDs auth.roles_password_policy: %v", err)
		return ms, err
	}
	return ms, nil
}

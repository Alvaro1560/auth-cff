package roles_password_policy

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/helper"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
)

// psql estructura de conexión a la BD de postgresql
type psql struct {
	DB   *sqlx.DB
	user *models.User
	TxID string
}

func NewRolesPasswordPolicyPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) Create(m *RolesPasswordPolicy) error {
	m.IdUser = s.user.ID
	const sqlInsert = `INSERT INTO auth.roles_password_policy (id ,role_id, days_pass_valid, max_length, min_length, store_pass_not_repeated, failed_attempts, time_unlock, alpha, digits, special, upper_case, lower_case, enable, inactivity_time, timeout,created_at, updated_at, id_user, is_delete, required_2fa) VALUES (:id ,:role_id, :days_pass_valid, :max_length, :min_length, :store_pass_not_repeated, :failed_attempts, :time_unlock, :alpha, :digits, :special, :upper_case, :lower_case, :enable, :inactivity_time, :timeout,Now(), Now(), :id_user, false, :required_2fa) `
	_, err := s.DB.NamedExec(sqlInsert, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert RolesPasswordPolicy: %v", err)
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) Update(m *RolesPasswordPolicy) error {
	const sqlUpdate = `UPDATE auth.roles_password_policy SET role_id = :role_id, days_pass_valid = :days_pass_valid, max_length = :max_length, min_length = :min_length, store_pass_not_repeated = :store_pass_not_repeated, failed_attempts = :failed_attempts, time_unlock = :time_unlock, alpha = :alpha, digits = :digits, special = :special, upper_case = :upper_case, lower_case = :lower_case, enable = :enable, inactivity_time = :inactivity_time, timeout = :timeout, updated_at = Now(), required_2fa = :required_2fa WHERE id = :id `
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
func (s *psql) Delete(id string) error {
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
func (s *psql) GetByID(id string) (*RolesPasswordPolicy, error) {
	const sqlGetByID = `SELECT  id , role_id, days_pass_valid, max_length, min_length, store_pass_not_repeated, failed_attempts, time_unlock, alpha, digits, special, upper_case, lower_case, enable, inactivity_time, timeout, required_2fa, created_at, updated_at FROM auth.roles_password_policy    WHERE id = $1 `
	mdl := RolesPasswordPolicy{}
	err := s.DB.Get(&mdl, sqlGetByID, id)
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
func (s *psql) GetAll() ([]*RolesPasswordPolicy, error) {
	var ms []*RolesPasswordPolicy
	const sqlGetAll = `SELECT  id , role_id, days_pass_valid, max_length, min_length, store_pass_not_repeated, failed_attempts, time_unlock, alpha, digits, special, upper_case, lower_case, enable, inactivity_time, timeout, required_2fa, created_at, updated_at FROM auth.roles_password_policy   `

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

func (s *psql) GetAllRolesPasswordPolicyByRolesIDs(RolesIDs []string) ([]*RolesPasswordPolicy, error) {
	var ms []*RolesPasswordPolicy
	const sqlGetPasswordPolicyByRolesIDs = `SELECT  id ,  role_id, days_pass_valid, max_length, min_length, store_pass_not_repeated, failed_attempts, time_unlock, alpha, digits, special, upper_case, lower_case, enable, inactivity_time, timeout, required_2fa, created_at, updated_at FROM auth.roles_password_policy   WHERE role_id in (%s)`

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

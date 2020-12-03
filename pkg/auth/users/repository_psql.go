package users

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

func NewUserPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) Create(m *User) error {
	const sqlInsert = `INSERT INTO auth.users (id ,username, name, lastname, password, email_notifications, identification_number, identification_type, status, failed_attempts, change_password, is_block, is_disabled,created_at, updated_at) 
					  VALUES (:id ,:username, :name, :lastname, :password, :email_notifications, :identification_number, :identification_type, :status, :failed_attempts, :change_password, :is_block, :is_disabled,Now(), Now()) `
	rs, err := s.DB.NamedExec(sqlInsert, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert User: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) Update(m *User) error {
	const sqlUpdate = `UPDATE auth.users SET username = :username, name = :name, lastname = :lastname, password = :password, email_notifications = :email_notifications, identification_number = :identification_number, identification_type = :identification_type, status = :status, failed_attempts = :failed_attempts, last_change_password = :last_change_password, block_date = :block_date, disabled_date = :disabled_date, change_password = :change_password, is_block = :is_block, is_disabled = :is_disabled, last_login = :last_login, updated_at = Now() WHERE id = :id `
	rs, err := s.DB.NamedExec(sqlUpdate, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't update User: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Delete elimina un registro de la BD
func (s *psql) Delete(id string) error {
	const sqlDelete = `DELETE FROM auth.users WHERE id = :id `
	m := User{ID: id}
	rs, err := s.DB.NamedExec(sqlDelete, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't delete User: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *psql) GetByID(id string) (*User, error) {
	const sqlGetByID = `SELECT id , username, name, lastname, password, email_notifications, identification_number, identification_type, status, failed_attempts, last_change_password, block_date, disabled_date, change_password, is_block, is_disabled, last_login, created_at, updated_at FROM auth.users    WHERE id = $1 `
	mdl := User{}
	err := s.DB.Get(&mdl, sqlGetByID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByID User: %v", err)
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *psql) GetAll() ([]*User, error) {
	var ms []*User
	const sqlGetAll = `SELECT id , username, name, lastname, password, email_notifications, identification_number, identification_type, status, failed_attempts, last_change_password, block_date, disabled_date, change_password, is_block, is_disabled, last_login, created_at, updated_at FROM auth.users   `

	err := s.DB.Select(&ms, sqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetAll auth.users: %v", err)
		return ms, err
	}
	return ms, nil
}

// GetByUsername consulta un registro por su ID
func (s *psql) GetByUsername(username string) (*User, error) {
	const sqlGetByUsername = `SELECT id , username, name, lastname, password, email_notifications, identification_number, identification_type, status, failed_attempts, last_change_password, block_date, disabled_date, change_password, is_block, is_disabled, last_login, created_at, updated_at FROM auth.users WHERE username = $1 `
	mdl := User{}
	err := s.DB.Get(&mdl, sqlGetByUsername, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByUsername User: %v", err)
		return &mdl, err
	}
	return &mdl, nil
}

// GetUsersByIDs consulta todos los registros de la BD
func (s *psql) GetUsersByIDs(ids []string) ([]*User, error) {
	var ms []*User
	const sqlGetUsersByIDs = `SELECT id , username, name, lastname, password, email_notifications, identification_number, identification_type, status, failed_attempts, last_change_password, block_date, disabled_date, change_password, is_block, is_disabled, last_login, created_at, updated_at FROM auth.users   WHERE id IN (%s) `
	querySqlGetUsersByIDs := fmt.Sprintf(sqlGetUsersByIDs, helper.SliceToString(ids))
	err := s.DB.Select(&ms, querySqlGetUsersByIDs)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute sqlGetUsersByIDs auth.users: %v", err)
		return ms, err
	}
	return ms, nil
}

//Bloquea el Usuario
func (s *psql) BlockUser(id string) error {
	const sqlUpdateBlockUser = `UPDATE auth.users SET status = 16, block_date = Now(), is_block = 1, updated_at = Now() WHERE id = :id `
	m := User{ID: id}
	rs, err := s.DB.NamedExec(sqlUpdateBlockUser, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't BlockUser User: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

//Desbloquea el Usuario
func (s *psql) UnblockUser(id string) error {
	const sqlUpdateUnblockUser = `UPDATE auth.users SET status = 1, block_date = Now(), is_block = 0, updated_at = Now() WHERE id = :id `
	m := User{ID: id}
	rs, err := s.DB.NamedExec(sqlUpdateUnblockUser, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't UnblockUser User: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

//Cierra Sesion del Usuario
func (s *psql) LogoutUser(id string) error {
	const sqlLogoutUser = `DELETE FROM auth.users_loggeds WHERE user_id = :id `
	m := User{ID: id}
	rs, err := s.DB.NamedExec(sqlLogoutUser, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't LogoutUser User: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

//Cambia la clave y guarda el historial de clave
func (s *psql) ChangePassword(id string, password string) error {
	tx, err := s.DB.Beginx()
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't begin Tx to ChangePassword: %v", err)
		return err
	}

	const sqlUpdate = `UPDATE auth.users SET password = :password, last_change_password = Now(), updated_at = Now() WHERE id = :id `
	m := User{ID: id, Password: password}

	rs, err := tx.NamedExec(sqlUpdate, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't update User: %v", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error.Printf(s.TxID, " - UpdatePassword: unable to rollback: %v", rollbackErr)
		}
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	const sqlInsert = `INSERT INTO auth.users_password_history (id ,user_id, password,created_at) VALUES (uuid_generate_v4() , :id, :password, Now()) `

	_, err = tx.NamedExec(sqlInsert, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert UsersPasswordHistory: %v", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error.Printf(s.TxID, " - insert UsersPasswordHistory: unable to rollback: %v", rollbackErr)
		}
		return err
	}
	if commitErr := tx.Commit(); commitErr != nil {
		logger.Error.Printf(s.TxID, " - ChangePassword: unable to commit: %v", commitErr)
	}
	return nil
}

func (s *psql) UpdatePasswordByUser(id string, password string) error {
	tx, err := s.DB.Beginx()
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't begin Tx to UpdatePasswordByUser: %v", err)
		return err
	}

	const sqlUpdate = `UPDATE auth.users SET password = :password, last_change_password = Now(), updated_at = Now() WHERE id = :id `
	m := User{ID: id, Password: password}

	rs, err := tx.NamedExec(sqlUpdate, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't update User: %v", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error.Printf(s.TxID, " - UpdatePasswordByUser: unable to rollback: %v", rollbackErr)
		}
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	const sqlInsert = `INSERT INTO auth.users_password_history (id ,user_id, password,created_at) VALUES (uuid_generate_v4() , :id, :password, Now()) `

	_, err = tx.NamedExec(sqlInsert, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert UsersPasswordHistory: %v", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error.Printf(s.TxID, " - insert UsersPasswordHistory: unable to rollback: %v", rollbackErr)
		}
		return err
	}
	if commitErr := tx.Commit(); commitErr != nil {
		logger.Error.Printf(s.TxID, " - UpdatePasswordByUser: unable to commit: %v", commitErr)
	}
	return nil
}

func (s *psql) GetByUsernameAndIdentificationNumber(username string, identificationNumber string) (*User, error) {
	const sqlGetByUsername = `SELECT  id , username, name, lastname, password, email_notifications, identification_number, identification_type, status, failed_attempts, last_change_password, block_date, disabled_date, change_password, is_block, is_disabled, last_login, created_at, updated_at FROM auth.users  WHERE username = '$1' AND identification_number = '$2'`
	mdl := User{}
	err := s.DB.Get(&mdl, sqlGetByUsername, username, identificationNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByUsername User: %v", err)
		return &mdl, err
	}
	return &mdl, nil
}

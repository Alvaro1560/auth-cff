package users

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

func NewUserOrclRepository(db *sqlx.DB, user *models.User, txID string) *orcl {
	return &orcl{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *orcl) Create(m *User) error {
	const osqlInsert = `INSERT INTO auth.users (id ,user, name, lastname, password, email_notifications, identification_number, identification_type, status, failed_attempts, last_change_password, block_date, disabled_date, change_password, is_block, is_disabled, last_login)  VALUES (:id ,:user, :name, :lastname, :password, :email_notifications, :identification_number, :identification_type, :status, :failed_attempts, :last_change_password, :block_date, :disabled_date, :change_password, :is_block, :is_disabled, :last_login)`
	_, err := s.DB.NamedExec(osqlInsert, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert User: %v", err)
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *orcl) Update(m *User) error {
	const osqlUpdate = `UPDATE auth.users SET user = :user, name = :name, lastname = :lastname, password = :password, email_notifications = :email_notifications, identification_number = :identification_number, identification_type = :identification_type, status = :status, failed_attempts = :failed_attempts, last_change_password = :last_change_password, block_date = :block_date, disabled_date = :disabled_date, change_password = :change_password, is_block = :is_block, is_disabled = :is_disabled, last_login = :last_login, updated_at = sysdate WHERE id = :id  `
	rs, err := s.DB.NamedExec(osqlUpdate, &m)
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
func (s *orcl) Delete(id string) error {
	const osqlDelete = `DELETE FROM auth.users WHERE id = :id `
	m := User{ID: id}
	rs, err := s.DB.NamedExec(osqlDelete, &m)
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
func (s *orcl) GetByID(id string) (*User, error) {
	const osqlGetByID = `SELECT id , user, name, lastname, password, email_notifications, identification_number, identification_type, status, failed_attempts, last_change_password, block_date, disabled_date, change_password, is_block, is_disabled, last_login, created_at, updated_at FROM auth.users WHERE id = :1 `
	mdl := User{}
	err := s.DB.Get(&mdl, osqlGetByID, id)
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
func (s *orcl) GetAll() ([]*User, error) {
	var ms []*User
	const osqlGetAll = ` SELECT id , user, name, lastname, password, email_notifications, identification_number, identification_type, status, failed_attempts, last_change_password, block_date, disabled_date, change_password, is_block, is_disabled, last_login, created_at, updated_at FROM auth.users `

	err := s.DB.Select(&ms, osqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetAll User: %v", err)
		return ms, err
	}
	return ms, nil
}

// GetByUsername consulta un registro por su ID
func (s *orcl) GetByUsername(username string) (*User, error) {
	const osqlGetByUsername = `SELECT id , user, name, lastname, password, email_notifications, identification_number, identification_type, status, failed_attempts, last_change_password, block_date, disabled_date, change_password, is_block, is_disabled, last_login, created_at, updated_at FROM auth.users WHERE username = :1 `
	mdl := User{}
	err := s.DB.Get(&mdl, osqlGetByUsername, username)
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
func (s *orcl) GetUsersByIDs(ids []string) ([]*User, error) {
	var ms []*User
	const sqlGetUsersByIDs = `SELECT convert(nvarchar(50), id) id , username, name, lastname, password, email_notifications, identification_number, identification_type, status, failed_attempts, last_change_password, block_date, disabled_date, change_password, is_block, is_disabled, last_login, created_at, updated_at FROM auth.users  WITH (NOLOCK) WHERE id IN (%s) `

	err := s.DB.Select(&ms, fmt.Sprintf(sqlGetUsersByIDs, helper.SliceToString(ids)))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute sqlGetUsersByIDs auth.users: %v", err)
		return ms, err
	}
	return ms, nil
}

// Bloquea el Usuario
func (s *orcl) BlockUser(id string) error {
	const osqlUpdateBlockUser = `UPDATE auth.users SET [status] = 16, block_date = GETDATE(), is_block = 1, updated_at = GETDATE() WHERE id = :id `
	m := User{ID: id}
	rs, err := s.DB.NamedExec(osqlUpdateBlockUser, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't BlockUser User: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

func (s *orcl) UnblockUser(id string) error {
	const osqlUpdateUnblockUser = `UPDATE auth.users SET [status] = 1, block_date = GETDATE(), is_block = 0, updated_at = GETDATE() WHERE id = :id `
	m := User{ID: id}
	rs, err := s.DB.NamedExec(osqlUpdateUnblockUser, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't UnblockUser User: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

func (s *orcl) LogoutUser(id string) error {
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

func (s *orcl) ChangePassword(id string, password string) error {
	tx, err := s.DB.Beginx()
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't begin Tx to ChangePassword: %v", err)
		return err
	}

	const osqlUpdate = `UPDATE auth.users SET password = :password, last_change_password = GETDATE(), updated_at = GETDATE() WHERE id = :id `
	m := User{ID: id, Password: password}

	rs, err := tx.NamedExec(osqlUpdate, &m)
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
	const osqlInsert = `INSERT INTO auth.users_password_history (id ,[user_id], [password],created_at, updated_at) VALUES (NEWID() , :id, :password, GETDATE(), GETDATE()) `
	_, err = tx.NamedExec(osqlInsert, &m)
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

func (s *orcl) UpdatePasswordByUser(id string, password string) error {
	tx, err := s.DB.Beginx()
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't begin Tx to UpdatePasswordByUser: %v", err)
		return err
	}

	const osqlUpdate = `UPDATE auth.users SET password = :password, last_change_password = GETDATE(), updated_at = GETDATE() WHERE id = :id `
	m := User{ID: id, Password: password}

	rs, err := tx.NamedExec(osqlUpdate, &m)
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
	const osqlInsert = `INSERT INTO auth.users_password_history (id ,[user_id], [password],created_at) VALUES (NEWID() , :id, :password, GETDATE()) `

	_, err = tx.NamedExec(osqlInsert, &m)
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

func (s *orcl) GetByUsernameAndIdentificationNumber(username string, identificationNumber string) (*User, error) {
	const osqlGetByUsername = `SELECT convert(nvarchar(50), id) id , username, name, lastname, password, email_notifications, identification_number, identification_type, status, failed_attempts, last_change_password, block_date, disabled_date, change_password, is_block, is_disabled, last_login, created_at, updated_at FROM auth.users WITH(NOLOCK) WHERE username = '%s' AND identification_number = '%s'`
	mdl := User{}
	err := s.DB.Get(&mdl, fmt.Sprintf(osqlGetByUsername, username, identificationNumber))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByUsername User: %v", err)
		return &mdl, err
	}
	return &mdl, nil
}

// DeleteUserPasswordHistory elimina un registro de la BD
func (s *orcl) DeleteUserPasswordHistory(id string) error {
	const sqlDelete = `DELETE FROM auth.users_password_history WHERE user_id = :id `
	m := User{ID: id}
	rs, err := s.DB.NamedExec(sqlDelete, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't delete users_password_history: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

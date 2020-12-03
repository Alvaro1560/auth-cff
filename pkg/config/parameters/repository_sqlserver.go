package parameters

import (
	"database/sql"
	"fmt"

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

func NewParameterSqlServerRepository(db *sqlx.DB, user *models.User, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *sqlserver) Create(m *Parameter) error {
	var id int64
	const sqlInsert = `INSERT INTO cfg.parameters (name, [value], type, description, client_id, is_cipher,created_at, updated_at) VALUES (@name, @value, @type, @description, @client_id, @is_cipher,GetDate(), GetDate())   SELECT ID = convert(bigint, SCOPE_IDENTITY()) `
	stmt, err := s.DB.Prepare(sqlInsert)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert Parameters: %v", err)
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		sql.Named("name", m.Name),
		sql.Named("value", m.Value),
		sql.Named("type", m.Type),
		sql.Named("description", m.Description),
		sql.Named("client_id", m.ClientId),
		sql.Named("is_cipher", m.IsCipher),
	).Scan(&id)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't insert Parameters: %v", err)
		return err
	}
	m.ID = int(id)
	return nil
}

// Update actualiza un registro en la BD
func (s *sqlserver) Update(m *Parameter) error {
	const sqlUpdate = `UPDATE cfg.parameters SET name = :name, [value] = :value, type = :type, description = :description, client_id = :client_id, is_cipher = :is_cipher, updated_at = GetDate() WHERE id = :id `
	rs, err := s.DB.NamedExec(sqlUpdate, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't update Parameter: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Delete elimina un registro de la BD
func (s *sqlserver) Delete(id int) error {
	const sqlDelete = `DELETE FROM cfg.parameters WHERE id = :id `
	m := Parameter{ID: id}
	rs, err := s.DB.NamedExec(sqlDelete, &m)
	if err != nil {
		logger.Error.Printf(s.TxID, " - couldn't delete Parameter: %v", err)
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *sqlserver) GetByID(id int) (*Parameter, error) {
	const sqlGetByID = `SELECT convert(nvarchar(50), id) id ,name ,[value] , type, description, client_id, is_cipher, created_at, updated_at FROM cfg.parameters  WITH (NOLOCK)  WHERE id = @id `
	mdl := Parameter{}
	err := s.DB.Get(&mdl, sqlGetByID, sql.Named("id", id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByID Parameter: %v", err)
		return &mdl, err
	}
	return &mdl, nil
}

// GetByName consulta un registro por su name
func (s *sqlserver) GetByName(name string) (*Parameter, error) {
	const osqlGetByName = `SELECT id ,name ,[value] , type, description, client_id, is_cipher, created_at, updated_at FROM cfg.parameters WHERE name = @name `
	mdl := Parameter{}
	err := s.DB.Get(&mdl, osqlGetByName, sql.Named("name", name))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetByName Parameter: %v", err)
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *sqlserver) GetAll() ([]*Parameter, error) {
	var ms []*Parameter
	const sqlGetAll = `SELECT convert(nvarchar(50), id) id ,name ,[value] , type, description, client_id, is_cipher, created_at, updated_at FROM cfg.parameters  WITH (NOLOCK) `

	err := s.DB.Select(&ms, sqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error.Printf(s.TxID, " - couldn't execute GetAll cfg.parameters: %v", err)
		return ms, err
	}
	return ms, nil
}

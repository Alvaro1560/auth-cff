package verification_email

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
)

// sqlServer estructura de conexión a la BD de mssql
type sqlserver struct {
	DB   *sqlx.DB
	user *models.User
	TxID string
}

func newVerificationEmailSqlServerRepository(db *sqlx.DB, user *models.User, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *sqlserver) create(m *VerificationEmail) error {
	var id int64
	date := time.Now()
	m.UpdatedAt = date
	m.CreatedAt = date
	const sqlInsert = `INSERT INTO auth.verification_email (email, verification_code, identification, verification_date, created_at, updated_at) VALUES (@email, @verification_code, @identification, @verification_date, @created_at, @updated_at) SELECT ID = convert(bigint, SCOPE_IDENTITY()) `
	stmt, err := s.DB.Prepare(sqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		sql.Named("email", m.Email),
		sql.Named("verification_code", m.VerificationCode),
		sql.Named("identification", m.Identification),
		sql.Named("verification_date", m.VerificationDate),
		sql.Named("created_at", m.CreatedAt),
		sql.Named("updated_at", m.UpdatedAt),
	).Scan(&id)
	if err != nil {
		return err
	}
	m.ID = id
	return nil
}

// Update actualiza un registro en la BD
func (s *sqlserver) update(m *VerificationEmail) error {
	date := time.Now()
	m.UpdatedAt = date
	const sqlUpdate = `UPDATE auth.verification_email SET email = :email, identification = :identification, verification_date = :verification_date, updated_at = :updated_at WHERE id = :id `
	rs, err := s.DB.NamedExec(sqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Delete elimina un registro de la BD
func (s *sqlserver) delete(id int64) error {
	const sqlDelete = `DELETE FROM auth.verification_email WHERE id = :id `
	m := VerificationEmail{ID: id}
	rs, err := s.DB.NamedExec(sqlDelete, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *sqlserver) getByID(id int64) (*VerificationEmail, error) {
	const sqlGetByID = `SELECT convert(nvarchar(50), id) id , email, verification_code, identification, verification_date, created_at, updated_at FROM auth.verification_email  WITH (NOLOCK)  WHERE id = @id `
	mdl := VerificationEmail{}
	err := s.DB.Get(&mdl, sqlGetByID, sql.Named("id", id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *sqlserver) getAll() ([]*VerificationEmail, error) {
	var ms []*VerificationEmail
	const sqlGetAll = `SELECT convert(nvarchar(50), id) id , email, verification_code, identification, verification_date, created_at, updated_at FROM auth.verification_email  WITH (NOLOCK) `

	err := s.DB.Select(&ms, sqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

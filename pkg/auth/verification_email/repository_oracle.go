package verification_email

import (
	"database/sql"
	"fmt"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

// Orcl estructura de conexión a la BD de Oracle
type orcl struct {
	DB   *sqlx.DB
	user *models.User
	TxID string
}

func newVerificationEmailOrclRepository(db *sqlx.DB, user *models.User, txID string) *orcl {
	return &orcl{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *orcl) create(m *VerificationEmail) error {

	const osqlInsert = `INSERT INTO auth.verification_email (email, verification_code, identification, verification_date, created_at, updated_at)  VALUES (:email, :verification_code, :identification, :verification_date, sysdate, sysdate) RETURNING id into id, created_at into created_at, updated_at into updated_at `
	stmt, err := s.DB.Prepare(osqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		m.Email,
		m.VerificationCode,
		m.Identification,
		m.VerificationDate,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *orcl) update(m *VerificationEmail) error {
	date := time.Now()
	m.UpdatedAt = date
	const osqlUpdate = `UPDATE auth.verification_email SET email = :email, identification = :identification, verification_date = :verification_date, updated_at = :updated_at WHERE id = :id  `
	rs, err := s.DB.NamedExec(osqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Delete elimina un registro de la BD
func (s *orcl) delete(id int64) error {
	const osqlDelete = `DELETE FROM auth.verification_email WHERE id = :id `
	m := VerificationEmail{ID: id}
	rs, err := s.DB.NamedExec(osqlDelete, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *orcl) getByID(id int64) (*VerificationEmail, error) {
	const osqlGetByID = `SELECT id , email, verification_code, identification, verification_date, created_at, updated_at FROM auth.verification_email WHERE id = :1 `
	mdl := VerificationEmail{}
	err := s.DB.Get(&mdl, osqlGetByID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *orcl) getAll() ([]*VerificationEmail, error) {
	var ms []*VerificationEmail
	const osqlGetAll = ` SELECT id , email, verification_code, identification, verification_date, created_at, updated_at FROM auth.verification_email `

	err := s.DB.Select(&ms, osqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

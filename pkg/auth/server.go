package auth

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth/users"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth/verification_email"
)

type ServerAuth struct {
	SrvUsers             users.Service
	SrvVerificationEmail verification_email.PortsServerVerificationEmail
}

func NewServerAuth(db *sqlx.DB, usr *models.User, txID string) *ServerAuth {
	repoUsers := users.FactoryStorage(db, usr, txID)
	srvUsers := users.NewUserService(repoUsers, usr, txID)

	repoVerificationEmail := verification_email.FactoryStorage(db, usr, txID)
	srvVerificationEmail := verification_email.NewVerificationEmailService(repoVerificationEmail, usr, txID)

	return &ServerAuth{
		SrvUsers:             srvUsers,
		SrvVerificationEmail: srvVerificationEmail,
	}

}

package test

import (
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/dbx"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth/verification_email"
	"testing"
)

func TestCreateVerificationEmail(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := verification_email.FactoryStorage(db, nil, txID)
	srvRole := verification_email.NewVerificationEmailService(repoRole, nil, txID)

	_, cod, err := srvRole.CreateVerificationEmail("test@ecapture.co", "123", "12345678", nil)
	if err != nil {
		t.Fatalf("Role CreateVerificationEmail error, cod: %d error: %v", cod, err)
	}
	t.Log("Role CreateVerificationEmail was successful")
	fmt.Print("Role CreateVerificationEmail was successful")
}

func TestUpdateVerificationEmail(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := verification_email.FactoryStorage(db, nil, txID)
	srvRole := verification_email.NewVerificationEmailService(repoRole, nil, txID)

	_, cod, err := srvRole.UpdateVerificationEmail(2533, "test@ecapture.co", "123", "12345678", nil)
	if err != nil {
		t.Fatalf("Role UpdateVerificationEmail error, cod: %d error: %v", cod, err)
	}
	t.Log("Role UpdateVerificationEmail was successful")
	fmt.Print("Role UpdateVerificationEmail was successful")
}

func TestDeleteVerificationEmail(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := verification_email.FactoryStorage(db, nil, txID)
	srvRole := verification_email.NewVerificationEmailService(repoRole, nil, txID)

	cod, err := srvRole.DeleteVerificationEmail(123)
	if err != nil {
		t.Fatalf("Role DeleteVerificationEmail error, cod: %d error: %v", cod, err)
	}
	t.Log("Role DeleteVerificationEmail was successful")
	fmt.Print("Role DeleteVerificationEmail was successful")
}

func TestGetVerificationEmailByID(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := verification_email.FactoryStorage(db, nil, txID)
	srvRole := verification_email.NewVerificationEmailService(repoRole, nil, txID)

	_, cod, err := srvRole.GetVerificationEmailByID(123)
	if err != nil {
		t.Fatalf("Role GetVerificationEmailByID error, cod: %d error: %v", cod, err)
	}
	t.Log("Role GetVerificationEmailByID was successful")
	fmt.Print("Role GetVerificationEmailByID was successful")
}

func TestGetAllVerificationEmail(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := verification_email.FactoryStorage(db, nil, txID)
	srvRole := verification_email.NewVerificationEmailService(repoRole, nil, txID)

	_, err := srvRole.GetAllVerificationEmail()
	if err != nil {
		t.Fatalf("Role GetAllVerificationEmail error, error: %v", err)
	}
	t.Log("Role GetAllVerificationEmail was successful")
	fmt.Print("Role GetAllVerificationEmail was successful")
}

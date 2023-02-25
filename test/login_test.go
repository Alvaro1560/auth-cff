package test

import (
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/dbx"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth/login"
	"testing"
)

func TestLogin(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	serviceLogin := login.NewLoginService(db, txID)

	_, cod, err := serviceLogin.Login("U2FsdGVkX19FENMwumk63jrB6xoByxpwUpAMcgVyijo=", "", "U2FsdGVkX19DwbpampI9FBdXLVYIDbZMaW7iol0ThBM=", 9926, "hostname", "127.0.0.1")
	if err != nil {
		t.Fatalf("no se pudo validar el login, cod: %d error: %v", cod, err)
	}
	t.Log("User Login was successful")
	fmt.Print("User Login was successful")

}

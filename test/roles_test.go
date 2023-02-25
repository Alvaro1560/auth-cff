package test

import (
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/dbx"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth/roles"
	"testing"
)

func TestCreateRole(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	usr := models.User{
		ID:   "25c59a3b-7f3f-4dcd-b7fd-a4a7071d57f1",
		Name: "manager",
	}
	repoRole := roles.FactoryStorage(db, &usr, txID)
	srvRole := roles.NewRoleService(repoRole, nil, txID)

	_, cod, err := srvRole.CreateRole("767efeac-714d-483a-88be-48943344b8f9", "test role", "test role", 1)
	if err != nil {
		t.Fatalf("Role creation error, cod: %d error: %v", cod, err)
	}
	t.Log("Role creation was successful")
	fmt.Print("Role creation was successful")
}

func TestUpdateRole(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := roles.FactoryStorage(db, nil, txID)
	srvRole := roles.NewRoleService(repoRole, nil, txID)

	_, cod, err := srvRole.UpdateRole("767efeac-714d-483a-88be-48943344b8f9", "test role", "test role", 1)
	if err != nil {
		t.Fatalf("Role update error, cod: %d error: %v", cod, err)
	}
	t.Log("Role update was successful")
	fmt.Print("Role update was successful")
}

/*
func TestDeleteRole(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := roles.FactoryStorage(db, nil, txID)
	srvRole := roles.NewRoleService(repoRole, nil, txID)

	cod, err := srvRole.DeleteRole("767efeac-714d-483a-88be-48943344b8f9")
	if err != nil {
		t.Fatalf("Role delete error, cod: %d error: %v", cod, err)
	}
	t.Log("Role delete was successful")
	fmt.Print("Role delete was successful")
}*/

func TestGetRoleByID(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := roles.FactoryStorage(db, nil, txID)
	srvRole := roles.NewRoleService(repoRole, nil, txID)

	_, cod, err := srvRole.GetRoleByID("767efeac-714d-483a-88be-48943344b8f9")
	if err != nil {
		t.Fatalf("Role get role by id error, cod: %d error: %v", cod, err)
	}
	t.Log("Role get role by id was successful")
	fmt.Print("Role get role by id was successful")
}

func TestGetAllRole(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := roles.FactoryStorage(db, nil, txID)
	srvRole := roles.NewRoleService(repoRole, nil, txID)

	_, err := srvRole.GetAllRole()
	if err != nil {
		t.Fatalf("Role get all role error, error: %v", err)
	}
	t.Log("Role get all role was successful")
	fmt.Print("Role get all role was successful")
}

func TestGetRolesByUserID(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := roles.FactoryStorage(db, nil, txID)
	srvRole := roles.NewRoleService(repoRole, nil, txID)

	_, cod, err := srvRole.GetRolesByUserID("767efeac-714d-483a-88be-48943344b8f9")
	if err != nil {
		t.Fatalf("Role GetRolesByUserID error, cod: %d error: %v", cod, err)
	}
	t.Log("Role GetRolesByUserID was successful")
	fmt.Print("Role GetRolesByUserID was successful")
}

func TestGetRolesByProcessIDs(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := roles.FactoryStorage(db, nil, txID)
	srvRole := roles.NewRoleService(repoRole, nil, txID)
	requestProcess := []string{"767efeac-714d-483a-88be-48943344b8f9", "767efeac-714d-483a-88be-48943344b8f9"}
	_, err := srvRole.GetRolesByProcessIDs(requestProcess)
	if err != nil {
		t.Fatalf("Role GetRolesByProcessIDs error, error: %v", err)
	}
	t.Log("Role GetRolesByProcessIDs was successful")
	fmt.Print("Role GetRolesByProcessIDs was successful")
}

func TestGetRolesByIDs(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := roles.FactoryStorage(db, nil, txID)
	srvRole := roles.NewRoleService(repoRole, nil, txID)
	requestProcess := []string{"767efeac-714d-483a-88be-48943344b8f9", "767efeac-714d-483a-88be-48943344b8f9"}
	_, err := srvRole.GetRolesByIDs(requestProcess)
	if err != nil {
		t.Fatalf("Role GetRolesByIDs error, error: %v", err)
	}
	t.Log("Role GetRolesByIDs was successful")
	fmt.Print("Role GetRolesByIDs was successful")
}

func TestGetRolesByUserIDs(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := roles.FactoryStorage(db, nil, txID)
	srvRole := roles.NewRoleService(repoRole, nil, txID)
	idsUser := []string{"25c59a3b-7f3f-4dcd-b7fd-a4a7071d57f1"}
	_, cod, err := srvRole.GetRolesByUserIDs(idsUser)
	if err != nil {
		t.Fatalf("Role GetRolesByUserIDs error, cod: %d error: %v", cod, err)
	}
	t.Log("Role GetRolesByUserIDs was successful")
	fmt.Print("Role GetRolesByUserIDs was successful")
}

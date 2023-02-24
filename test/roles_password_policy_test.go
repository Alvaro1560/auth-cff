package test

import (
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/dbx"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/models"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth/roles_password_policy"
	"testing"
)

func TestCreateRolesPasswordPolicy(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	usr := models.User{
		ID:   "25c59a3b-7f3f-4dcd-b7fd-a4a7071d57f1",
		Name: "manager",
	}

	repoRolesPasswordPolicy := roles_password_policy.FactoryStorage(db, &usr, txID)
	srvRolesPasswordPolicy := roles_password_policy.NewRolesPasswordPolicyService(repoRolesPasswordPolicy, nil, txID)

	_, cod, err := srvRolesPasswordPolicy.CreateRolesPasswordPolicy("767efeac-714d-483a-88be-48943344b8f9",
		"38856274-5311-4173-97d1-fb7e2c8c9e16", 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, true, 1, 1)
	if err != nil {
		t.Fatalf("Error in CreateRolesPasswordPolicy rolesPasswordPolicy, cod: %d error: %v", cod, err)
	}
	t.Log("RolesPasswordPolicy CreateRolesPasswordPolicy was successful")
	fmt.Print("RolesPasswordPolicy CreateRolesPasswordPolicy was successful")

}

func TestUpdateRolesPasswordPolicy(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRolesPasswordPolicy := roles_password_policy.FactoryStorage(db, nil, txID)
	srvRolesPasswordPolicy := roles_password_policy.NewRolesPasswordPolicyService(repoRolesPasswordPolicy, nil, txID)

	_, cod, err := srvRolesPasswordPolicy.UpdateRolesPasswordPolicy("767efeac-714d-483a-88be-48943344b8f9",
		"38856274-5311-4173-97d1-fb7e2c8c9e16", 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, true, 1, 1)
	if err != nil {
		t.Fatalf("Error in UpdateRolesPasswordPolicy rolesPasswordPolicy, cod: %d error: %v", cod, err)
	}
	t.Log("RolesPasswordPolicy UpdateRolesPasswordPolicy was successful")
	fmt.Print("RolesPasswordPolicy UpdateRolesPasswordPolicy was successful")
}

func TestDeleteRolesPasswordPolicy(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRolesPasswordPolicy := roles_password_policy.FactoryStorage(db, nil, txID)
	srvRolesPasswordPolicy := roles_password_policy.NewRolesPasswordPolicyService(repoRolesPasswordPolicy, nil, txID)

	_, err := srvRolesPasswordPolicy.DeleteRolesPasswordPolicy("767efeac-714d-483a-88be-48943344b8f9")
	if err != nil {
		t.Fatalf("Error in DeleteRolesPasswordPolicy rolesPasswordPolicy, error: %v", err)
	}
	t.Log("RolesPasswordPolicy DeleteRolesPasswordPolicy was successful")
	fmt.Print("RolesPasswordPolicy DeleteRolesPasswordPolicy was successful")
}

func TestGetRolesPasswordPolicyByID(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRolesPasswordPolicy := roles_password_policy.FactoryStorage(db, nil, txID)
	srvRolesPasswordPolicy := roles_password_policy.NewRolesPasswordPolicyService(repoRolesPasswordPolicy, nil, txID)

	_, cod, err := srvRolesPasswordPolicy.GetRolesPasswordPolicyByID("767efeac-714d-483a-88be-48943344b8f9")
	if err != nil {
		t.Fatalf("Error in GetRolesPasswordPolicyByID rolesPasswordPolicy, cod: %d error: %v", cod, err)
	}
	t.Log("RolesPasswordPolicy GetRolesPasswordPolicyByID was successful")
	fmt.Print("RolesPasswordPolicy GetRolesPasswordPolicyByID was successful")
}

func TestGetAllRolesPasswordPolicy(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRolesPasswordPolicy := roles_password_policy.FactoryStorage(db, nil, txID)
	srvRolesPasswordPolicy := roles_password_policy.NewRolesPasswordPolicyService(repoRolesPasswordPolicy, nil, txID)

	_, err := srvRolesPasswordPolicy.GetAllRolesPasswordPolicy()
	if err != nil {
		t.Fatalf("Error in GetAllRolesPasswordPolicy rolesPasswordPolicy, error: %v", err)
	}
	t.Log("RolesPasswordPolicy GetAllRolesPasswordPolicy was successful")
	fmt.Print("RolesPasswordPolicy GetAllRolesPasswordPolicy was successful")
}

func TestGetAllRolesPasswordPolicyByRolesIDs(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRolesPasswordPolicy := roles_password_policy.FactoryStorage(db, nil, txID)
	srvRolesPasswordPolicy := roles_password_policy.NewRolesPasswordPolicyService(repoRolesPasswordPolicy, nil, txID)
	requestProcess := []string{"767efeac-714d-483a-88be-48943344b8f9", "767efeac-714d-483a-88be-48943344b8f9"}

	_, err := srvRolesPasswordPolicy.GetAllRolesPasswordPolicyByRolesIDs(requestProcess)
	if err != nil {
		t.Fatalf("Error in GetAllRolesPasswordPolicyByRolesIDs rolesPasswordPolicy, error: %v", err)
	}
	t.Log("RolesPasswordPolicy GetAllRolesPasswordPolicyByRolesIDs was successful")
	fmt.Print("RolesPasswordPolicy GetAllRolesPasswordPolicyByRolesIDs was successful")
}

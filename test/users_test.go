package test

import (
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/dbx"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth/users"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := users.FactoryStorage(db, nil, txID)
	srvRole := users.NewUserService(repoRole, nil, txID)

	_, cod, err := srvRole.CreateUser("767efeac-714d-483a-88be-48943344b8f9", "test123", "test", "test", "tes123", "test@e-capture.co", "12345678", "CC")
	if err != nil {
		t.Fatalf("Users CreateUser error, cod: %d error: %v", cod, err)
	}
	t.Log("Users CreateUser was successful")
	fmt.Print("Users CreateUser was successful")
}

/*
func TestUpdateUser(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := users.FactoryStorage(db, nil, txID)
	srvRole := users.NewUserService(repoRole, nil, txID)

	_, cod, err := srvRole.UpdateUser("767efeac-714d-483a-88be-48943344b8f9", "test123", "test", "test", "tes123", "test@e-capture.co", "12345678", "CC")
	if err != nil {
		t.Fatalf("Users UpdateUser error, cod: %d error: %v", cod, err)
	}
	t.Log("Users UpdateUser was successful")
	fmt.Print("Users UpdateUser was successful")
}*/

func TestGetUserByID(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := users.FactoryStorage(db, nil, txID)
	srvRole := users.NewUserService(repoRole, nil, txID)

	_, cod, err := srvRole.GetUserByID("767efeac-714d-483a-88be-48943344b8f9")
	if err != nil {
		t.Fatalf("Users GetUserByID error, cod: %d error: %v", cod, err)
	}
	t.Log("Users GetUserByID was successful")
	fmt.Print("Users GetUserByID was successful")
}

func TestGetUserByUsername(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := users.FactoryStorage(db, nil, txID)
	srvRole := users.NewUserService(repoRole, nil, txID)

	_, cod, err := srvRole.GetUserByUsername("767efeac-714d-483a-88be-48943344b8f9")
	if err != nil {
		t.Fatalf("Users GetUserByUsername error, cod: %d error: %v", cod, err)
	}
	t.Log("Users GetUserByUsername was successful")
	fmt.Print("Users GetUserByUsername was successful")
}

func TestGetAllUser(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := users.FactoryStorage(db, nil, txID)
	srvRole := users.NewUserService(repoRole, nil, txID)

	_, err := srvRole.GetAllUser()
	if err != nil {
		t.Fatalf("Users GetAllUser error, error: %v", err)
	}
	t.Log("Users GetAllUser was successful")
	fmt.Print("Users GetAllUser was successful")
}

func TestGetUsersByIDs(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := users.FactoryStorage(db, nil, txID)
	srvRole := users.NewUserService(repoRole, nil, txID)
	requestProcess := []string{"767efeac-714d-483a-88be-48943344b8f9", "767efeac-714d-483a-88be-48943344b8f9"}

	_, err := srvRole.GetUsersByIDs(requestProcess)
	if err != nil {
		t.Fatalf("Users GetUsersByIDs error, error: %v", err)
	}
	t.Log("Users GetUsersByIDs was successful")
	fmt.Print("Users GetUsersByIDs was successful")
}

func TestBlockUser(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := users.FactoryStorage(db, nil, txID)
	srvRole := users.NewUserService(repoRole, nil, txID)
	err := srvRole.BlockUser("767efeac-714d-483a-88be-48943344b8f9")
	if err != nil {
		t.Fatalf("Users BlockUser error, error: %v", err)
	}
	t.Log("Users BlockUser was successful")
	fmt.Print("Users BlockUser was successful")
}

func TestUnblockUser(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := users.FactoryStorage(db, nil, txID)
	srvRole := users.NewUserService(repoRole, nil, txID)
	err := srvRole.UnblockUser("767efeac-714d-483a-88be-48943344b8f9")
	if err != nil {
		t.Fatalf("Users UnblockUser error, error: %v", err)
	}
	t.Log("Users UnblockUser was successful")
	fmt.Print("Users UnblockUser was successful")
}

/*func TestLogoutUser(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := users.FactoryStorage(db, nil, txID)
	srvRole := users.NewUserService(repoRole, nil, txID)

	cod, err := srvRole.DeleteUserPasswordHistory("767efeac-714d-483a-88be-48943344b8f9")
	if err != nil {
		t.Fatalf("Users LogoutUser error, cod: %d, error: %v", cod, err)
	}

	cod, err = srvRole.LogoutUser("767efeac-714d-483a-88be-48943344b8f9")
	if err != nil {
		t.Fatalf("Users LogoutUser error, cod: %d, error: %v", cod, err)
	}
	t.Log("Users LogoutUser was successful")
	fmt.Print("Users LogoutUser was successful")
}*/

func TestChangePassword(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := users.FactoryStorage(db, nil, txID)
	srvRole := users.NewUserService(repoRole, nil, txID)
	cod, err := srvRole.ChangePassword("767efeac-714d-483a-88be-48943344b8f9", "test123", "test123")
	if err != nil {
		t.Fatalf("Users ChangePassword error, cod: %d, error: %v", cod, err)
	}
	t.Log("Users ChangePassword was successful")
	fmt.Print("Users ChangePassword was successful")
}

/*
	func TestUpdatePasswordByUser(t *testing.T) {
		db := dbx.GetConnection()
		txID := uuid.New().String()
		usr := models.User{
			ID:   "25c59a3b-7f3f-4dcd-b7fd-a4a7071d57f1",
			Name: "manager",
		}
		repoRole := users.FactoryStorage(db, &usr, txID)
		srvRole := users.NewUserService(repoRole, &usr, txID)
		cod, err := srvRole.UpdatePasswordByUser("767efeac-714d-483a-88be-48943344b8f9", "test123", "test123", "test123")
		if err != nil {
			t.Fatalf("Users UpdatePasswordByUser error, cod: %d, error: %v", cod, err)
		}
		t.Log("Users UpdatePasswordByUser was successful")
		fmt.Print("Users UpdatePasswordByUser was successful")
	}
*/
/*
func TestUpdatePasswordByAdministrator(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	usr := models.User{
		ID:   "25c59a3b-7f3f-4dcd-b7fd-a4a7071d57f1",
		Name: "manager",
	}
	repoRole := users.FactoryStorage(db, &usr, txID)
	srvRole := users.NewUserService(repoRole, &usr, txID)

	cod, err := srvRole.DeleteUserPasswordHistory("767efeac-714d-483a-88be-48943344b8f9")
	if err != nil {
		t.Fatalf("Users LogoutUser error, cod: %d, error: %v", cod, err)
	}

	cod, err = srvRole.UpdatePasswordByAdministrator("767efeac-714d-483a-88be-48943344b8f9", "test123", "test123")
	if err != nil {
		t.Fatalf("Users UpdatePasswordByAdministrator error, cod: %d, error: %v", cod, err)
	}
	t.Log("Users UpdatePasswordByAdministrator was successful")
	fmt.Print("Users UpdatePasswordByAdministrator was successful")
}*/
/*
func TestGetUserByUsernameAndIdentificationNumber(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := users.FactoryStorage(db, nil, txID)
	srvRole := users.NewUserService(repoRole, nil, txID)
	_, cod, err := srvRole.GetUserByUsernameAndIdentificationNumber("manager", "80188284")
	if err != nil {
		t.Fatalf("Users GetUserByUsernameAndIdentificationNumber error, cod: %d, error: %v", cod, err)
	}
	t.Log("Users GetUserByUsernameAndIdentificationNumber was successful")
	fmt.Print("Users GetUserByUsernameAndIdentificationNumber was successful")
}*/

/*func TestValidatePasswordPolicy(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := users.FactoryStorage(db, nil, txID)
	srvRole := users.NewUserService(repoRole, nil, txID)
	_, cod, err := srvRole.ValidatePasswordPolicy("test123", 1, 1, 1, 1, 1, 1, 1, true)
	if err != nil {
		t.Fatalf("Users ValidatePasswordPolicy error, cod: %d, error: %v", cod, err)
	}
	t.Log("Users ValidatePasswordPolicy was successful")
	fmt.Print("Users ValidatePasswordPolicy was successful")
}*/

func TestDeleteUser(t *testing.T) {
	db := dbx.GetConnection()
	txID := uuid.New().String()
	repoRole := users.FactoryStorage(db, nil, txID)
	srvRole := users.NewUserService(repoRole, nil, txID)

	cod, err := srvRole.DeleteUserPasswordHistory("767efeac-714d-483a-88be-48943344b8f9")
	if err != nil {
		t.Fatalf("Users LogoutUser error, cod: %d, error: %v", cod, err)
	}

	_, err = srvRole.DeleteUser("767efeac-714d-483a-88be-48943344b8f9")
	if err != nil {
		t.Fatalf("Users DeleteUser error, error: %v", err)
	}
	t.Log("Users DeleteUser was successful")
	fmt.Print("Users DeleteUser was successful")
}

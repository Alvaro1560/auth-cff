package transact

import (
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/dbx"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/transactions/loggedusers"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
)

func RegisterConfig(action string, description string, user string) {

}

func RegisterLogUsr(Event string, HostName string, IpRequest string, IpRemote string, UserId string) {
	conn := dbx.GetConnection()

	repoLoggedUsers := loggedusers.FactoryStorage(conn, nil, "")
	srvLoggedUsers := loggedusers.NewTxLoggedUserService(repoLoggedUsers, nil, "")

	_, _, err := srvLoggedUsers.CreateTxLoggedUser(Event, HostName, IpRequest, IpRemote, UserId)
	if err != nil {
		logger.Error.Println("", " - couldn't create loggedUsers :", err)
	}
}

func RegisterTrace(typeMessage string, fileName string, codeLine int, message string, transaction string) {

}

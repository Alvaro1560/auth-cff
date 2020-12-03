package msgs

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/config/messages"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/dbx"
)

type Model struct {
	db *sqlx.DB
}

func (m *Model) GetByCode(code int) (int, string, string) {

	db := dbx.GetConnection()
	repoMsg := messages.FactoryStorage(db, nil, "")
	srvMsg := messages.NewMessageService(repoMsg, nil, "")
	msg, _, err := srvMsg.GetMessageByID(code)
	if err != nil {
		return 70, "", ""
	}

	return msg.ID, msg.TypeMessage, msg.Spa

}

package main

import (
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/api"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/env"
)

func main() {
	c := env.NewConfiguration()
	api.Start(c.App.Port, c.App.ServiceName, c.App.LoggerHttp, c.App.AllowedDomains)
}

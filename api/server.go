package api

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/fatih/color"
)

const (
	version = "0.1.0"
	website = "http://www.ecapture.com.co"
	banner  = `
███████  ██████  █████  ████████  ██████ ██   ██ 
██      ██      ██   ██    ██    ██      ██   ██ 
█████   ██      ███████    ██    ██      ███████ 
██      ██      ██   ██    ██    ██      ██   ██ 
███████  ██████ ██   ██    ██     ██████ ██   ██ 
                                                 `

	description = `Ecatch BPM - %s - Port: %s
by Ecapture
Version: %s
%s`
)

type server struct {
	listening string
	app       string
	fb        *fiber.App
}

func newServer(listening int, app string, fb *fiber.App) *server {
	return &server{fmt.Sprintf(":%d", listening), app, fb}
}

func (srv *server) Start() {
	color.Blue(banner)
	color.Cyan(fmt.Sprintf(description, srv.app, srv.listening, version, website))
	log.Fatal(srv.fb.Listen(srv.listening))
}

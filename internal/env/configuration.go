package env

import (
	"encoding/json"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/ciphers"
	"io/ioutil"
	"log"
	"sync"
)

var (
	once   sync.Once
	config = &configuration{}
)

type configuration struct {
	App      App      `json:"app"`
	DB       DB       `json:"db"`
	Smtp     Smtp     `json:"smtp"`
	Template Template `json:"template"`
}

type App struct {
	ServiceName       string `json:"service_name"`
	Port              int    `json:"port"`
	AllowedDomains    string `json:"allowed_domains"`
	PathLog           string `json:"path_log"`
	LogReviewInterval int    `json:"log_review_interval"`
	RegisterLog       bool   `json:"register_log"`
	RSAPrivateKey     string `json:"rsa_private_key"`
	RSAPublicKey      string `json:"rsa_public_key"`
	UrlPortal         string `json:"url_portal"`
	LoggerHttp        bool   `json:"logger_http"`
	IsCipher          bool   `json:"is_cipher"`
	KeywordAutologin  string `json:"keyword_autologin"`
	Autologin         bool   `json:"autologin"`
	User              string `json:"user"`
	Password          string `json:"password"`
}

type DB struct {
	Engine   string `json:"engine"`
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Instance string `json:"instance"`
	IsSecure bool   `json:"is_secure"`
	SSLMode  string `json:"ssl_mode"`
}
type Smtp struct {
	Port     int    `json:"port"`
	Host     string `json:"host"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Template struct {
	Recovery         string `json:"recovery"`
	EmailCode        string `json:"email_code"`
	EmailCodeSubject string `json:"email_code_subject"`
	EmailSender      string `json:"email_sender"`
}

func NewConfiguration() *configuration {
	fromFile()
	return config
}

// LoadConfiguration lee el archivo configuration.json
// y lo carga en un objeto de la estructura Configuration
func fromFile() {
	once.Do(func() {
		b, err := ioutil.ReadFile("config.json")
		if err != nil {
			log.Fatalf("no se pudo leer el archivo de configuración: %s", err.Error())
		}

		err = json.Unmarshal(b, config)
		if err != nil {
			log.Fatalf("no se pudo parsear el archivo de configuración: %s", err.Error())
		}

		if config.DB.Engine == "" {
			log.Fatal("no se ha cargado la información de configuración")
		}

		if config.App.IsCipher {
			if config.DB.Password = ciphers.Decrypt(config.DB.Password); config.DB.Password == "" {
				log.Fatal("no se pudo obtener config.DB.Password Decrypt")
			}
		}

	})
}

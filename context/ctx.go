package context

import (
	"log"

	"github.com/caarlos0/env/v10"
)

type Env struct {
	Creds           string   `env:"PROXY_CREDENTIALS" envDefault:""`
	ProxyUser       string   `env:"PROXY_USER" envDefault:""`
	ProxyPassword   string   `env:"PROXY_PASS" envDefault:""`
	ProxyPort       string   `env:"PROXY_PORT" envDefault:"1080"`
	AllowedDestFqdn string   `env:"PROXY_ALLOWED_DEST_FQDN" envDefault:""`
	AllowedIPs      []string `env:"PROXY_ALLOWED_IPS" envDefault:""`

	// web ui
	WebUIPort string `env:"WEBUI_PORT" envDefault:"8080"`
	WebUIUser string `env:"WEBUI_USER" envDefault:""`
	WebUIPass string `env:"WEBUI_PASS" envDefault:""`

	// security
	JwtSecret string `env:"JWT_SECRET" envDefault:"secret"`

	// database
	DbPath string `env:"DB_PATH" envDefault:"main.db"`
}

var ctx = Env{}

func GetEnv() Env {
	if err := env.Parse(&ctx); err != nil {
		log.Fatal(err)
	}

	return ctx
}

package main

import (
	"log"
	"log/slog"

	"github.com/alexng353/ihostproxy/credentials"
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
}

var creds = credentials.Get()

var ctx = Env{}

var B = "hello world"

func main() {
	if err := env.Parse(&ctx); err != nil {
		log.Fatal(err)
	}

	if ctx.WebUIPass != "" && ctx.WebUIUser != "" {
		err := creds.AddAdmin(ctx.WebUIUser, ctx.WebUIPass)

		if err != nil {
			slog.Error("Failed to add admin", "error", err)
		}
	}

	go startProxy(ctx)
	go startWebUI(ctx)

	for {
	}
}

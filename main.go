package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexng353/ihostproxy/context"
	"github.com/alexng353/ihostproxy/credentials"
	"github.com/alexng353/ihostproxy/webui"
	"github.com/caarlos0/env/v10"
)

var creds = credentials.Get()

var ctx = context.Env{}

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

	disableProxy := os.Getenv("DISABLE_PROXY") == "1"
	if !disableProxy {
		go startProxy(ctx)
	}

	disableWebui := os.Getenv("DISABLE_WEBUI") == "1"
	if !disableWebui {
		go webui.StartWebUI(ctx)
	}

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs

	slog.String("signal", sig.String())
}

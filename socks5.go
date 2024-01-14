package main

import (
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/alexng353/ihostproxy/context"
	"github.com/alexng353/ihostproxy/credentials"
	"github.com/alexng353/ihostproxy/helpers"
	"github.com/things-go/go-socks5"
)

func startProxy(ctx context.Env) {
	var c = credentials.Get()
	if ctx.Creds != "" {
		err := c.Load(ctx.Creds)
		if err != nil {
			slog.Error("Failed to load credentials", "error", err)
		}
	}

	if ctx.ProxyUser != "" && ctx.ProxyPassword != "" {
		err := c.AddEntry(ctx.ProxyUser, ctx.ProxyPassword)
		if err != nil {
			slog.Error("Failed to add entry", "error", err)
		}
	}

	cator := socks5.UserPassAuthenticator{Credentials: c}
	server := socks5.NewServer(
		socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "socks5: ", log.LstdFlags))),
		socks5.WithAuthMethods([]socks5.Authenticator{cator}),
	)

	port := strconv.FormatInt(int64(helpers.ValidatePort(ctx.ProxyPort)), 10)

	slog.Info("Starting Socks5 Proxy", "port", port)
	if err := server.ListenAndServe("tcp", ":"+port); err != nil {
		panic(err)
	}
}

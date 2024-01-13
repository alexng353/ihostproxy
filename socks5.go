package main

import (
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/things-go/go-socks5"
)

func startProxy(ctx Env) {
	var credentials = NewSQLiteCredentialStore()
	if ctx.Creds != "" {
		err := credentials.Load(ctx.Creds)
		if err != nil {
			slog.Error("Failed to load credentials", "error", err)
		}
	}

	if ctx.ProxyUser != "" && ctx.ProxyPassword != "" {
		err := credentials.AddEntry(ctx.ProxyUser, ctx.ProxyPassword)
		if err != nil {
			slog.Error("Failed to add entry", "error", err)
		}
	}

	cator := socks5.UserPassAuthenticator{Credentials: credentials}
	server := socks5.NewServer(
		socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "socks5: ", log.LstdFlags))),
		socks5.WithAuthMethods([]socks5.Authenticator{cator}),
	)

	port := strconv.FormatInt(int64(validatePort(ctx.ProxyPort)), 10)

	slog.Info("Starting Socks5 Proxy", "port", port)
	if err := server.ListenAndServe("tcp", ":"+port); err != nil {
		panic(err)
	}
}

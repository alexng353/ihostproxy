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
	// "github.com/armon/go-socks5"
)

var c = credentials.Get()

func startProxy(ctx context.Env) {
	if ctx.Creds != "" {
		err := c.Load(ctx.Creds)
		if err != nil {
			slog.Error("Failed to load credentials", "error", err)
		}
	}

	// if ctx.ProxyUser != "" && ctx.ProxyPassword != "" {
	// 	err := c.AddEntry(ctx.ProxyUser, ctx.ProxyPassword)
	// 	if err != nil {
	// 		slog.Error("Failed to add entry", "error", err)
	// 	}
	// }

	// socks5conf := &socks5.Config{
	// 	Logger: log.New(os.Stdout, "socks5: ", log.LstdFlags),
	// }

	// cator := socks5.UserPassAuthenticator{Credentials: c}

	// cator := socks5.UserPassAuthenticator{
	// 	Credentials: socks5.StaticCredentials{
	// 		ctx.ProxyUser: ctx.ProxyPassword,
	// 	},
	// }

	// socks5conf.AuthMethods = []socks5.Authenticator{cator}

	// server, err := socks5.New(socks5conf)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	cator := socks5.UserPassAuthenticator{Credentials: c}
	// cator := socks5.UserPassAuthenticator{
	// 	Credentials: socks5.StaticCredentials{
	// 		ctx.ProxyUser: ctx.ProxyPassword,
	// 	},
	// }
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

package main

import (
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/postfinance/discovery/internal/cmd/server"
	"github.com/postfinance/flash"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zbindenren/king"
)

const appName = "discovery"

//nolint:gochecknoglobals //this vars are set on build by goreleaser
var (
	version = "0.0.0"
	commit  = "12345678"
	date    = "2020-09-23T07:03:55+02:00"
)

func main() {
	registry := prometheus.NewRegistry()
	cli := server.CLI{}
	l := flash.New(flash.WithColor(), flash.WithStacktrace(), flash.WithPrometheus(appName, registry))

	b, err := king.NewBuildInfo(version,
		king.WithDateString(date),
		king.WithRevision(commit),
	)
	if err != nil {
		l.Fatal(err)
	}

	app := kong.Parse(&cli, king.DefaultOptions(
		king.Config{
			Name:        appName,
			Description: "GRPC discovery service.",
			BuildInfo:   b,
		},
	)...)

	l.SetDebug(cli.Debug)

	if cli.Profiler.Enabled {
		cli.Profiler.New(syscall.SIGUSR2).Start()
	}

	if err := app.Run(&cli.Globals, l.Get(), registry); err != nil {
		l.Fatal(err)
	}
}

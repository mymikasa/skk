package main

import (
	"github.com/mymikasa/skk/pkg/events"
	"github.com/mymikasa/skk/pkg/grpcx"
)

type App struct {
	consumers []events.Consumer
	server    *grpcx.Server
}

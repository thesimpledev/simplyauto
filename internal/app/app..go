package app

import "log"

type App struct {
	Log *log.Logger
}

func New() *App {
	// Load settings from registry
	a := &App{}

	return a
}

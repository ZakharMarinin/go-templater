package app

import (
	"go-templater/internal/config"
	"log/slog"
)

type App struct {
	cfg *config.Config
	log *slog.Logger
	di  *Container
}

func New(cfg *config.Config, log *slog.Logger) *App {
	return &App{
		cfg: cfg,
		log: log,
		di: NewContainer(cfg, log),
	}
}

func (a *App) MustRun() {
	err := a.Run()
	if err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	a.di.Cobra().Execute()

	return nil
}

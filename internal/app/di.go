package app

import (
	"github.com/ZakharMarinin/go-templater/internal/cli"
	"github.com/ZakharMarinin/go-templater/internal/config"
	"github.com/ZakharMarinin/go-templater/internal/ui"
	"github.com/ZakharMarinin/go-templater/internal/usecase"
	"log/slog"
)

type Container struct {
	cfg     *config.Config
	log     *slog.Logger
	ui      *ui.UI
	usecase *usecase.UseCase
	cobra   *cli.Cobra
}

func NewContainer(cfg *config.Config, log *slog.Logger) *Container {
	return &Container{cfg: cfg, log: log}
}

func (di *Container) UI() *ui.UI {
	if di.ui == nil {
		di.ui = ui.New(di.cfg)
	}

	return di.ui
}

func (di *Container) UseCase() *usecase.UseCase {
	if di.usecase == nil {
		di.usecase = usecase.New(di.cfg, di.log, di.UI())
	}

	return di.usecase
}

func (di *Container) Cobra() *cli.Cobra {
	if di.cobra == nil {
		di.cobra = cli.New(di.cfg, di.UseCase())
	}

	return di.cobra
}

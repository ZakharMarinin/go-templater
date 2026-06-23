package ui

import "go-templater/internal/config"

const (
	enter  = "enter"
	tab    = "tab"
	shift  = "shift"
	up     = "up"
	down   = "down"
	esc    = "esc"
	cancel = "ctrl+c"
)

type UI struct {
	cfg *config.Config
}

func New(cfg *config.Config) *UI {
	return &UI{cfg: cfg}
}

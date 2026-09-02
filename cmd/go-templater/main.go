package main

import (
	"fmt"
	"github.com/ZakharMarinin/go-templater/internal/app"
	"github.com/ZakharMarinin/go-templater/internal/config"
	"github.com/ZakharMarinin/go-templater/internal/libs/logger"
)

func main() {
	cfg := config.MustLoad()
	log, err := logger.Logger(cfg.Env, cfg.Routes.LogsDir)
	if err != nil {
		fmt.Printf("error: %s", err)
		
		return
	}

	app := app.New(cfg, log)

	app.MustRun()
}

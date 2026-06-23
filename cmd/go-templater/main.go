package main

import (
	"fmt"
	"go-templater/internal/app"
	"go-templater/internal/config"
	"go-templater/internal/libs/logger"
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

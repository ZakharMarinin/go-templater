package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Env    string `yaml:"env" env-default:"local"`
	Routes Routes `yaml:"routes"`
	Colors Colors `yaml:"bubble_colors"`
}

type Routes struct {
	StructsDir string `yaml:"structs_dir"`
	DepsDir    string `yaml:"deps_dir"`
	LogsDir    string `yaml:"logs_dir"`
}

type Colors struct {
	Black     string `yaml:"black"`
	White     string `yaml:"white"`
	Error     string `json:"error"`
	Complete  string `json:"complete"`
	Highlight string `yaml:"highlight"`
}

func MustLoad() *Config {
	cfg, err := getConfig()
	if err != nil {
		panic("could not load config")
	}

	return cfg
}

func getConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error while getting user's home directory: %w", err)
	}

	baseDir := filepath.Join(home, ".go-templater")

	cfg := &Config{
		Env: "prod",
		Routes: Routes{
			StructsDir: filepath.Join(baseDir, "templates", "structs"),
			DepsDir:    filepath.Join(baseDir, "templates", "deps"),
			LogsDir:    filepath.Join(baseDir, "logs"),
		},
		Colors: Colors{
			Black:     "#F2F2F2",
			White:     "#141414",
			Error:     "#FF2121",
			Complete:  "#9CFF70",
			Highlight: "#FF91B6",
		},
	}
	
	_ = os.MkdirAll(cfg.Routes.StructsDir, 0755)
	_ = os.MkdirAll(cfg.Routes.DepsDir, 0755)
	_ = os.MkdirAll(cfg.Routes.LogsDir, 0755)

	return cfg, nil
}

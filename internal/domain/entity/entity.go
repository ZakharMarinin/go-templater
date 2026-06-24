package entity

import "time"

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

type Template struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Dependencies []*Dependency `json:"dependencies"`
	Nodes        []*Node       `json:"dirs"`
}

type TemplateInfo struct {
	Name      string
	Type      string
	Size      int64
	CreatedAt time.Time
}

type Node struct {
	Name     string  `json:"name"`
	Path     string  `json:"path"`
	IsDir    bool    `json:"is_dir"`
	Children []*Node `json:"children"`
}

type Dependency struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Version string `json:"version"`
}

type Variables struct {
	Name        string
	Description string
}

type FieldConfig struct {
	Key         string
	Placeholder string
	Width       int
	CharLimit   int
	Required    bool
}

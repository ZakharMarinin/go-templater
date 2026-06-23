package usecase

import (
	"fmt"
	"go-templater/internal/domain/entity"
	"go-templater/pkg/response"
	"os"
	"os/exec"
	"path/filepath"
)

func (u *UseCase) InsertDepsTemplate(homeDir string) error {
	const op = "usecase.InsertDepsTemplate"

	log := u.log.With("operation", op)

	templates, err := getTemplates(u.cfg.Routes.DepsDir)
	if err != nil {
		log.Error("could not get templates", "error", err)

		return err
	}

	template, err := u.ui.Select(templates)
	if err != nil {
		log.Error("error with template selection", "error", err)

		return err
	}

	task := func() error {
		return insertDeps(template.Dependencies, homeDir)
	}

	err = u.ui.NewSpinner("Template dependencies installation", task)
	if err != nil {
		log.Warn("could not install template dependencies", "error", err)

		return err
	}

	return nil
}

func insertDeps(deps []*entity.Dependency, path string) error {
	goFile := filepath.Join(path, "go.mod")
	_, err := os.Stat(goFile)
	if os.IsNotExist(err) {
		return response.ErrNotExist
	}

	args := make([]string, 0, len(deps)+1)
	args = append(args, "get")

	for _, dep := range deps {
		version := dep.Version
		if version == "" {
			version = "latest"
		}

		target := fmt.Sprintf("%s@%s", dep.URL, version)
		args = append(args, target)
	}

	cmd := exec.Command("go", args...)
	cmd.Dir = path

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("could not get dependencies: %w. out: %s", err, string(out))
	}

	return nil
}

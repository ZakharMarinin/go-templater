package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"go-templater/internal/config"
	"go-templater/internal/domain/entity"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type UI interface {
	Select(templates []*entity.Template) (*entity.Template, error)
	Input() (*entity.Variables, error)
	ConfirmOverwrite(fileName string) (bool, error)
	ShowStatus(msg string, duration time.Duration) error 
	NewSpinner(title string, task func() error) error
}

type UseCase struct {
	cfg *config.Config
	log *slog.Logger
	ui  UI
}

func New(cfg *config.Config, log *slog.Logger, ui UI) *UseCase {
	return &UseCase{cfg: cfg, log: log, ui: ui}
}

func InsertTemplate(ctx context.Context) {

}

func DeleteTemplate(ctx context.Context) {

}

func isUnique(path, name string) bool {
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if d.Name() == name+".json" {
			return fmt.Errorf("not unique name")
		}

		return nil
	})

	return err == nil
}

func getTemplates(path string) ([]*entity.Template, error) {
	var templates []*entity.Template

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if err != nil {
			return err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var template entity.Template

		err = json.Unmarshal(data, &template)
		if err != nil {
			return err
		}

		templates = append(templates, &template)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return templates, nil
}
package usecase

import (
	"fmt"
	"go-templater/internal/domain/entity"
	"os"
	"os/exec"
)

func (u *UseCase) InsertDirTemplate(homeDir string) error {
	const op = "usecase.InsertDirTemplate"

	log := u.log.With("operation", op)

	templates, err := getTemplates(u.cfg.Routes.StructsDir)
	if err != nil {
		log.Error("could not get templates", "error", err)

		return err
	}

	template, err := u.ui.Select(templates)
	if err != nil {
		log.Error("error with template selection", "error", err)

		return err
	}

	nameKey := "project"

	configs := []*entity.FieldConfig{
		{Key: nameKey, Placeholder: "Name your project", Required: true},
	}

	variables, err := u.ui.DynamicInput("Create project", configs)
	if err != nil {
		log.Error("error with name input", "error", err)

		return err
	}
	
	err = insertDirs(template.Nodes, homeDir)
	if err != nil {
		log.Error("could not insert template", "error", err)

		return err
	}

	err = creaProject(homeDir, variables[nameKey])
	if err != nil {
		log.Error("could not create project", "error", err)

		return err
	}

	return nil
}

func insertDirs(nodes []*entity.Node, path string) error {
	for _, i := range nodes {
		name := path + i.Path

		if i.IsDir {
			err := os.Mkdir(name, 0777)
			if err != nil {
				return err
			}

			err = insertDirs(i.Children, path)
			if err != nil {
				return err
			}
		} else {
			err := os.WriteFile(name, []byte{}, 0777)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func creaProject(path, projectName string) error {
	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Dir = path

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error: %w. out: %s", err, string(out))
	}

	return nil
}
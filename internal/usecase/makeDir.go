package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-templater/internal/domain/entity"
	"go-templater/pkg/response"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func (u *UseCase) MakeStructTemplate(homeDir string) error {
	const op = "usecase.MakeStructTemplate"
	log := u.log.With("operation", op)

	nameKey := "name"
	descKey := "desc"
	
	configs := []*entity.FieldConfig {
	    {Key: nameKey, Placeholder: "Name (required)", Required: true, Width: 32},
	    {Key: descKey, Placeholder: "Description (optional)", Required: false, Width: 64, CharLimit: 64},
	}
	
	variables, err := u.ui.DynamicInput("Create template", configs)
	if err != nil {
		if errors.Is(err, response.ErrCanceled) {
			return nil
		}
		
		log.Error("could not get template variables: %w", "error", err)

		return err
	}

	if !isUnique(u.cfg.Routes.StructsDir, variables[nameKey]) {
		err := u.confirmStatus(variables[nameKey])
		if err != nil {
			return err
		}
	}

	root, err := getUserTree(homeDir)
	if err != nil {
		return fmt.Errorf("could not get user's directory tree. error: %w", err)
	}

	template := &entity.Template{
		Name:        variables[nameKey],
		Description: variables[descKey],
		Nodes:       root,
	}
	
	err = createDirFile(u.cfg.Routes.StructsDir, homeDir, template)
	if err != nil {
		log.Error("could not created template", "error", err)

		return err
	}

	return nil
}

func createDirFile(path, homeDir string, template *entity.Template) error {
	data, err := json.Marshal(template)
	if err != nil {
		return fmt.Errorf("could not marshal data. error: %w", err)
	}

	name := template.Name + ".json"

	fullPath := filepath.Join(path, name)

	err = os.WriteFile(fullPath, data, 0777)
	if err != nil {
		return fmt.Errorf("could not create file. error: %w", err)
	}

	return nil
}

func getUserTree(homeDir string) ([]*entity.Node, error) {
	nodesMap := make(map[string]*entity.Node)
	var rootNode *entity.Node

	fullPath, err := filepath.Abs(homeDir)
	if err != nil {
		return nil, err
	}
	
	err = filepath.WalkDir(fullPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		err = nameFilter(d.Name(), d.IsDir())
		if err != nil {
			if errors.Is(err, filepath.SkipDir) {
				return err
			}

			return nil
		}

		folderPath, _ := strings.CutPrefix(path, fullPath)

		node := &entity.Node{
			Name:  d.Name(),
			Path:  folderPath,
			IsDir: d.IsDir(),
		}

		nodesMap[path] = node

		parentPath := filepath.Dir(path)
		parent, ok := nodesMap[parentPath]
		if ok {
			parent.Children = append(parent.Children, node)
		}

		if path == fullPath {
			rootNode = node

			return nil
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	
	return rootNode.Children, nil
}

func nameFilter(name string, isDir bool) error {
	if isDir {
		switch name {
		case ".git":
			return filepath.SkipDir
		case ".idea":
			return filepath.SkipDir
		case ".vscode":
			return filepath.SkipDir
		}
	} else {
		switch name {
		case "go.mod":
			return fmt.Errorf("skip file")
		case "go.sum":
			return fmt.Errorf("skip file")
		}
	}

	return nil
}

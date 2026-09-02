package usecase

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ZakharMarinin/go-templater/internal/domain/entity"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func (u *UseCase) MakeDepsTemplate(homeDir string) error {
	const op = "usecase.MakeDepsTemplate"

	log := u.log.With("operation", op)

	nameKey := "name"
	descKey := "desc"
	
	configs := []*entity.FieldConfig {
	    {Key: nameKey, Placeholder: "Name (required)", Required: true, Width: 32},
	    {Key: descKey, Placeholder: "Description (optional)", Required: false, Width: 64, CharLimit: 64},
	}
	
	variables, err := u.ui.DynamicInput("Create template", configs)
	if err != nil {
		log.Error("could not get template variables: %w", "error", err)

		return err
	}
	
	var deps []*entity.Dependency
	
	if homeDir == "" {
		deps, err = readDeps()
		if err != nil {
			log.Error("could not read dependecies", "error", err)

			return err
		}
	} else {
		deps, err = copyDeps(homeDir)
		if err != nil {
			log.Error("could not copy dependecies", "error", err)

			return err
		}
	}

	if !isUnique(u.cfg.Routes.DepsDir, variables[nameKey]) {
		err := u.confirmStatus(variables[nameKey])
		if err != nil {
			return err
		}
	}
	
	template := &entity.Template{
		Name: variables[nameKey],
		Description: variables[descKey],
		Dependencies: deps,
	}
	
	err = createDepsFile(u.cfg.Routes.DepsDir, template)
	if err != nil {
		log.Error("could not create template", "error", err)

		return err
	}
	
	return nil
}

func (u *UseCase) confirmStatus(name string) error {
	isIt, err := u.ui.ConfirmOverwrite(name)
	if err != nil {
		return err
	}

	if !isIt {
		err := u.ui.ShowStatus("canceling", 500*time.Microsecond)
		if err != nil {
			return err
		}
		
		return nil
	}

	err = u.ui.ShowStatus("overwriting", 1*time.Second)
	if err != nil {
		return err
	}

	return nil
}

func copyDeps(path string) ([]*entity.Dependency, error) {
	cmd := exec.Command("go", "list", "-m", "all")
	cmd.Dir = path

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("could not read dependencies: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	var deps []*entity.Dependency

	for _, line := range lines {
		dep, err := parseDependencyLine(line)
		if err != nil {
			continue
		}

		if !strings.Contains(dep.URL, ".") && dep.URL != "go" && dep.URL != "os" {
			continue
		}

		deps = append(deps, dep)
	}

	return deps, nil
}

func parseDependencyLine(line string) (*entity.Dependency, error) {
	fields := strings.Fields(line)
	if len(fields) != 2 {
		return nil, fmt.Errorf("invalid dependency line: %q", line)
	}

	url := fields[0]
	version := fields[1]

	nameParts := strings.Split(url, "/")
	name := strings.Join(nameParts[1:], "/")

	return &entity.Dependency{Name: name, URL: url, Version: version}, nil
}

func readDeps() ([]*entity.Dependency, error) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Please, write dependecies like 'github.com/spf13/cobra v1.10.2' and separate them with enter: ")

	var dependencies []*entity.Dependency

	for scanner.Scan() && scanner.Text() != "" {
		dep, err := parseDependencyLine(scanner.Text())
		if err != nil {
			return nil, err
		}

		dependencies = append(dependencies, dep)
	}
	err := scanner.Err()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading standard input:", err)
	}
	
	return dependencies, nil
}

func createDepsFile(path string, deps *entity.Template) error {
	data, err := json.Marshal(deps)
	if err != nil {
		return fmt.Errorf("could not marshal data. error: %w", err)
	}

	name := deps.Name + ".json"

	fullPath := filepath.Join(path, name)

	err = os.WriteFile(fullPath, data, 0600)
	if err != nil {
		return fmt.Errorf("could not create file. error: %w", err)
	}
	
	return nil
}
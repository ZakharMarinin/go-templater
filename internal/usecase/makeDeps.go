package usecase

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go-templater/internal/domain/entity"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func (u *UseCase) MakeDepsTemplate(homeDir string) error {
	const op = "usecase.MakeDepsTemplate"

	log := u.log.With("operation", op)

	variables, err := u.ui.Input()
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

	if !isUnique(u.cfg.Routes.DepsDir, variables.Name) {
		err := u.confirmStatus(variables.Name)
		if err != nil {
			return err
		}
	}
	
	template := &entity.Template{
		Name: variables.Name,
		Description: variables.Description,
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
		parts := strings.Fields(line)
		if len(parts) > 0 {
			if len(parts) != 2 {
				continue
			}
			
			version := parts[1]
			
			nameParts := strings.Split(parts[0], "/")

			name := strings.Join(nameParts[1:], "/")
			url := strings.Join(nameParts[0:], "/")
			
			dep := entity.Dependency{
				Name: name,
				URL: url,
				Version: version,
			}

			if !strings.Contains(url, ".") && url != "go" && url != "os" {
				continue
			}

			deps = append(deps, &dep)
		}
	}
	
	return deps, nil
}

func readDeps() ([]*entity.Dependency, error) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Please, write dependecies like 'github.com/spf13/cobra v1.10.2' and separate them with enter: ")

	var dependencies []*entity.Dependency
	
	for scanner.Scan() && scanner.Text() != "" {
		line := scanner.Text()
		
		parts := strings.Split(line, " ")

		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid dependency")
		}
		
		version := parts[1]
		
		nameParts := strings.Split(parts[0], "/")

		name := strings.Join(nameParts[1:], "/")
		url := strings.Join(nameParts[0:], "/")
		
		dep := entity.Dependency{
			Name: name,
			URL: url,
			Version: version,
		}

		dependencies = append(dependencies, &dep)
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
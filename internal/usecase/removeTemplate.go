package usecase

import (
	"fmt"
	"os"
	"path/filepath"
)

func (u *UseCase) RemoveTemplate(templateType string) error {
	const op = "usecase.RemoveTemplate"
	
	log := u.log.With("operation", op, "type", templateType)

	var targetDir string
	switch templateType {
	case "struct":
		targetDir = u.cfg.Routes.StructsDir
	case "deps":
		targetDir = u.cfg.Routes.DepsDir
	default:
		log.Warn("Unknown template type")
	
		return fmt.Errorf("unknown template type: %s", templateType)
	}

	templates, err := getTemplates(targetDir)
	if err != nil {
		log.Error("could not get templates for removal", "error", err)
		
		return err
	}

	if len(templates) == 0 {
		fmt.Printf("No templates found in '%s' category.\n", templateType)
		
		return nil
	}

	selectedTemplate, err := u.ui.Select(templates)
	if err != nil {
		log.Error("error or cancellation during template selection", "error", err)
		
		return err
	}

	fileName := fmt.Sprintf("%s.json", selectedTemplate.Name)
	filePath := filepath.Join(targetDir, fileName)

	err = os.Remove(filePath)
	if err != nil {
		log.Error("failed to delete template file", "path", filePath, "error", err)
		fmt.Printf("Failed to delete template: %v\n", err)
		
		return err
	}

	fmt.Printf("Template '%s' successfully removed from %s!\n", selectedTemplate.Name, templateType)

	return nil
}
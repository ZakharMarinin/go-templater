package usecase

import (
	"go-templater/internal/domain/entity"
	"os"
)

func (u *UseCase) ListTemplates() error {
	const op = "usecase.ListTemplates"

	log := u.log.With("operation", op)

	var list []*entity.TemplateInfo

	structs, err := scanDir(u.cfg.Routes.StructsDir, "struct")
	if err != nil {
		log.Error("could not read struct templates", "error", err)

		return err
	}

	list = append(list, structs...)

	deps, err := scanDir(u.cfg.Routes.DepsDir, "deps")
	if err != nil {
		log.Error("could not read deps templates", "error", err)

		return err
	}

	list = append(list, deps...)

	u.ui.ShowTemplatesTable(list)

	return nil
}

func scanDir(dirPath string, tType string) ([]*entity.TemplateInfo, error) {
	var templates []*entity.TemplateInfo

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		templates = append(templates, &entity.TemplateInfo{
			Name:      file.Name(),
			Type:      tType,
			Size:      info.Size(),
			CreatedAt: info.ModTime(),
		})
	}

	return templates, nil
}

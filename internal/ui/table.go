package ui

import (
	"fmt"
	"go-templater/internal/domain/entity"
	"strconv"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

func (ui *UI) ShowTemplatesTable(templates []*entity.TemplateInfo) {
	if len(templates) == 0 {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("No templates found. Use 'make' to create one."))

		return
	}

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ui.cfg.Colors.Highlight)).
		Bold(true).
		Padding(0, 1)

	cellStyle := lipgloss.NewStyle().Padding(0, 1)

	var rows [][]string
	for _, t := range templates {
		sizeStr := fmt.Sprintf("%.2f KB", float64(t.Size)/1024.0)
		if t.Size < 1024 {
			sizeStr = strconv.FormatInt(t.Size, 10) + " B"
		}

		rows = append(rows, []string{
			t.Name,
			t.Type,
			sizeStr,
			t.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
		Headers("TEMPLATE NAME", "TYPE", "SIZE", "CREATED AT").
		Rows(rows...)

	t.StyleFunc(func(row, col int) lipgloss.Style {
		if row == 0 {
			return headerStyle
		}
		
		return cellStyle
	})

	fmt.Println(t.Render())
}
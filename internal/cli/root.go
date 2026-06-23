package cli

import (
	"fmt"
	"go-templater/internal/config"
	"os"

	"github.com/spf13/cobra"
)

const (
	structs = "struct"
	deps = "deps"
)

type UseCase interface {
	MakeStructTemplate(homeDir string) error
	MakeDepsTemplate(homeDir string) error
	InsertDirTemplate(homeDir string) error
	InsertDepsTemplate(homeDir string) error
	RemoveTemplate(removeType string) error
}

type Cobra struct {
	rootCmd *cobra.Command
	uc      UseCase
}

func New(cfg *config.Config, uc UseCase) *Cobra {
	c := &Cobra{
		uc: uc,
	}

	c.rootCmd = &cobra.Command{
        Use:   "go-templater",
        Short: "A CLI tool to manage Go project architecture templates",
        Long: `go-templater allows you to create and insert project structures and dependencies quickly.

Example:
  go-templater make struct --dir ./my-app
  go-templater insert struct --dir ./new-project`,
    }

	c.setupCommands()

	return c
}

func (c *Cobra) Execute() {
	err := c.rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (c *Cobra) setupCommands() {
	var makeDir string
	var insertDir string

	makeCmd := &cobra.Command{
		Use:   "make",
		Short: "Commands for creating templates",
	}

	insertCmd := &cobra.Command{
		Use:   "insert",
		Short: "Commands for inserting templates into projects",
	}

	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Commands for removing templates",
	}

	removeStructCmd := &cobra.Command{
		Use:   structs,
		Short: "Select and remove a structure template",
		Run: func(cmd *cobra.Command, args []string) {
			err := c.uc.RemoveTemplate("struct")
			if err != nil {
				os.Exit(1)
			}
		},
	}

	removeDepsCmd := &cobra.Command{
		Use:   deps,
		Short: "Select and remove a dependencies template",
		Run: func(cmd *cobra.Command, args []string) {
			err := c.uc.RemoveTemplate("deps")
			if err != nil {
				os.Exit(1)
			}
		},
	}

	makeCmd.PersistentFlags().StringVarP(&makeDir, "dir", "d", "", "Directory to read from")
	insertCmd.PersistentFlags().StringVarP(&insertDir, "dir", "d", "", "Target directory to insert into")

	makeStructCmd := &cobra.Command{
		Use:   structs,
		Short: "Create a structure template from the current directory",
		Run: func(cmd *cobra.Command, args []string) {
			dir := resolveDirectory(makeDir)

			err := c.uc.MakeStructTemplate(dir)
			if err != nil {
				os.Exit(1)
			}
		},
	}

	makeDepsCmd := &cobra.Command{
		Use:   deps,
		Short: "Create a dependencies template from go.mod",
		Run: func(cmd *cobra.Command, args []string) {
			dir := resolveDirectory(makeDir)

			err := c.uc.MakeDepsTemplate(dir)
			if err != nil {
				os.Exit(1)
			}
		},
	}

	insertStructCmd := &cobra.Command{
		Use:   structs,
		Short: "Insert a structure template into a directory",
		Run: func(cmd *cobra.Command, args []string) {
			dir := resolveDirectory(insertDir)

			err := c.uc.InsertDirTemplate(dir)
			if err != nil {
				os.Exit(1)
			}
		},
	}

	insertDepsCmd := &cobra.Command{
		Use:   deps,
		Short: "Insert dependencies from a template into a directory",
		Run: func(cmd *cobra.Command, args []string) {
			dir := resolveDirectory(insertDir)

			err := c.uc.InsertDepsTemplate(dir)
			if err != nil {
				os.Exit(1)
			}
		},
	}

	removeCmd.AddCommand(removeStructCmd, removeDepsCmd)
	makeCmd.AddCommand(makeStructCmd, makeDepsCmd)
	insertCmd.AddCommand(insertStructCmd, insertDepsCmd)

	c.rootCmd.AddCommand(makeCmd, insertCmd, removeCmd)
}

func resolveDirectory(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("❌ Error getting current working directory: %v\n", err)
		os.Exit(1)
	}

	return currentDir
}

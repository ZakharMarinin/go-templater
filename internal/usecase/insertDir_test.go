package usecase

import (
	"github.com/ZakharMarinin/go-templater/internal/domain/entity"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInsertDirs(t *testing.T) {
	dir := t.TempDir()

	nodes := []*entity.Node{
		{
			Name: "sub", Path: "/sub", IsDir: true,
			Children: []*entity.Node{
				{Name: "file.go", Path: "/sub/file.go", IsDir: false, Content: "package sub"},
			},
		},
		{Name: "README.md", Path: "/README.md", IsDir: false, Content: "hello"},
		{Name: "empty.go", Path: "/empty.go", IsDir: false},
	}

	err := insertDirs(nodes, dir)
	if err != nil {
		t.Fatalf("insertDirs() error = %v", err)
	}

	assertFileContent := func(path, want string) {
		t.Helper()

		data, err := os.ReadFile(filepath.Join(dir, path))
		if err != nil {
			t.Fatalf("could not read %s: %v", path, err)
		}

		if string(data) != want {
			t.Errorf("%s content = %q, want %q", path, data, want)
		}
	}

	assertFileContent("sub/file.go", "package sub")
	assertFileContent("README.md", "hello")
	assertFileContent("empty.go", "")
}

func TestCreaProject(t *testing.T) {
	requireGoToolchain(t)

	dir := t.TempDir()

	err := creaProject(dir, "myproj")
	if err != nil {
		t.Fatalf("creaProject() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		t.Fatalf("could not read go.mod: %v", err)
	}

	if !strings.Contains(string(data), "module myproj") {
		t.Errorf("go.mod = %q, want it to contain %q", data, "module myproj")
	}
}

func TestCreaProjectInvalidDir(t *testing.T) {
	err := creaProject(filepath.Join(t.TempDir(), "does-not-exist"), "myproj")
	if err == nil {
		t.Fatal("expected error for a non-existent directory, got nil")
	}
}

func TestInsertDirTemplate(t *testing.T) {
	requireGoToolchain(t)

	dstDir := t.TempDir()

	tpl := &entity.Template{
		Name: "tpl",
		Nodes: []*entity.Node{
			{Name: "main.go", Path: "/main.go", IsDir: false, Content: "package main"},
		},
	}

	ui := &fakeUI{
		selectTemplate: tpl,
		dynamicInputs:  map[string]string{"project": "myproj"},
	}

	uc, _ := newTestUseCase(t, ui)

	err := uc.InsertDirTemplate(dstDir)
	if err != nil {
		t.Fatalf("InsertDirTemplate() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dstDir, "main.go"))
	if err != nil || string(data) != "package main" {
		t.Errorf("main.go content = %q, err = %v, want %q", data, err, "package main")
	}

	_, err = os.Stat(filepath.Join(dstDir, "go.mod"))
	if err != nil {
		t.Errorf("expected go.mod to be created: %v", err)
	}
}

func TestInsertDirTemplateSelectError(t *testing.T) {
	ui := &fakeUI{selectErr: os.ErrInvalid}
	uc, _ := newTestUseCase(t, ui)

	err := uc.InsertDirTemplate(t.TempDir())
	if err == nil {
		t.Fatal("expected error when template selection fails, got nil")
	}
}

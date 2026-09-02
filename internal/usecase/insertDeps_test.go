package usecase

import (
	"errors"
	"github.com/ZakharMarinin/go-templater/internal/domain/entity"
	"github.com/ZakharMarinin/go-templater/pkg/response"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestInsertDepsNoGoMod(t *testing.T) {
	err := insertDeps(nil, t.TempDir())
	if !errors.Is(err, response.ErrNotExist) {
		t.Fatalf("insertDeps() error = %v, want %v", err, response.ErrNotExist)
	}
}

func TestInsertDepsEmptyList(t *testing.T) {
	requireGoToolchain(t)

	t.Setenv("GOPROXY", "off")
	t.Setenv("GOSUMDB", "off")

	dir := t.TempDir()

	cmd := exec.Command("go", "mod", "init", "testmod")
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go mod init failed: %v, out: %s", err, out)
	}

	// `go get` with no packages needs at least one buildable package in the
	// module, otherwise it errors with "no package to get in current directory".
	err = os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0600)
	if err != nil {
		t.Fatalf("could not write fixture: %v", err)
	}

	err = insertDeps(nil, dir)
	if err != nil {
		t.Fatalf("insertDeps() error = %v", err)
	}
}

func TestInsertDepsTemplateNoGoMod(t *testing.T) {
	dstDir := t.TempDir()

	tpl := &entity.Template{Name: "tpl"}
	ui := &fakeUI{selectTemplate: tpl}
	uc, _ := newTestUseCase(t, ui)

	err := uc.InsertDepsTemplate(dstDir)
	if !errors.Is(err, response.ErrNotExist) {
		t.Fatalf("InsertDepsTemplate() error = %v, want %v", err, response.ErrNotExist)
	}
}

func TestInsertDepsTemplateSelectError(t *testing.T) {
	ui := &fakeUI{selectErr: response.ErrCanceled}
	uc, _ := newTestUseCase(t, ui)

	err := uc.InsertDepsTemplate(t.TempDir())
	if !errors.Is(err, response.ErrCanceled) {
		t.Fatalf("InsertDepsTemplate() error = %v, want %v", err, response.ErrCanceled)
	}
}

func TestInsertDepsTemplateGetTemplatesError(t *testing.T) {
	ui := &fakeUI{}
	uc, _ := newTestUseCase(t, ui)
	uc.cfg.Routes.DepsDir = filepath.Join(t.TempDir(), "does-not-exist")

	err := uc.InsertDepsTemplate(t.TempDir())
	if err == nil {
		t.Fatal("expected error when deps templates directory is missing, got nil")
	}
}

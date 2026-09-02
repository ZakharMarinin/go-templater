package usecase

import (
	"encoding/json"
	"errors"
	"github.com/ZakharMarinin/go-templater/internal/domain/entity"
	"github.com/ZakharMarinin/go-templater/pkg/response"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestParseDependencyLine(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		wantErr     bool
		wantName    string
		wantURL     string
		wantVersion string
	}{
		{
			name: "valid", line: "github.com/spf13/cobra v1.10.2",
			wantName: "spf13/cobra", wantURL: "github.com/spf13/cobra", wantVersion: "v1.10.2",
		},
		{
			name: "single segment path", line: "os v0.0.0",
			wantName: "", wantURL: "os", wantVersion: "v0.0.0",
		},
		{
			name: "collapses extra whitespace", line: "  github.com/a/b   v1.0.0  ",
			wantName: "a/b", wantURL: "github.com/a/b", wantVersion: "v1.0.0",
		},
		{name: "missing version", line: "github.com/spf13/cobra", wantErr: true},
		{name: "empty line", line: "", wantErr: true},
		{name: "too many fields", line: "a b c", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dep, err := parseDependencyLine(tt.line)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("parseDependencyLine(%q) expected error, got dep=%+v", tt.line, dep)
				}

				return
			}

			if err != nil {
				t.Fatalf("parseDependencyLine(%q) unexpected error: %v", tt.line, err)
			}

			if dep.Name != tt.wantName || dep.URL != tt.wantURL || dep.Version != tt.wantVersion {
				t.Errorf("parseDependencyLine(%q) = %+v, want name=%q url=%q version=%q",
					tt.line, dep, tt.wantName, tt.wantURL, tt.wantVersion)
			}
		})
	}
}

func TestCreateDepsFile(t *testing.T) {
	dir := t.TempDir()

	tpl := &entity.Template{
		Name:        "my-deps",
		Description: "desc",
		Dependencies: []*entity.Dependency{
			{Name: "cobra", URL: "github.com/spf13/cobra", Version: "v1.10.2"},
		},
	}

	err := createDepsFile(dir, tpl)
	if err != nil {
		t.Fatalf("createDepsFile() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "my-deps.json"))
	if err != nil {
		t.Fatalf("could not read created file: %v", err)
	}

	var got entity.Template

	err = json.Unmarshal(data, &got)
	if err != nil {
		t.Fatalf("could not unmarshal created file: %v", err)
	}

	if got.Name != tpl.Name || len(got.Dependencies) != 1 || got.Dependencies[0].URL != tpl.Dependencies[0].URL {
		t.Errorf("got %+v, want %+v", got, tpl)
	}
}

func withStdin(t *testing.T, content string) {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("could not create pipe: %v", err)
	}

	origStdin := os.Stdin
	os.Stdin = r

	t.Cleanup(func() { os.Stdin = origStdin })

	go func() {
		_, _ = io.WriteString(w, content)

		_ = w.Close()
	}()
}

func TestReadDeps(t *testing.T) {
	withStdin(t, "github.com/spf13/cobra v1.10.2\ngithub.com/a/b v1.0.0\n\n")

	deps, err := readDeps()
	if err != nil {
		t.Fatalf("readDeps() error = %v", err)
	}

	if len(deps) != 2 {
		t.Fatalf("expected 2 deps, got %d: %+v", len(deps), deps)
	}
}

func TestReadDepsInvalidLine(t *testing.T) {
	withStdin(t, "not-a-valid-dependency-line\n")

	_, err := readDeps()
	if err == nil {
		t.Fatal("expected error for invalid dependency line, got nil")
	}
}

func TestReadDepsEmptyInput(t *testing.T) {
	withStdin(t, "\n")

	deps, err := readDeps()
	if err != nil {
		t.Fatalf("readDeps() error = %v", err)
	}

	if len(deps) != 0 {
		t.Errorf("expected no deps, got %d", len(deps))
	}
}

func TestCopyDeps(t *testing.T) {
	requireGoToolchain(t)

	dir := t.TempDir()

	cmd := exec.Command("go", "mod", "init", "testmod")
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go mod init failed: %v, out: %s", err, out)
	}

	deps, err := copyDeps(dir)
	if err != nil {
		t.Fatalf("copyDeps() error = %v", err)
	}

	if len(deps) != 0 {
		t.Errorf("expected no dependencies for a bare module, got %+v", deps)
	}
}

func TestCopyDepsInvalidDir(t *testing.T) {
	_, err := copyDeps(filepath.Join(t.TempDir(), "does-not-exist"))
	if err == nil {
		t.Fatal("expected error for a non-module directory, got nil")
	}
}

func TestMakeDepsTemplateFromStdin(t *testing.T) {
	ui := &fakeUI{dynamicInputs: map[string]string{"name": "stdin-deps", "desc": "from stdin"}}
	uc, _ := newTestUseCase(t, ui)

	withStdin(t, "github.com/spf13/cobra v1.10.2\n\n")

	err := uc.MakeDepsTemplate("")
	if err != nil {
		t.Fatalf("MakeDepsTemplate() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(uc.cfg.Routes.DepsDir, "stdin-deps.json"))
	if err != nil {
		t.Fatalf("could not read created template: %v", err)
	}

	var got entity.Template

	err = json.Unmarshal(data, &got)
	if err != nil {
		t.Fatalf("could not unmarshal created template: %v", err)
	}

	if len(got.Dependencies) != 1 || got.Dependencies[0].URL != "github.com/spf13/cobra" {
		t.Errorf("unexpected dependencies: %+v", got.Dependencies)
	}
}

func TestMakeDepsTemplateCanceledInput(t *testing.T) {
	ui := &fakeUI{dynamicInputErr: response.ErrCanceled}
	uc, _ := newTestUseCase(t, ui)

	err := uc.MakeDepsTemplate("")
	if !errors.Is(err, response.ErrCanceled) {
		t.Fatalf("MakeDepsTemplate() error = %v, want %v", err, response.ErrCanceled)
	}
}

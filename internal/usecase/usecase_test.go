package usecase

import (
	"github.com/ZakharMarinin/go-templater/internal/config"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}

func newTestUseCase(t *testing.T, ui UI) (*UseCase, string) {
	t.Helper()

	structsDir := t.TempDir()
	depsDir := t.TempDir()
	logsDir := t.TempDir()

	cfg := &config.Config{
		Env: "local",
		Routes: config.Routes{
			StructsDir: structsDir,
			DepsDir:    depsDir,
			LogsDir:    logsDir,
		},
	}

	return New(cfg, discardLogger(), ui), structsDir
}

func TestIsUnique(t *testing.T) {
	dir := t.TempDir()

	err := os.WriteFile(filepath.Join(dir, "existing.json"), []byte("{}"), 0600)
	if err != nil {
		t.Fatalf("could not write fixture: %v", err)
	}

	tests := []struct {
		name string
		want bool
	}{
		{name: "existing", want: false},
		{name: "brand-new", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isUnique(dir, tt.name)
			if got != tt.want {
				t.Errorf("isUnique(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestGetTemplates(t *testing.T) {
	dir := t.TempDir()

	writeJSON := func(name, content string) {
		err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0600)
		if err != nil {
			t.Fatalf("could not write fixture: %v", err)
		}
	}

	writeJSON("one.json", `{"name":"one","description":"first"}`)
	writeJSON("two.json", `{"name":"two","description":"second"}`)

	templates, err := getTemplates(dir)
	if err != nil {
		t.Fatalf("getTemplates() error = %v", err)
	}

	if len(templates) != 2 {
		t.Fatalf("expected 2 templates, got %d", len(templates))
	}

	names := map[string]bool{}
	for _, tpl := range templates {
		names[tpl.Name] = true
	}

	if !names["one"] || !names["two"] {
		t.Errorf("expected templates 'one' and 'two', got %v", names)
	}
}

func TestGetTemplatesInvalidJSON(t *testing.T) {
	dir := t.TempDir()

	err := os.WriteFile(filepath.Join(dir, "broken.json"), []byte("not json"), 0600)
	if err != nil {
		t.Fatalf("could not write fixture: %v", err)
	}

	_, err = getTemplates(dir)
	if err == nil {
		t.Fatal("expected error for invalid json, got nil")
	}
}

func TestGetTemplatesMissingDir(t *testing.T) {
	_, err := getTemplates(filepath.Join(t.TempDir(), "does-not-exist"))
	if err == nil {
		t.Fatal("expected error for missing directory, got nil")
	}
}

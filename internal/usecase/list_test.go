package usecase

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanDir(t *testing.T) {
	dir := t.TempDir()

	err := os.WriteFile(filepath.Join(dir, "a.json"), []byte("{}"), 0600)
	if err != nil {
		t.Fatalf("could not write fixture: %v", err)
	}

	err = os.Mkdir(filepath.Join(dir, "subdir"), 0755)
	if err != nil {
		t.Fatalf("could not create fixture dir: %v", err)
	}

	infos, err := scanDir(dir, "struct")
	if err != nil {
		t.Fatalf("scanDir() error = %v", err)
	}

	if len(infos) != 1 {
		t.Fatalf("expected 1 entry (directories skipped), got %d", len(infos))
	}

	if infos[0].Name != "a.json" || infos[0].Type != "struct" {
		t.Errorf("got %+v, want name=a.json type=struct", infos[0])
	}
}

func TestScanDirMissing(t *testing.T) {
	_, err := scanDir(filepath.Join(t.TempDir(), "missing"), "struct")
	if err == nil {
		t.Fatal("expected error for a missing directory, got nil")
	}
}

func TestListTemplates(t *testing.T) {
	ui := &fakeUI{}
	uc, structsDir := newTestUseCase(t, ui)

	err := os.WriteFile(filepath.Join(structsDir, "s1.json"), []byte("{}"), 0600)
	if err != nil {
		t.Fatalf("could not write fixture: %v", err)
	}

	err = os.WriteFile(filepath.Join(uc.cfg.Routes.DepsDir, "d1.json"), []byte("{}"), 0600)
	if err != nil {
		t.Fatalf("could not write fixture: %v", err)
	}

	err = uc.ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates() error = %v", err)
	}

	if len(ui.shownTemplates) != 2 {
		t.Fatalf("expected 2 shown templates, got %d", len(ui.shownTemplates))
	}
}

func TestListTemplatesMissingStructsDir(t *testing.T) {
	ui := &fakeUI{}
	uc, _ := newTestUseCase(t, ui)
	uc.cfg.Routes.StructsDir = filepath.Join(t.TempDir(), "missing")

	err := uc.ListTemplates()
	if err == nil {
		t.Fatal("expected error when structs directory is missing, got nil")
	}
}

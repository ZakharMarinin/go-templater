package usecase

import (
	"github.com/ZakharMarinin/go-templater/internal/domain/entity"
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveTemplateStruct(t *testing.T) {
	ui := &fakeUI{selectTemplate: &entity.Template{Name: "tpl"}}
	uc, structsDir := newTestUseCase(t, ui)

	tplPath := filepath.Join(structsDir, "tpl.json")

	err := os.WriteFile(tplPath, []byte(`{"name":"tpl"}`), 0600)
	if err != nil {
		t.Fatalf("could not write fixture: %v", err)
	}

	err = uc.RemoveTemplate("struct")
	if err != nil {
		t.Fatalf("RemoveTemplate() error = %v", err)
	}

	_, statErr := os.Stat(tplPath)
	if !os.IsNotExist(statErr) {
		t.Errorf("expected template file to be removed, stat err = %v", statErr)
	}
}

func TestRemoveTemplateUnknownType(t *testing.T) {
	ui := &fakeUI{}
	uc, _ := newTestUseCase(t, ui)

	err := uc.RemoveTemplate("bogus")
	if err == nil {
		t.Fatal("expected error for unknown template type, got nil")
	}
}

func TestRemoveTemplateNoTemplatesFound(t *testing.T) {
	ui := &fakeUI{}
	uc, _ := newTestUseCase(t, ui)

	err := uc.RemoveTemplate("deps")
	if err != nil {
		t.Fatalf("RemoveTemplate() error = %v, want nil for an empty directory", err)
	}
}

func TestRemoveTemplateSelectError(t *testing.T) {
	ui := &fakeUI{selectErr: os.ErrInvalid}
	uc, structsDir := newTestUseCase(t, ui)

	err := os.WriteFile(filepath.Join(structsDir, "tpl.json"), []byte(`{"name":"tpl"}`), 0600)
	if err != nil {
		t.Fatalf("could not write fixture: %v", err)
	}

	err = uc.RemoveTemplate("struct")
	if err == nil {
		t.Fatal("expected error when selection fails, got nil")
	}
}

package usecase

import (
	"encoding/json"
	"github.com/ZakharMarinin/go-templater/internal/domain/entity"
	"github.com/ZakharMarinin/go-templater/pkg/response"
	"os"
	"path/filepath"
	"testing"
)

func TestMakeStructTemplate(t *testing.T) {
	srcDir := t.TempDir()

	err := os.WriteFile(filepath.Join(srcDir, "main.go"), []byte("package main"), 0600)
	if err != nil {
		t.Fatalf("could not write fixture: %v", err)
	}

	ui := &fakeUI{dynamicInputs: map[string]string{"name": "struct-tpl", "desc": "d"}}
	uc, structsDir := newTestUseCase(t, ui)

	err = uc.MakeStructTemplate(srcDir)
	if err != nil {
		t.Fatalf("MakeStructTemplate() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(structsDir, "struct-tpl.json"))
	if err != nil {
		t.Fatalf("could not read created template: %v", err)
	}

	var got entity.Template

	err = json.Unmarshal(data, &got)
	if err != nil {
		t.Fatalf("could not unmarshal created template: %v", err)
	}

	if len(got.Nodes) != 1 || got.Nodes[0].Content != "package main" {
		t.Errorf("got nodes %+v, want a single main.go node with content 'package main'", got.Nodes)
	}
}

func TestMakeStructTemplateCanceledInput(t *testing.T) {
	ui := &fakeUI{dynamicInputErr: response.ErrCanceled}
	uc, structsDir := newTestUseCase(t, ui)

	err := uc.MakeStructTemplate(t.TempDir())
	if err != nil {
		t.Fatalf("MakeStructTemplate() error = %v, want nil on cancel", err)
	}

	entries, err := os.ReadDir(structsDir)
	if err != nil {
		t.Fatalf("could not read structs dir: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("expected no template file to be created on cancel, found %d", len(entries))
	}
}

func TestMakeStructTemplateCanceledExclusions(t *testing.T) {
	srcDir := t.TempDir()

	err := os.WriteFile(filepath.Join(srcDir, "a.txt"), []byte("x"), 0600)
	if err != nil {
		t.Fatalf("could not write fixture: %v", err)
	}

	ui := &fakeUI{
		dynamicInputs: map[string]string{"name": "tpl2", "desc": ""},
		exclusionsErr: response.ErrCanceled,
	}
	uc, structsDir := newTestUseCase(t, ui)

	err = uc.MakeStructTemplate(srcDir)
	if err != nil {
		t.Fatalf("MakeStructTemplate() error = %v, want nil on cancel", err)
	}

	_, statErr := os.Stat(filepath.Join(structsDir, "tpl2.json"))
	if !os.IsNotExist(statErr) {
		t.Errorf("expected no template file to be created on cancel, stat err = %v", statErr)
	}
}

// TestMakeStructTemplateDeclinedOverwriteStillWrites documents existing behavior:
// confirmStatus only surfaces an error when ConfirmOverwrite itself errors, never
// when the user declines (isIt == false) - so a decline logs "canceling" but the
// template is written anyway. This looks like a pre-existing bug in confirmStatus.
func TestMakeStructTemplateDeclinedOverwriteStillWrites(t *testing.T) {
	srcDir := t.TempDir()

	err := os.WriteFile(filepath.Join(srcDir, "a.txt"), []byte("x"), 0600)
	if err != nil {
		t.Fatalf("could not write fixture: %v", err)
	}

	ui := &fakeUI{
		dynamicInputs:    map[string]string{"name": "dup", "desc": ""},
		confirmOverwrite: false,
	}
	uc, structsDir := newTestUseCase(t, ui)

	err = os.WriteFile(filepath.Join(structsDir, "dup.json"), []byte(`{"name":"dup"}`), 0600)
	if err != nil {
		t.Fatalf("could not write fixture: %v", err)
	}

	err = uc.MakeStructTemplate(srcDir)
	if err != nil {
		t.Fatalf("MakeStructTemplate() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(structsDir, "dup.json"))
	if err != nil {
		t.Fatalf("could not read template: %v", err)
	}

	var got entity.Template

	err = json.Unmarshal(data, &got)
	if err != nil {
		t.Fatalf("could not unmarshal template: %v", err)
	}

	if len(got.Nodes) != 1 {
		t.Errorf("expected the template to be overwritten with new nodes, got %+v", got.Nodes)
	}
}

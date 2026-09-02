package usecase

import (
	"github.com/ZakharMarinin/go-templater/internal/domain/entity"
	"os/exec"
	"testing"
	"time"
)

func requireGoToolchain(t *testing.T) {
	t.Helper()

	_, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain not available")
	}
}

type fakeUI struct {
	selectTemplate *entity.Template
	selectErr      error

	dynamicInputs   map[string]string
	dynamicInputErr error

	confirmOverwrite bool
	confirmErr       error

	exclusionsFn  func(nodes []*entity.Node) error
	exclusionsErr error

	statusCalls []string

	shownTemplates []*entity.TemplateInfo

	spinnerErr error
}

func (f *fakeUI) Select(_ []*entity.Template) (*entity.Template, error) {
	return f.selectTemplate, f.selectErr
}

func (f *fakeUI) DynamicInput(_ string, _ []*entity.FieldConfig) (map[string]string, error) {
	return f.dynamicInputs, f.dynamicInputErr
}

func (f *fakeUI) ConfirmOverwrite(_ string) (bool, error) {
	return f.confirmOverwrite, f.confirmErr
}

func (f *fakeUI) ShowStatus(msg string, _ time.Duration) error {
	f.statusCalls = append(f.statusCalls, msg)

	return nil
}

func (f *fakeUI) NewSpinner(_ string, task func() error) error {
	if f.spinnerErr != nil {
		return f.spinnerErr
	}

	return task()
}

func (f *fakeUI) ShowTemplatesTable(templates []*entity.TemplateInfo) {
	f.shownTemplates = templates
}

func (f *fakeUI) SelectContentExclusions(nodes []*entity.Node) error {
	if f.exclusionsFn != nil {
		return f.exclusionsFn(nodes)
	}

	return f.exclusionsErr
}

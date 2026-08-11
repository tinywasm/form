package form_test

import (
	"strings"
	"testing"

	"github.com/tinywasm/form"
	"github.com/tinywasm/input"
	"github.com/tinywasm/model"
)

type anatomyStruct struct {
	model.Fielder
	Nombre string
	Gender string
}

func (s *anatomyStruct) Schema() []model.Field {
	return []model.Field{
		{Name: "nombre", NotNull: true, Type: input.Text()},
		{Name: "gender", Type: input.Radio()},
	}
}

func (s *anatomyStruct) Pointers() []any { return []any{&s.Nombre, &s.Gender} }
func (s *anatomyStruct) Values() []any   { return []any{s.Nombre, s.Gender} }

func TestForm_Anatomy(t *testing.T) {
	s := &anatomyStruct{}
	f, _ := form.New("app", s)
	html := f.String()

	// 1. Root container class should be 'tw-field'
	if !strings.Contains(html, "class='tw-field'") {
		t.Errorf("Expected html to contain root class 'tw-field', got: %s", html)
	}

	// 2. Label class should be 'tw-field__label'
	if !strings.Contains(html, "class='tw-field__label'") {
		t.Errorf("Expected html to contain label class 'tw-field__label', got: %s", html)
	}

	// 3. Error span class should be 'tw-field__error'
	if !strings.Contains(html, "class='tw-field__error'") {
		t.Errorf("Expected html to contain error class 'tw-field__error', got: %s", html)
	}

	// 4. Radio group should contain 'tw-field__radio-group'
	if !strings.Contains(html, "class='tw-field__radio-group'") {
		t.Errorf("Expected html to contain radio-group class 'tw-field__radio-group', got: %s", html)
	}
}

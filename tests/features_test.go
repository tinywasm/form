package form_test

import (
	"testing"

	"github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/form"
	"github.com/tinywasm/form/input"
)

func TestForm_SetClass(t *testing.T) {
	s := &submitStruct{Nombre: "Jules"}

	t.Run("Per-form class only", func(t *testing.T) {
		f, _ := form.New("app", s)
		f.SetClass("cms-form")
		html := f.String()
		if !fmt.Contains(html, "class='cms-form'") {
			t.Errorf("Expected class='cms-form' in HTML, got %s", html)
		}
	})

	t.Run("Global + per-form class", func(t *testing.T) {
		form.SetGlobalClass("base-form")
		f, _ := form.New("app", s)
		f.SetClass("extra-form")
		html := f.String()
		if !fmt.Contains(html, "class='base-form extra-form'") {
			t.Errorf("Expected combined classes, got %s", html)
		}
	})
}

type customInput struct {
	input.Base
}

func (c *customInput) Clone(parentID, name string) input.Input {
	cl := *c
	cl.InitBase(parentID, name, c.HTMLName())
	return &cl
}

func (c *customInput) RenderInput(value *dom.SignalString, onInput func(string)) *dom.Element {
	return dom.NewElement("div").Class("my-widget").Text("Custom Widget")
}

func TestForm_Renderer(t *testing.T) {
	ci := &customInput{}
	ci.InitBase("custom", "Custom Field", "text")

	fx := &kindFixture{inp: ci}
	f, _ := form.New("tid", fx)

	html := f.String()

	// Must contain custom markup
	if !fmt.Contains(html, "class='my-widget'") {
		t.Errorf("Expected 'my-widget' class, got %s", html)
	}

	// Must contain the standard error span
	if !fmt.Contains(html, "class='tw-field-error'") {
		t.Errorf("Expected 'tw-field-error' class, got %s", html)
	}
}

func TestForm_RendererValidation(t *testing.T) {
	// Custom input with a validation rule
	ci := &customInput{}
	ci.InitBase("custom", "Custom Field", "text")
	ci.SetRequired(true)

	fx := &kindFixture{inp: ci}
	f, _ := form.New("tid", fx)

	// Empty required field should fail validation
	f.SetValues("tfield", "")
	err := f.Validate()
	if err == nil {
		t.Error("Expected validation error for required custom field, got nil")
	}

	// Non-empty should pass
	f.SetValues("tfield", "valid")
	err = f.Validate()
	if err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}
}

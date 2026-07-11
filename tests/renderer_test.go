package form_test

import (
	"testing"

	"github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/form"
	"github.com/tinywasm/form/input"
)

type customInput struct {
	input.Base
}

func (c *customInput) Clone(parentID, name string) input.Input {
	nc := *c
	nc.InitBase(parentID, name, "text")
	return &nc
}

func (c *customInput) RenderInput(value *dom.SignalString, onInput func(string)) *dom.Element {
	return dom.NewElement("div").Class("my-widget")
}

func TestForm_Renderer(t *testing.T) {
	ci := &customInput{}
	fx := &kindFixture{inp: ci}

	f, _ := form.New("app", fx)
	html := f.String()

	// Must contain custom markup
	if !fmt.Contains(html, "class='my-widget'") {
		t.Errorf("Expected html to contain 'class='my-widget'', got: %s", html)
	}

	// Must still contain the standard error span
	if !fmt.Contains(html, "class='tw-field-error'") {
		t.Errorf("Expected html to contain 'class='tw-field-error'', got: %s", html)
	}
}

func TestForm_Renderer_Validation(t *testing.T) {
	// customInput inherits Validate from input.Base which returns nil
	// Let's make a custom input with validation
	ci := &customInputWithValidation{}
	fx := &kindFixture{inp: ci}

	f, _ := form.New("app", fx)

	f.SetValues("tfield", "invalid")
	err := f.Validate()
	if err == nil {
		t.Error("Expected validation error for 'invalid' value, got nil")
	}
}

type customInputWithValidation struct {
	input.Base
}

func (c *customInputWithValidation) Clone(parentID, name string) input.Input {
	nc := *c
	nc.InitBase(parentID, name, "text")
	return &nc
}

func (c *customInputWithValidation) RenderInput(value *dom.SignalString, onInput func(string)) *dom.Element {
	return dom.NewElement("div").Class("my-widget")
}

func (c *customInputWithValidation) Validate(val string) error {
	if val == "invalid" {
		return fmt.Err("tfield", "is invalid")
	}
	return nil
}

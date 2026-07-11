package form_test

import (
	"testing"

	"github.com/tinywasm/fmt"
	"github.com/tinywasm/form"
)

func TestForm_SetClass(t *testing.T) {
	s := &submitStruct{}
	f, _ := form.New("app", s)

	f.SetClass("cms-form")

	html := f.String()
	expected := "class='cms-form'"
	if !fmt.Contains(html, expected) {
		t.Errorf("Expected html to contain %q, got: %s", expected, html)
	}
}

func TestForm_SetClass_Append(t *testing.T) {
	form.SetGlobalClass("global-class")
	defer form.SetGlobalClass("") // Reset global state

	s := &submitStruct{}
	f, _ := form.New("app", s)
	f.SetClass("local-class")

	html := f.String()
	// New() uses globalClass as initial f.class. SetClass appends.
	// Initial f.class = "global-class"
	// After SetClass("local-class"), f.class = "global-class local-class"
	expected := "class='global-class local-class'"
	if !fmt.Contains(html, expected) {
		t.Errorf("Expected html to contain %q, got: %s", expected, html)
	}
}

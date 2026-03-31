// Package input tests.
// This file contains the shared registry and helpers used by all test files.
// To add a new input:
//  1. Add a case in buildInput()
//  2. Add validation cases in validation_test.go
//  3. Add render cases in render_test.go (if it has custom RenderHTML)
package input

import (
	"testing"

	"github.com/tinywasm/fmt"
	_ "github.com/tinywasm/fmt/dictionary"
)

// tc is a compact validation test case.
type tc struct {
	t    string // input type name (must match a case in buildInput)
	name string // subtest name
	val  string // input value
	err  string // expected error substring (empty = no error expected)
	opts []fmt.KeyValue
	req  bool
}

// rc is a compact render test case.
type rc struct {
	t       string // input type name
	name    string // subtest name
	values  []string
	opts    []fmt.KeyValue
	contain string // expected substring in HTML output
}

// Shared option sets used across test files.
var opts12 = []fmt.KeyValue{{Key: "1", Value: "Admin"}, {Key: "2", Value: "Editor"}}
var optsGender = []fmt.KeyValue{{Key: "m", Value: "Male"}, {Key: "f", Value: "Female"}}

// buildInput creates a fresh input instance by kind. Add new inputs here.
func buildInput(t *testing.T, kind string, opts []fmt.KeyValue) Input {
	t.Helper()
	id, name := "tid", "tfield"
	var inp Input
	switch kind {
	case "Address":
		inp = Address()
	case "Checkbox":
		inp = Checkbox()
	case "Datalist":
		dl := Datalist()
		if len(opts) > 0 {
			dl.(interface{ SetOptions(...fmt.KeyValue) }).SetOptions(opts...)
		}
		inp = dl
	case "Date":
		inp = Date()
	case "Email":
		inp = Email()
	case "Filepath":
		inp = Filepath()
	case "Gender":
		g := Gender()
		if len(opts) > 0 {
			g.(interface{ SetOptions(...fmt.KeyValue) }).SetOptions(opts...)
		}
		inp = g
	case "Hour":
		inp = Hour()
	case "IP":
		inp = IP()
	case "Number":
		inp = Number()
	case "Password":
		inp = Password()
	case "Phone":
		inp = Phone()
	case "Radio":
		r := Radio()
		if len(opts) > 0 {
			r.(interface{ SetOptions(...fmt.KeyValue) }).SetOptions(opts...)
		}
		inp = r
	case "Rut":
		inp = Rut()
	case "Select":
		s := Select()
		if len(opts) > 0 {
			s.(interface{ SetOptions(...fmt.KeyValue) }).SetOptions(opts...)
		}
		inp = s
	case "Text":
		inp = Text()
	case "Textarea":
		inp = Textarea()
	default:
		t.Fatalf("unknown input type: %q — add it to buildInput()", kind)
		return nil
	}
	return inp.Clone(id, name).(Input)
}

// checkErr asserts the error matches the expected substring (case-insensitive).
func checkErr(t *testing.T, err error, expected string) {
	t.Helper()
	if expected == "" {
		if err != nil {
			t.Errorf("expected no error, got %q", err.Error())
		}
		return
	}
	if err == nil {
		t.Errorf("expected error containing %q, got nil", expected)
		return
	}
	got := fmt.Convert(err.Error()).ToLower().String()
	exp := fmt.Convert(expected).ToLower().String()
	if !fmt.Contains(got, exp) {
		t.Errorf("expected error containing %q, got %q", expected, err.Error())
	}
}

func TestClone_Preservation(t *testing.T) {
	// Create a prototype with custom configuration
	proto := Text()
	if setter, ok := proto.(interface{ SetPlaceholder(string) }); ok {
		setter.SetPlaceholder("Custom Placeholder")
	}
	if setter, ok := proto.(interface{ SetTitle(string) }); ok {
		setter.SetTitle("Custom Title")
	}
	proto.AddAttribute("data-test", "value")
	proto.SetRequired(true)

	// Clone it
	cloned := proto.Clone("parent", "field").(Input)

	// Verify ID and name are updated
	if cloned.GetID() != "parent.field" {
		t.Errorf("Expected ID 'parent.field', got %q", cloned.GetID())
	}
	if cloned.FieldName() != "field" {
		t.Errorf("Expected name 'field', got %q", cloned.FieldName())
	}

	// Verify custom configuration is preserved
	if getter, ok := cloned.(interface{ GetPlaceholder() string }); ok {
		if getter.GetPlaceholder() != "Custom Placeholder" {
			t.Errorf("Expected placeholder 'Custom Placeholder', got %q", getter.GetPlaceholder())
		}
	}
	if getter, ok := cloned.(interface{ GetTitle() string }); ok {
		if getter.GetTitle() != "Custom Title" {
			t.Errorf("Expected title 'Custom Title', got %q", getter.GetTitle())
		}
	}

	// Verify attributes are preserved
	html := cloned.RenderHTML()
	if !fmt.Contains(html, `data-test="value"`) {
		t.Errorf("Expected attribute data-test=\"value\" in HTML, got %s", html)
	}
	if !fmt.Contains(html, `required`) {
		t.Errorf("Expected required attribute in HTML, got %s", html)
	}
}

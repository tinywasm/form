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
	switch kind {
	case "Address":
		return Address(id, name)
	case "Checkbox":
		return Checkbox(id, name)
	case "Datalist":
		dl := Datalist(id, "datalist_field")
		if len(opts) > 0 {
			dl.(*datalist).SetOptions(opts...)
		}
		return dl
	case "Date":
		return Date(id, name)
	case "Email":
		return Email(id, name)
	case "Filepath":
		return Filepath(id, name)
	case "Gender":
		g := Radio(id, name).(*radio)
		g.SetOptions(optsGender...)
		return g
	case "Hour":
		return Hour(id, name)
	case "IP":
		return IP(id, name)
	case "Number":
		return Number(id, name)
	case "Password":
		return Password(id, name)
	case "Phone":
		return Phone(id, name)
	case "Radio":
		r := Radio(id, name).(*radio)
		r.SetOptions(optsGender...)
		return r
	case "Rut":
		return Rut(id, name)
	case "Select":
		s := Select(id, name)
		if len(opts) > 0 {
			s.(*select_).SetOptions(opts...)
		}
		return s
	case "Text":
		return Text(id, name)
	case "Textarea":
		return Textarea(id, name)
	default:
		t.Fatalf("unknown input type: %q — add it to buildInput()", kind)
		return nil
	}
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

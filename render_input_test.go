package form

import (
	"testing"

	"github.com/tinywasm/fmt"
	"github.com/tinywasm/form/input"
)

// Test_Render verifies String output for inputs with custom rendering.
func Test_Render(t *testing.T) {
	cases := []rc{
		// ── Checkbox ─────────────────────────────────────────────────────────
		{
			t: "Checkbox", name: "renders checkbox input",
			contain: `type="checkbox"`,
		},

		// ── Datalist ─────────────────────────────────────────────────────────
		{
			t: "Datalist", name: "contains datalist element",
			opts:    opts12,
			contain: `<datalist`,
		},
		{
			t: "Datalist", name: "contains option values",
			opts:    opts12,
			contain: `value="1"`,
		},
		{
			t: "Datalist", name: "links input to datalist via list attribute",
			opts:    opts12,
			contain: `list="`,
		},

		// ── Radio ────────────────────────────────────────────────────────────
		{
			t: "Radio", name: "renders label per option",
			opts:    optsGender,
			contain: `<label>`,
		},
		{
			t: "Radio", name: "renders male option value",
			opts:    optsGender,
			contain: `value="m"`,
		},
		{
			t: "Radio", name: "renders checked when value matches",
			opts:    optsGender,
			values:  []string{"f"},
			contain: `checked`,
		},

		// ── Select ───────────────────────────────────────────────────────────
		{
			t: "Select", name: "renders select element",
			opts:    opts12,
			contain: `<select`,
		},
		{
			t: "Select", name: "renders option element",
			opts:    opts12,
			contain: `<option`,
		},
		{
			t: "Select", name: "renders selected when value matches",
			opts:    opts12,
			values:  []string{"2"},
			contain: `selected`,
		},

		// ── Textarea ─────────────────────────────────────────────────────────
		{
			t: "Textarea", name: "uses textarea tag",
			contain: `<textarea`,
		},
		{
			t: "Textarea", name: "renders value inside tag",
			values:  []string{"hello world"},
			contain: `hello world`,
		},

		// ── Rut ──────────────────────────────────────────────────────────────

		// ── Search ───────────────────────────────────────────────────────────
		{
			t: "Search", name: "renders search input",
			contain: `type="search"`,
		},

		// ── Text ─────────────────────────────────────────────────────────────
		{
			t: "Text", name: "renders input tag",
			contain: `<input`,
		},
		{
			t: "Text", name: "has type text",
			contain: `type="text"`,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.t+"/"+c.name, func(t *testing.T) {
			inp := buildInput(t, c.t, c.opts)
			if setter, ok := inp.(interface{ SetValues(...string) }); ok && len(c.values) > 0 {
				setter.SetValues(c.values...)
			}
			html := RenderInput(inp)
			if !fmt.Contains(html, c.contain) {
				t.Errorf("RenderInput() missing %q\ngot: %s", c.contain, html)
			}
		})
	}
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

// buildInput creates a fresh input instance by kind.
func buildInput(t *testing.T, kind string, opts []fmt.KeyValue) input.Input {
	t.Helper()
	id, name := "tid", "tfield"
	var inp input.Input
	switch kind {
	case "Address":
		inp = input.Address()
	case "Checkbox":
		inp = input.Checkbox()
	case "Datalist":
		dl := input.Datalist()
		if len(opts) > 0 {
			dl.(interface{ SetOptions(...fmt.KeyValue) }).SetOptions(opts...)
		}
		inp = dl
	case "Date":
		inp = input.Date()
	case "Email":
		inp = input.Email()
	case "Filepath":
		inp = input.Filepath()
	case "Gender":
		g := input.Gender()
		if len(opts) > 0 {
			g.(interface{ SetOptions(...fmt.KeyValue) }).SetOptions(opts...)
		}
		inp = g
	case "Hour":
		inp = input.Hour()
	case "IP":
		inp = input.IP()
	case "Number":
		inp = input.Number()
	case "Password":
		inp = input.Password()
	case "Phone":
		inp = input.Phone()
	case "Radio":
		r := input.Radio()
		if len(opts) > 0 {
			r.(interface{ SetOptions(...fmt.KeyValue) }).SetOptions(opts...)
		}
		inp = r
	case "Rut":
		inp = input.Rut()
	case "Search":
		inp = input.Search()
	case "Select":
		s := input.Select()
		if len(opts) > 0 {
			s.(interface{ SetOptions(...fmt.KeyValue) }).SetOptions(opts...)
		}
		inp = s
	case "Text":
		inp = input.Text()
	case "Textarea":
		inp = input.Textarea()
	default:
		t.Fatalf("unknown input type: %q", kind)
		return nil
	}
	return inp.Clone(id, name).(input.Input)
}

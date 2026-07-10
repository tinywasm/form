package form_test

import (
	"testing"

	"github.com/tinywasm/fmt"
	"github.com/tinywasm/form"
	"github.com/tinywasm/form/input"
	"github.com/tinywasm/model"
)

// kindFixture is a one-field Fielder used to render a single input kind
// through the real form pipeline.
type kindFixture struct {
	inp     input.Input
	valText string
	valInt  int64
	valBool bool
}

func (k *kindFixture) Schema() []model.Field {
	return []model.Field{{Name: "tfield", Type: k.inp}}
}

func (k *kindFixture) Pointers() []any {
	switch k.inp.Storage() {
	case model.FieldInt:
		return []any{&k.valInt}
	case model.FieldBool:
		return []any{&k.valBool}
	default:
		return []any{&k.valText}
	}
}

func (k *kindFixture) Values() []any {
	switch k.inp.Storage() {
	case model.FieldInt:
		return []any{k.valInt}
	case model.FieldBool:
		return []any{k.valBool}
	default:
		return []any{k.valText}
	}
}

// Test_Render verifies String output for inputs with custom rendering.
func Test_Render(t *testing.T) {
	cases := []rc{
		// ── Checkbox ─────────────────────────────────────────────────────────
		{
			t: "Checkbox", name: "renders checkbox input",
			contain: `type='checkbox'`,
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
			contain: `value='1'`,
		},
		{
			t: "Datalist", name: "links input to datalist via list attribute",
			opts:    opts12,
			contain: `list='`,
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
			contain: `value='m'`,
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

		// ── Search ───────────────────────────────────────────────────────────
		{
			t: "Search", name: "renders search input",
			contain: `type='search'`,
		},

		// ── Text ─────────────────────────────────────────────────────────────
		{
			t: "Text", name: "renders input tag",
			contain: `<input`,
		},
		{
			t: "Text", name: "has type text",
			contain: `type='text'`,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.t+"/"+c.name, func(t *testing.T) {
			inp := buildInput(t, c.t)
			fx := &kindFixture{inp: inp}
			f, err := form.New("tid", fx)
			if err != nil {
				t.Fatalf("form.New failed: %v", err)
			}
			if len(c.opts) > 0 {
				f.SetOptions("tfield", c.opts...)
			}
			if len(c.values) > 0 {
				f.SetValues("tfield", c.values...)
			}
			html := f.String()
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
func buildInput(t *testing.T, kind string) input.Input {
	t.Helper()
	var inp input.Input
	switch kind {
	case "Address":
		inp = input.Address()
	case "Checkbox":
		inp = input.Checkbox()
	case "Datalist":
		inp = input.Datalist()
	case "Date":
		inp = input.Date()
	case "Email":
		inp = input.Email()
	case "Filepath":
		inp = input.Filepath()
	case "Gender":
		inp = input.Gender()
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
		inp = input.Radio()
	case "Rut":
		inp = input.Rut()
	case "Search":
		inp = input.Search()
	case "Select":
		inp = input.Select()
	case "Text":
		inp = input.Text()
	case "Textarea":
		inp = input.Textarea()
	default:
		t.Fatalf("unknown input type: %q", kind)
		return nil
	}
	return inp
}

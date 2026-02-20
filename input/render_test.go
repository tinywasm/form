package input

import (
	"testing"

	"github.com/tinywasm/fmt"
)

// Test_Render verifies RenderHTML output for inputs with custom rendering.
// Standard inputs (text, email, etc.) delegate to Base.RenderInput and are not verified here.
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
			html := inp.RenderHTML()
			if !fmt.Contains(html, c.contain) {
				t.Errorf("RenderHTML() missing %q\ngot: %s", c.contain, html)
			}
		})
	}
}

package form

import "github.com/tinywasm/fmt"

// SetSSR enables or disables SSR mode for this form.
func (f *Form) SetSSR(enabled bool) *Form {
	f.ssrMode = enabled
	return f
}

// RenderHTML renders the form based on its SSR mode.
func (f *Form) RenderHTML() string {
	out := fmt.GetConv()

	out.Write(`<form id="`).Write(f.ID()).Write(`"`)

	if f.class != "" {
		out.Write(` class="`).Write(f.class).Write(`"`)
	}

	// SSR mode: render method and action
	if f.ssrMode {
		out.Write(` method="`).Write(f.method).Write(`"`)
		out.Write(` action="`).Write(f.action).Write(`"`)
	}

	out.Write(`>`)

	for _, inp := range f.Inputs {
		out.Write(inp.RenderHTML())
	}

	// SSR mode: render submit button
	if f.ssrMode {
		out.Write(`<button type="submit">Submit</button>`)
	}

	out.Write("</form>")
	return out.String()
}

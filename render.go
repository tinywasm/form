package form

import "github.com/tinywasm/fmt"

// SetSSR enables or disables SSR mode for this form.
func (f *Form) SetSSR(enabled bool) *Form {
	f.ssrMode = enabled
	return f
}

// String serializes the form to its HTML string representation.
func (f *Form) String() string {
	out := fmt.GetConv()

	out.Write(`<form id="`).Write(f.GetID()).Write(`"`)

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
		out.Write(RenderInput(inp))
	}

	// Render submit button in both SSR and WASM modes unless explicitly
	// disabled via HideSubmit. Every form needs a way to submit; the dev can
	// customize the label with SubmitLabel(text).
	if !f.noSubmit {
		label := f.submitLabel
		if label == "" {
			label = "Submit"
		}
		out.Write(`<button type="submit" id="`).Write(f.id).Write(`.submit">`).Write(label).Write(`</button>`)
	}

	out.Write("</form>")
	return out.String()
}

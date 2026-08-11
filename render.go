package form

import (
	"github.com/tinywasm/dom"
	"github.com/tinywasm/widget"
)

// SetSSR enables or disables SSR mode for this form.
func (f *Form) SetSSR(enabled bool) *Form {
	f.ssrMode = enabled
	return f
}

// String serializes the form to its HTML string representation.
func (f *Form) String() string {
	return f.Render().String()
}

// Render returns a reactive dom.Element tree for the form.
func (f *Form) Render() *dom.Element {
	el := dom.NewElement("form").ID(f.GetID())

	if f.class != "" {
		el.Class(f.class)
	}

	// SSR mode: render method and action
	if f.ssrMode {
		el.Attr("method", f.method).Attr("action", f.action)
	}

	for _, child := range f.Children() {
		el.Child(child)
	}

	// Submit button
	if !f.noSubmit {
		btn := dom.NewElement("button").
			Attr("type", "submit").
			Class(widget.NameField.Class(widget.PartSubmit).String()).
			ID(f.id + ".submit")

		btn.BindAttrBool("disabled", f.submitting)

		btn.BindTextFunc(func() string {
			if f.submitting.Get() {
				label := f.submitLoadingLabel
				if label == "" {
					label = f.resolveSubmitLabel() + "..."
				}
				return label
			}
			return f.resolveSubmitLabel()
		})

		// Wrapped in the same tw-field box every field gets, not appended
		// bare: the <form> carries no class, so a field's own inset IS its
		// spacing (fieldset's Root pads by Space2 for exactly that reason).
		// A bare button skips that inset and comes out wider than the inputs
		// above it on both edges — the one misaligned edge in an otherwise
		// squared-off stack. Sharing the wrapper aligns it by construction
		// rather than by a margin tuned to match fieldset's padding.
		el.Child(dom.NewElement("div").
			Class(widget.NameField.Root().String()).
			Child(btn))
	}

	// Bind submit event
	el.On("submit", func(e dom.Event) {
		e.PreventDefault()
		f.Submit()
	})

	return el
}

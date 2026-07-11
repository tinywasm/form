package form

import "github.com/tinywasm/dom"

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

		el.Child(btn)
	}

	// Bind submit event
	el.On("submit", func(e dom.Event) {
		e.PreventDefault()
		f.Submit()
	})

	return el
}

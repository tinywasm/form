//go:build wasm

package form

import (
	"github.com/tinywasm/dom"
)

// OnMount implements dom.Mountable.
// It sets up event delegation for all inputs within the form.
func (f *Form) OnMount() {
	el, ok := dom.Get(f.GetID())
	if !ok {
		return
	}

	// 1. Input/Change listener for live sync
	onInput := func(e dom.Event) {
		id := e.TargetID()

		// Find input by ID
		for _, inp := range f.Inputs {
			if inp.GetID() == id {
				val := e.TargetValue()
				if setter, ok := inp.(interface{ SetValues(...string) }); ok {
					setter.SetValues(val)
				}

				// Real-time validation feedback
				errID := inp.GetID() + ".error"
				if ref, ok := dom.Get(errID); ok {
					if err := inp.Validate(val); err != nil {
						ref.SetText(err.Error())
						ref.SetAttr("class", "tw-field-error tw-field-error--visible")
					} else {
						ref.SetText("")
						ref.SetAttr("class", "tw-field-error")
					}
				}
				break
			}
		}
	}

	// 2. Submit listener
	onSubmit := func(e dom.Event) {
		e.PreventDefault()

		// Sync all values to struct
		f.SyncValues(f.data)

		// Validate all (final check)
		if err := f.Validate(); err != nil {
			// Handle validation error (could be passed back to user)
			return
		}

		// Disable button + show loading state via Reference (no re-render)
		submitID := f.GetID() + ".submit"
		loadingLabel := f.submitLoadingLabel
		if loadingLabel == "" {
			loadingLabel = f.resolveSubmitLabel() + "..."
		}

		if btnRef, ok := dom.Get(submitID); ok {
			btnRef.SetAttr("disabled", "")
			btnRef.SetText(loadingLabel)
		}

		// Call OnSubmit callback with async done helper
		if f.onSubmit != nil {
			f.onSubmit(f.data, func(err error) {
				if err == nil && !f.noResetOnSuccess {
					f.reset()
				}

				// Restore button state
				if btnRef, ok := dom.Get(submitID); ok {
					btnRef.RemoveAttr("disabled")
					btnRef.SetText(f.resolveSubmitLabel())
				}
			})
		}
	}

	el.On("input", onInput)
	el.On("change", onInput)
	el.On("submit", onSubmit)
}

// OnUnmount implements dom.Mountable.
func (f *Form) OnUnmount() {
	// Cleanup is handled automatically by tinywasm/dom on Unmount
}

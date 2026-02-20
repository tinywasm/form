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
				inp.ValidateField(val)
				break
			}
		}
	}

	// 2. Submit listener
	onSubmit := func(e dom.Event) {
		e.PreventDefault()

		// Sync all values to struct
		f.SyncValues()

		// Validate all (final check)
		if err := f.Validate(); err != nil {
			// Handle validation error (could be passed back to user)
			return
		}

		// Call OnSubmit callback
		if f.onSubmit != nil {
			f.onSubmit(f.Value)
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

package form

// Validate validates all inputs and returns the first error found.
func (f *Form) Validate() error {
	for i, inp := range f.Inputs {
		// Skip validation if requested via tag
		if skipper, ok := inp.(interface{ GetSkipValidation() bool }); ok && skipper.GetSkipValidation() {
			continue
		}

		// Signal is the source of truth in WASM mode.
		val := f.valueSignals[i].Get()

		// Fallback only if we are somehow in SSR mode where signals might be empty
		if val == "" && f.ssrMode {
			if valuer, ok := inp.(interface{ GetSelectedValue() string }); ok {
				val = valuer.GetSelectedValue()
			}
		}

		if err := inp.Validate(val); err != nil {
			return err
		}
	}
	return nil
}

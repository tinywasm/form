package form

// Validate validates all inputs and returns the first error found.
func (f *Form) Validate() error {
	for _, inp := range f.Inputs {
		if valuer, ok := inp.(interface{ GetSelectedValue() string }); ok {
			val := valuer.GetSelectedValue()
			if err := inp.ValidateField(val); err != nil {
				return err
			}
		}
	}
	return nil
}

package input

import "github.com/tinywasm/model"

// decimal is a distinct type from number — Number() keeps reporting FieldInt
// unconditionally; nothing about it changes.
type decimal struct{ Base }

// Decimal creates a number input whose storage is float64, for fields that need
// fractional precision (price, measurements, percentages, ...). Renders as the
// same HTML <input type="number"> as Number(); the only difference is Storage().
func Decimal() Input {
	d := &decimal{}
	d.Numbers = true
	d.Extra = []rune{'.', '-'} // allow the decimal point and a leading minus sign
	d.Minimum = 1
	d.Maximum = 20
	d.InitBase("", "", "number")
	return d
}

// Storage satisfies model.Kind.Storage(): a decimal field stores as FieldFloat.
func (d *decimal) Storage() model.FieldType { return model.FieldFloat }

// Clone creates a new decimal input with the given parentID and name.
func (d *decimal) Clone(parentID, name string) Input {
	c := *d
	c.InitBase(parentID, name, "number")
	return &c
}

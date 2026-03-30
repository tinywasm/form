package input

import "github.com/tinywasm/fmt"

// radio represents a standard radio button input.
// NewRadio returns a template instance for use in fmt.Field.Widget (no position).
func NewRadio() fmt.Widget { return Radio("", "") }

type radio struct{ Base }

// Radio creates a new Radio input instance.
func Radio(parentID, name string) Input {
	r := &radio{}
	r.Letters = true
	r.Numbers = true
	r.Minimum = 1
	r.InitBase(parentID, name, "radio")
	return r
}

// Clone creates a new Radio input.
func (r *radio) Clone(parentID, name string) fmt.Widget { return Radio(parentID, name) }

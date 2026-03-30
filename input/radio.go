package input

import "github.com/tinywasm/fmt"

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

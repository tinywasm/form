package input

import "github.com/tinywasm/fmt"

type radio struct{ Base }

// Radio creates a new Radio input instance.
func Radio() Input {
	r := &radio{}
	r.Letters = true
	r.Numbers = true
	r.Minimum = 1
	r.InitBase("", "", "radio")
	return r
}

// Clone creates a new Radio input.
func (r *radio) Clone(parentID, name string) fmt.Widget {
	c := *r
	c.InitBase(parentID, name, "radio")
	return &c
}

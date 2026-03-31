package input

import "github.com/tinywasm/fmt"

type number struct{ Base }

// Number creates a new number input instance.
func Number() Input {
	n := &number{}
	n.Numbers = true
	n.Minimum = 1
	n.Maximum = 20
	n.InitBase("", "", "number")
	return n
}

// Clone creates a new number input with the given parentID and name.
func (n *number) Clone(parentID, name string) fmt.Widget {
	c := *n
	c.InitBase(parentID, name, "number")
	return &c
}

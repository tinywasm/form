package input

import "github.com/tinywasm/fmt"

// number represents a numeric input field.
// NewNumber returns a template instance for use in fmt.Field.Widget (no position).
func NewNumber() fmt.Widget { return Number("", "") }

type number struct{ Base }

// Number creates a new number input instance.
func Number(parentID, name string) Input {
	n := &number{}
	n.Numbers = true
	n.Minimum = 1
	n.Maximum = 20
	n.InitBase(parentID, name, "number", "num", "amount", "price", "age")
	return n
}

// Clone creates a new number input with the given parentID and name.
func (n *number) Clone(parentID, name string) fmt.Widget { return Number(parentID, name) }

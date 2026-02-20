package input

// number represents a numeric input field.
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
func (n *number) Clone(parentID, name string) Input { return Number(parentID, name) }

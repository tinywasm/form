package input

// number represents a numeric input field.
type number struct {
	Base
	Permitted Permitted
}

// Number creates a new number input instance.
func Number(parentID, name string) Input {
	n := &number{
		Permitted: Permitted{
			Numbers: true,
			Minimum: 1,
			Maximum: 20,
		},
	}
	// htmlName: "number", aliases: "num", "amount", "price", "age"
	n.Base.InitBase(parentID, name, "number", "num", "amount", "price", "age")
	return n
}

// HTMLName returns "number".
func (n *number) HTMLName() string {
	return n.Base.HTMLName()
}

// ValidateField validates the value against Permitted rules.
func (n *number) ValidateField(value string) error {
	return n.Permitted.Validate(value)
}

// RenderHTML delegates to Base.RenderInput.
func (n *number) RenderHTML() string {
	return n.Base.RenderInput()
}

// Clone creates a new number input with the given parentID and name.
func (n *number) Clone(parentID, name string) Input {
	return Number(parentID, name)
}

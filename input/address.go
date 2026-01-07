package input

// address represents a standard address input.
type address struct {
	Base
	Permitted Permitted
}

// Address creates a new Address input instance.
func Address(parentID, name string) Input {
	a := &address{
		Permitted: Permitted{
			Letters:     true,
			Numbers:     true,
			WhiteSpaces: true,
			Characters:  []rune{' ', '.', ',', '#', '-', '/', '(', ')'},
			Minimum:     5,
			Maximum:     200,
		},
	}
	// htmlName: "text", aliases: "address", "direccion", "dir", "location"
	a.Base.InitBase(parentID+"."+name, name, "text", "address", "direccion", "dir", "location")

	// Default placeholder
	a.Base.SetPlaceholder("Enter Address")

	return a
}

// HTMLName returns "text".
func (a *address) HTMLName() string {
	return "text"
}

// ValidateField validates the value against Permitted rules.
func (a *address) ValidateField(value string) error {
	return a.Permitted.Validate(value)
}

// RenderHTML delegates to Base.RenderInput.
func (a *address) RenderHTML() string {
	return a.Base.RenderInput()
}

// Clone creates a new Address input with the given parentID and name.
func (a *address) Clone(parentID, name string) Input {
	return Address(parentID, name)
}

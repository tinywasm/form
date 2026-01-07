package input

// text represents a standard text input.
type text struct {
	Base
	Permitted Permitted
}

// Text creates a new Text input instance.
func Text(parentID, name string) Input {
	t := &text{
		Permitted: Permitted{
			Letters:    true,
			Numbers:    true,
			Characters: []rune{' ', '.', ',', '(', ')'},
			Minimum:    2,
			Maximum:    100,
		},
	}
	// htmlName: "text", aliases: "name", "fullname", "username"
	t.Base.InitBase(parentID+"."+name, name, "text", "name", "fullname", "username")
	return t
}

// HTMLName returns "text".
func (t *text) HTMLName() string {
	return t.Base.HTMLName()
}

// ValidateField validates the value against Permitted rules.
func (t *text) ValidateField(value string) error {
	return t.Permitted.Validate(value)
}

// RenderHTML delegates to Base.RenderInput.
func (t *text) RenderHTML() string {
	return t.Base.RenderInput()
}

// Clone creates a new Text input with the given parentID and name.
func (t *text) Clone(parentID, name string) Input {
	return Text(parentID, name)
}

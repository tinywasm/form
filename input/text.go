package input

// text represents a standard text input.
type text struct {
	Base      // Embed generic state and logic
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
	// Initialize Base fields using the method from Base
	t.Base.InitBase(parentID+"."+name, name)
	return t
}

// HtmlName returns "text".
func (t *text) HtmlName() string {
	return "text"
}

// ValidateField validates the value against Permitted rules.
func (t *text) ValidateField(value string) error {
	return t.Permitted.Validate(value)
}

// RenderHTML delegates to Base.RenderInput to reuse logic.
func (t *text) RenderHTML() string {
	return t.Base.RenderInput("text")
}

package input

// textarea represents a textarea field.
type textarea struct {
	Base
	Permitted Permitted
}

// Textarea creates a new textarea input instance.
func Textarea(parentID, name string) Input {
	t := &textarea{
		Permitted: Permitted{
			Letters:    true,
			Numbers:    true,
			Characters: []rune{' ', '.', ',', '-', '_', ':', ';', '(', ')', '\n', '\r', '$', '#', '!', '?'},
			Minimum:    5,
			Maximum:    2000,
		},
	}
	// htmlName: "textarea", aliases: "text", "description", "details"
	t.Base.InitBase(parentID, name, "textarea", "description", "details", "comments")
	return t
}

// HTMLName returns "textarea".
func (t *textarea) HTMLName() string {
	return t.Base.HTMLName()
}

// ValidateField validates the value against Permitted rules.
func (t *textarea) ValidateField(value string) error {
	return t.Permitted.Validate(value)
}

// RenderHTML delegates to Base.RenderInput.
func (t *textarea) RenderHTML() string {
	return t.Base.RenderInput()
}

// Clone creates a new textarea input with the given parentID and name.
func (t *textarea) Clone(parentID, name string) Input {
	return Textarea(parentID, name)
}

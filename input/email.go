package input

// email represents an email input.
type email struct {
	Base
	Permitted Permitted
}

// Email creates a new Email input instance.
func Email(parentID, name string) Input {
	e := &email{
		Permitted: Permitted{
			Letters:    true,
			Numbers:    true,
			Characters: []rune{'@', '.', '_', '-'},
			Minimum:    5,
			Maximum:    100,
		},
	}
	// htmlName: "email", aliases: "mail", "correo"
	e.Base.InitBase(parentID, name, "email", "mail", "correo")
	return e
}

// HTMLName returns "email".
func (e *email) HTMLName() string {
	return e.Base.HTMLName()
}

// ValidateField validates the value against Permitted rules.
func (e *email) ValidateField(value string) error {
	return e.Permitted.Validate(value)
}

// RenderHTML delegates to Base.RenderInput.
func (e *email) RenderHTML() string {
	return e.Base.RenderInput()
}

// Clone creates a new Email input with the given parentID and name.
func (e *email) Clone(parentID, name string) Input {
	return Email(parentID, name)
}

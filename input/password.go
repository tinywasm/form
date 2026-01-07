package input

// password represents a password input.
type password struct {
	Base
	Permitted Permitted
}

// Password creates a new Password input instance.
func Password(parentID, name string) Input {
	p := &password{
		Permitted: Permitted{
			Letters:    true,
			Numbers:    true,
			Tilde:      true,
			Characters: []rune{'!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '-', '_', '=', '+'},
			Minimum:    5,
			Maximum:    50,
		},
	}
	// htmlName: "password", aliases: "pass", "clave", "pwd"
	p.Base.InitBase(parentID+"."+name, name, "password", "pass", "clave", "pwd")
	return p
}

// HtmlName returns "password".
func (p *password) HtmlName() string {
	return p.Base.GetHtmlName()
}

// ValidateField validates the value against Permitted rules.
func (p *password) ValidateField(value string) error {
	return p.Permitted.Validate(value)
}

// RenderHTML delegates to Base.RenderInput.
func (p *password) RenderHTML() string {
	return p.Base.RenderInput()
}

// Clone creates a new Password input with the given parentID and name.
func (p *password) Clone(parentID, name string) Input {
	return Password(parentID, name)
}

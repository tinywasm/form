package input

// phone represents a phone number input field.
type phone struct {
	Base
	Permitted Permitted
}

// Phone creates a new phone input instance.
func Phone(parentID, name string) Input {
	p := &phone{
		Permitted: Permitted{
			Numbers:    true,
			Characters: []rune{'+', ' ', '(', ')', '-'},
			Minimum:    7,
			Maximum:    15,
		},
	}
	// htmlName: "tel", aliases: "phone", "mobile", "cell"
	p.Base.InitBase(parentID, name, "tel", "phone", "mobile", "cell")
	return p
}

// HTMLName returns "tel".
func (p *phone) HTMLName() string {
	return p.Base.HTMLName()
}

// ValidateField validates the value against Permitted rules.
func (p *phone) ValidateField(value string) error {
	return p.Permitted.Validate(value)
}

// RenderHTML delegates to Base.RenderInput.
func (p *phone) RenderHTML() string {
	return p.Base.RenderInput()
}

// Clone creates a new phone input with the given parentID and name.
func (p *phone) Clone(parentID, name string) Input {
	return Phone(parentID, name)
}

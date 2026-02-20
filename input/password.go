package input

// password represents a password input.
type password struct{ Base }

// Password creates a new Password input instance.
func Password(parentID, name string) Input {
	p := &password{}
	p.Letters = true
	p.Numbers = true
	p.Tilde = true
	p.Characters = []rune{'!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '-', '_', '=', '+'}
	p.Minimum = 5
	p.Maximum = 50
	p.InitBase(parentID, name, "password", "pass", "clave", "pwd")
	return p
}

// Clone creates a new Password input with the given parentID and name.
func (p *password) Clone(parentID, name string) Input { return Password(parentID, name) }

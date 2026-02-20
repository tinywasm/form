package input

// email represents an email input.
type email struct{ Base }

// Email creates a new Email input instance.
func Email(parentID, name string) Input {
	e := &email{}
	e.Letters = true
	e.Numbers = true
	e.Characters = []rune{'@', '.', '_', '-'}
	e.Minimum = 5
	e.Maximum = 100
	e.InitBase(parentID, name, "email", "mail", "correo")
	return e
}

// Clone creates a new Email input with the given parentID and name.
func (e *email) Clone(parentID, name string) Input { return Email(parentID, name) }

package input

import "github.com/tinywasm/fmt"

type email struct{ Base }

// Email creates a new Email input instance.
func Email(parentID, name string) Input {
	e := &email{}
	e.Letters = true
	e.Numbers = true
	e.Extra = []rune{'@', '.', '_', '-'}
	e.Minimum = 5
	e.Maximum = 100
	e.InitBase(parentID, name, "email", "mail", "correo")
	return e
}

// Clone creates a new Email input with the given parentID and name.
func (e *email) Clone(parentID, name string) fmt.Widget { return Email(parentID, name) }

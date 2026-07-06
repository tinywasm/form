package input

import "github.com/tinywasm/model"


type email struct{ Base }

// Email creates a new Email input instance.
func Email() Input {
	e := &email{}
	e.Letters = true
	e.Numbers = true
	e.Extra = []rune{'@', '.', '_', '-'}
	e.Minimum = 5
	e.Maximum = 100
	e.InitBase("", "", "email")
	return e
}

// Clone creates a new Email input with the given parentID and name.
func (e *email) Clone(parentID, name string) model.Widget {
	c := *e
	c.InitBase(parentID, name, "email")
	return &c
}

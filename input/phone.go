package input

import "github.com/tinywasm/fmt"

// phone represents a phone number input field.
// NewPhone returns a template instance for use in fmt.Field.Widget (no position).
func NewPhone() fmt.Widget { return Phone("", "") }

type phone struct{ Base }

// Phone creates a new phone input instance.
func Phone(parentID, name string) Input {
	p := &phone{}
	p.Numbers = true
	p.Characters = []rune{'+', ' ', '(', ')', '-'}
	p.Minimum = 7
	p.Maximum = 15
	p.InitBase(parentID, name, "tel", "phone", "mobile", "cell")
	return p
}

// Clone creates a new phone input with the given parentID and name.
func (p *phone) Clone(parentID, name string) fmt.Widget { return Phone(parentID, name) }

package input

import "github.com/tinywasm/fmt"

type password struct{ Base }

// Password creates a new Password input instance.
func Password() Input {
	p := &password{}
	p.Letters = true
	p.Numbers = true
	p.Tilde = true
	p.Extra = []rune{'!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '-', '_', '=', '+'}
	p.Minimum = 5
	p.Maximum = 50
	p.InitBase("", "", "password")
	return p
}

// Clone creates a new Password input with the given parentID and name.
func (p *password) Clone(parentID, name string) fmt.Widget {
	c := *p
	c.InitBase(parentID, name, "password")
	return &c
}

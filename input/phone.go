package input

import "github.com/tinywasm/fmt"

type phone struct{ Base }

// Phone creates a new phone input instance.
func Phone() Input {
	p := &phone{}
	p.Numbers = true
	p.Spaces = true
	p.Extra = []rune{'+', '(', ')', '-'}
	p.Minimum = 7
	p.Maximum = 15
	p.InitBase("", "", "tel")
	return p
}

// Clone creates a new phone input with the given parentID and name.
func (p *phone) Clone(parentID, name string) fmt.Widget {
	c := *p
	c.InitBase(parentID, name, "tel")
	return &c
}

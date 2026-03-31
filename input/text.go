package input

import "github.com/tinywasm/fmt"

type text struct{ Base }

// Text creates a new Text input instance.
func Text() Input {
	t := &text{}
	t.Letters = true
	t.Tilde = true
	t.Numbers = true
	t.Spaces = true
	t.Extra = []rune{'.', ',', '(', ')'}
	t.Minimum = 2
	t.Maximum = 100
	t.InitBase("", "", "text")
	return t
}

// Clone creates a new Text input with the given parentID and name.
func (t *text) Clone(parentID, name string) fmt.Widget {
	c := *t
	c.InitBase(parentID, name, "text")
	return &c
}

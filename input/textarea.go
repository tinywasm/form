package input

import "github.com/tinywasm/fmt"

type textarea struct{ Base }

// Textarea creates a new textarea input instance.
func Textarea() Input {
	t := &textarea{}
	t.Letters = true
	t.Numbers = true
	t.Tilde = true
	t.Spaces = true
	t.BreakLine = true
	t.Extra = []rune{'.', ',', '-', '_', ':', ';', '(', ')', '$', '#', '!', '?'}
	t.Minimum = 5
	t.Maximum = 2000
	t.InitBase("", "", "textarea")
	return t
}

// Clone creates a new textarea input with the given parentID and name.
func (t *textarea) Clone(parentID, name string) fmt.Widget {
	c := *t
	c.InitBase(parentID, name, "textarea")
	return &c
}

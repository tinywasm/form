package input

import "github.com/tinywasm/fmt"

type textarea struct{ Base }

// Textarea creates a new textarea input instance.
func Textarea(parentID, name string) Input {
	t := &textarea{}
	t.Letters = true
	t.Numbers = true
	t.Tilde = true
	t.Spaces = true
	t.BreakLine = true
	t.Extra = []rune{'.', ',', '-', '_', ':', ';', '(', ')', '$', '#', '!', '?'}
	t.Minimum = 5
	t.Maximum = 2000
	t.InitBase(parentID, name, "textarea", "description", "details", "comments")
	return t
}

// Clone creates a new textarea input with the given parentID and name.
func (t *textarea) Clone(parentID, name string) fmt.Widget { return Textarea(parentID, name) }

package input

// textarea represents a textarea field.
type textarea struct{ Base }

// Textarea creates a new textarea input instance.
func Textarea(parentID, name string) Input {
	t := &textarea{}
	t.Letters = true
	t.Numbers = true
	t.Tilde = true
	t.WhiteSpaces = true
	t.BreakLine = true
	t.Characters = []rune{'.', ',', '-', '_', ':', ';', '(', ')', '$', '#', '!', '?'}
	t.Minimum = 5
	t.Maximum = 2000
	t.InitBase(parentID, name, "textarea", "description", "details", "comments")
	return t
}

// Clone creates a new textarea input with the given parentID and name.
func (t *textarea) Clone(parentID, name string) Input { return Textarea(parentID, name) }

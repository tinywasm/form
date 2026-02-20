package input

// text represents a standard text input.
type text struct{ Base }

// Text creates a new Text input instance.
func Text(parentID, name string) Input {
	t := &text{}
	t.Letters = true
	t.Tilde = true
	t.Numbers = true
	t.Characters = []rune{' ', '.', ',', '(', ')'}
	t.Minimum = 2
	t.Maximum = 100
	t.InitBase(parentID, name, "text", "name", "fullname", "username")
	return t
}

// Clone creates a new Text input with the given parentID and name.
func (t *text) Clone(parentID, name string) Input { return Text(parentID, name) }

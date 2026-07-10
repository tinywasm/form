package input







type textarea struct{ Base }

// Textarea creates a new textarea input instance.
//ormc:storage text
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
func (t *textarea) Clone(parentID, name string) Input {
	c := *t
	c.InitBase(parentID, name, "textarea")
	return &c
}

func (t *textarea) setTilde(v bool) { t.Tilde = v }

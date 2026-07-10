package input







type address struct{ Base }

// Address creates a new Address input instance.
func Address() Input {
	a := &address{}
	a.Letters = true
	a.Numbers = true
	a.Spaces = true
	a.Extra = []rune{'.', ',', '#', '-', '/', '(', ')'}
	a.Minimum = 5
	a.Maximum = 200
	a.InitBase("", "", "text")
	a.SetPlaceholder("Enter Address")
	return a
}

// Clone creates a new Address input with the given parentID and name.
func (a *address) Clone(parentID, name string) Input {
	c := *a
	c.InitBase(parentID, name, "text")
	return &c
}

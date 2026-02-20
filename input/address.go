package input

// Address creates a new Address input instance (semantic wrapper around text).
func Address(parentID, name string) Input {
	a := &text{}
	a.Letters = true
	a.Numbers = true
	a.WhiteSpaces = true
	a.Characters = []rune{' ', '.', ',', '#', '-', '/', '(', ')'}
	a.Minimum = 5
	a.Maximum = 200
	a.InitBase(parentID, name, "text", "address", "addr", "direccion", "dir", "location")
	a.SetPlaceholder("Enter Address")
	return a
}

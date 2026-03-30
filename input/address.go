package input

// address represents an address input (semantic wrapper around text).
type address struct{ Base }

// Address creates a new Address input instance.
func Address(parentID, name string) Input {
	a := &address{}
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

// Clone creates a new Address input with the given parentID and name.
func (a *address) Build(parentID, name string) Input { return Address(parentID, name) }

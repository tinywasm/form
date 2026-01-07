package input

// Address creates a new Address input instance.
// It is a semantic wrapper around Text.
func Address(parentID, name string) Input {
	a := Text(parentID, name).(*text)

	// Add specific aliases for address
	a.Base.aliases = append(a.Base.aliases, "address", "direccion", "dir", "location")

	// Pre-configure validation and placeholder
	a.Base.SetPlaceholder("Enter Address")
	a.Permitted.Minimum = 5
	a.Permitted.Maximum = 200
	a.Permitted.WhiteSpaces = true
	a.Permitted.Characters = []rune{' ', '.', ',', '#', '-', '/', '(', ')'}

	return a
}

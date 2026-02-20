package input

// select_ represents a dropdown selection.
type select_ struct{ Base }

// Select creates a new Select input instance.
func Select(parentID, name string) Input {
	s := &select_{}
	s.Letters = true
	s.Numbers = true
	s.Minimum = 1
	s.InitBase(parentID, name, "select", "role", "tipo")
	return s
}

// Clone creates a new Select input.
func (s *select_) Clone(parentID, name string) Input { return Select(parentID, name) }

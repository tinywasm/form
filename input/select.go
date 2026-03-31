package input

import "github.com/tinywasm/fmt"

type select_ struct{ Base }

// Select creates a new Select input instance.
func Select() Input {
	s := &select_{}
	s.Letters = true
	s.Numbers = true
	s.Minimum = 1
	s.InitBase("", "", "select")
	return s
}

// Clone creates a new Select input.
func (s *select_) Clone(parentID, name string) fmt.Widget {
	c := *s
	c.InitBase(parentID, name, "select")
	return &c
}

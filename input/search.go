package input

import "github.com/tinywasm/fmt"

type search_ struct{ Base }

// Search creates a new Search input instance.
func Search() Input {
	s := &search_{}
	s.Letters = true
	s.Numbers = true
	s.Spaces = true
	s.Minimum = 0
	s.Maximum = 100
	s.InitBase("", "", "search")
	return s
}

// Clone creates a new Search input with the given parentID and name.
func (s *search_) Clone(parentID, name string) fmt.Widget {
	c := *s
	c.InitBase(parentID, name, "search")
	return &c
}

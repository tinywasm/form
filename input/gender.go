package input

import "github.com/tinywasm/fmt"

type gender struct{ Base }

// Gender creates a new Gender input instance with default Male/Female options.
func Gender() Input {
	g := &gender{}
	g.Letters = true
	g.Numbers = true
	g.Minimum = 1
	g.InitBase("", "", "radio")
	g.SetOptions(
		fmt.KeyValue{Key: "m", Value: "Male"},
		fmt.KeyValue{Key: "f", Value: "Female"},
	)
	return g
}

// Clone creates a new Gender input with the given parentID and name.
func (g *gender) Clone(parentID, name string) fmt.Widget {
	c := *g
	c.InitBase(parentID, name, "radio")
	return &c
}

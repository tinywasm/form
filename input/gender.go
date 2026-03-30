package input

import "github.com/tinywasm/fmt"

// gender represents a gender input (semantic wrapper around radio).
type gender struct{ Base }

// Gender creates a new Gender input instance with default Male/Female options.
func Gender(parentID, name string) Input {
	g := &gender{}
	g.Letters = true
	g.Numbers = true
	g.Minimum = 1
	g.InitBase(parentID, name, "radio", "gender", "sexo")
	g.SetOptions(
		fmt.KeyValue{Key: "m", Value: "Male"},
		fmt.KeyValue{Key: "f", Value: "Female"},
	)
	return g
}

// Clone creates a new Gender input with the given parentID and name.
func (g *gender) Build(parentID, name string) Input { return Gender(parentID, name) }

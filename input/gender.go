package input

import "github.com/tinywasm/fmt"

// gender represents a gender input (semantic wrapper around radio).
// NewGender returns a template instance for use in fmt.Field.Widget (no position).
func NewGender() fmt.Widget { return Gender("", "") }

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
func (g *gender) Clone(parentID, name string) fmt.Widget { return Gender(parentID, name) }

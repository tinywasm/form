package input

import "github.com/tinywasm/fmt"

// Gender creates a new Gender input instance with default Male/Female options.
// It is a semantic wrapper around Radio.
func Gender(parentID, name string) Input {
	g := Radio(parentID, name).(*radio)

	// Add specific aliases for gender
	g.Base.aliases = append(g.Base.aliases, "gender", "sexo")

	// Default options
	g.Base.SetOptions(
		fmt.KeyValue{Key: "m", Value: "Male"},
		fmt.KeyValue{Key: "f", Value: "Female"},
	)

	return g
}

package form

import (
	"github.com/tinywasm/fmt"
)

// ParseOptionsTag parses "key1:text1,key2:text2" format into []fmt.KeyValue.
func ParseOptionsTag(tag string) []fmt.KeyValue {
	// Reusing the same logic as fmt.TagPairs but on the direct string
	return fmt.Convert(`tmp:"` + tag + `"`).TagPairs("tmp")
}

// GetTagOptions extracts options from a struct field tag using fmt.TagPairs.
func GetTagOptions(fieldTag string) ([]fmt.KeyValue, bool) {
	opts := fmt.Convert(fieldTag).TagPairs("options")
	if opts == nil {
		return nil, false
	}
	return opts, true
}

package input

import (
	"github.com/tinywasm/fmt"
)

// ip represents an IP address input field (IPv4 or IPv6).
type ip struct {
	Base
	Permitted Permitted
}

// IP creates a new IP input instance.
func IP(parentID, name string) Input {
	i := &ip{
		Permitted: Permitted{
			Numbers:    true,
			Letters:    true, // hex for ipv6
			Characters: []rune{'.', ':'},
			Minimum:    7,  // 1.1.1.1
			Maximum:    39, // full ipv6 length
		},
	}
	// htmlName: "text", aliases: "ip", "address"
	i.Base.InitBase(parentID, name, "text", "ip", "address")
	return i
}

// HTMLName returns "text".
func (i *ip) HTMLName() string {
	return i.Base.HTMLName()
}

// ValidateField validates the value format for ipv4 or ipv6.
func (i *ip) ValidateField(value string) error {
	if value == "0.0.0.0" {
		return fmt.Err("Format", "Invalid")
	}

	err := i.Permitted.Validate(value)
	if err != nil {
		return err
	}

	// rudimentary IP format checking (would use net.ParseIP if we had full stdlib, but we use fmt/Permitted)
	hasDot := false
	hasColon := false
	for _, char := range value {
		if char == '.' {
			hasDot = true
		}
		if char == ':' {
			hasColon = true
		}
	}
	if hasDot && hasColon {
		return fmt.Err("Format", "Invalid")
	}

	return nil
}

// RenderHTML delegates to Base.RenderInput.
func (i *ip) RenderHTML() string {
	return i.Base.RenderInput()
}

// Clone creates a new IP input with the given parentID and name.
func (i *ip) Clone(parentID, name string) Input {
	return IP(parentID, name)
}

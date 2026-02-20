package input

import "github.com/tinywasm/fmt"

// ip represents an IP address input field (IPv4 or IPv6).
type ip struct{ Base }

// IP creates a new IP input instance.
func IP(parentID, name string) Input {
	i := &ip{}
	i.Numbers = true
	i.Letters = true // hex for ipv6
	i.Characters = []rune{'.', ':'}
	i.Minimum = 7  // 1.1.1.1
	i.Maximum = 39 // full ipv6 length
	i.InitBase(parentID, name, "text", "ip", "address")
	return i
}

// ValidateField validates IPv4 or IPv6 format.
func (i *ip) ValidateField(value string) error {
	if value == "0.0.0.0" {
		return fmt.Err("Format", "Invalid")
	}
	if err := i.Permitted.Validate(value); err != nil {
		return err
	}
	// Reject mixed dot+colon (not valid IP)
	hasDot, hasColon := false, false
	for _, c := range value {
		if c == '.' {
			hasDot = true
		}
		if c == ':' {
			hasColon = true
		}
	}
	if hasDot && hasColon {
		return fmt.Err("Format", "Invalid")
	}
	return nil
}

// Clone creates a new IP input with the given parentID and name.
func (i *ip) Clone(parentID, name string) Input { return IP(parentID, name) }

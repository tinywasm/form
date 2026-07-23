package input

import "github.com/tinywasm/fmt"

type ip struct{ Base }

// IP creates a new IP input instance.
func IP() Input {
	i := &ip{}
	i.Numbers = true
	i.Letters = true // hex for ipv6
	i.Extra = []rune{'.', ':'}
	i.Minimum = 7  // 1.1.1.1
	i.Maximum = 39 // full ipv6 length
	i.InitBase("", "", "text")
	return i
}

// Validate validates IPv4 or IPv6 format.
func (i *ip) Validate(value string) error {
	if value == "0.0.0.0" {
		return fmt.Err("Format", "Invalid")
	}

	origMin := i.Minimum
	if value != "" {
		hasColon := false
		for _, c := range value {
			if c == ':' {
				hasColon = true
				break
			}
		}
		if hasColon {
			i.Minimum = 2
		}
	}

	err := i.Permitted.Validate(i.name, value)
	i.Minimum = origMin
	if err != nil {
		return err
	}

	hasDot, hasColon := false, false
	for _, c := range value {
		if c == '.' {
			hasDot = true
		}
		if c == ':' {
			hasColon = true
		}
	}

	switch {
	case hasDot && hasColon:
		return fmt.Err("Format", "Invalid")
	case !hasDot && !hasColon:
		return fmt.Err("Format", "Invalid")
	case hasDot:
		return validateIPv4(value)
	default:
		return validateIPv6(value)
	}
}

func validateIPv4(value string) error {
	parts := fmt.Convert(value).Split(".")
	if len(parts) != 4 {
		return fmt.Err("Format", "Invalid")
	}
	for _, part := range parts {
		if len(part) < 1 || len(part) > 3 {
			return fmt.Err("Format", "Invalid")
		}
		for _, c := range part {
			if c < '0' || c > '9' {
				return fmt.Err("Format", "Invalid")
			}
		}
		val, err := fmt.Convert(part).Int()
		if err != nil || val < 0 || val > 255 {
			return fmt.Err("Format", "Invalid")
		}
	}
	return nil
}

func validateIPv6(value string) error {
	// Count occurrences of "::"
	doubleColonCount := 0
	for i := 0; i < len(value)-1; i++ {
		if value[i] == ':' && value[i+1] == ':' {
			doubleColonCount++
			i++
		}
	}

	if doubleColonCount > 1 {
		return fmt.Err("Format", "Invalid")
	}

	parts := fmt.Convert(value).Split(":")

	startsWithDoubleColon := len(value) >= 2 && value[0] == ':' && value[1] == ':'
	endsWithDoubleColon := len(value) >= 2 && value[len(value)-2] == ':' && value[len(value)-1] == ':'

	if doubleColonCount == 1 {
		if startsWithDoubleColon && endsWithDoubleColon {
			// e.g. "::" -> parts is ["", "", ""]
			if len(parts) != 3 {
				return fmt.Err("Format", "Invalid")
			}
			for _, part := range parts {
				if part != "" {
					return fmt.Err("Format", "Invalid")
				}
			}
		} else if startsWithDoubleColon {
			// e.g. "::1" -> parts is ["", "", "1"]
			if len(parts) < 3 {
				return fmt.Err("Format", "Invalid")
			}
			if parts[0] != "" || parts[1] != "" {
				return fmt.Err("Format", "Invalid")
			}
			nonEmptyCount := 0
			for i := 2; i < len(parts); i++ {
				if parts[i] == "" {
					return fmt.Err("Format", "Invalid")
				}
				if !isValidHexPart(parts[i]) {
					return fmt.Err("Format", "Invalid")
				}
				nonEmptyCount++
			}
			if nonEmptyCount > 7 {
				return fmt.Err("Format", "Invalid")
			}
		} else if endsWithDoubleColon {
			// e.g. "1::" -> parts is ["1", "", ""]
			if len(parts) < 3 {
				return fmt.Err("Format", "Invalid")
			}
			if parts[len(parts)-1] != "" || parts[len(parts)-2] != "" {
				return fmt.Err("Format", "Invalid")
			}
			nonEmptyCount := 0
			for i := 0; i < len(parts)-2; i++ {
				if parts[i] == "" {
					return fmt.Err("Format", "Invalid")
				}
				if !isValidHexPart(parts[i]) {
					return fmt.Err("Format", "Invalid")
				}
				nonEmptyCount++
			}
			if nonEmptyCount > 7 {
				return fmt.Err("Format", "Invalid")
			}
		} else {
			// e.g. "2001::1" -> parts has exactly one empty part
			emptyCount := 0
			nonEmptyCount := 0
			for _, part := range parts {
				if part == "" {
					emptyCount++
				} else {
					if !isValidHexPart(part) {
						return fmt.Err("Format", "Invalid")
					}
					nonEmptyCount++
				}
			}
			if emptyCount != 1 || nonEmptyCount > 7 {
				return fmt.Err("Format", "Invalid")
			}
		}
	} else {
		// doubleColonCount == 0
		if len(parts) != 8 {
			return fmt.Err("Format", "Invalid")
		}
		for _, part := range parts {
			if part == "" || !isValidHexPart(part) {
				return fmt.Err("Format", "Invalid")
			}
		}
	}

	return nil
}

func isValidHexPart(part string) bool {
	if len(part) < 1 || len(part) > 4 {
		return false
	}
	for _, c := range part {
		isDigit := c >= '0' && c <= '9'
		isLowerHex := c >= 'a' && c <= 'f'
		isUpperHex := c >= 'A' && c <= 'F'
		if !isDigit && !isLowerHex && !isUpperHex {
			return false
		}
	}
	return true
}

// Clone satisfies input.Input — IP() returns Input which implements it.
func (i *ip) Clone(parentID, name string) Input {
	c := *i
	c.InitBase(parentID, name, "text")
	return &c
}

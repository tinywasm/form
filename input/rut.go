package input

import (
	"github.com/tinywasm/fmt"
)

// rut represents a Chilean RUT input field.
type rut struct {
	Base
	Permitted Permitted
}

// Rut creates a new RUT input instance.
func Rut(parentID, name string) Input {
	r := &rut{
		Permitted: Permitted{
			Numbers:    true,
			Characters: []rune{'-', 'k', 'K'},
			Minimum:    3,
			Maximum:    12,
		},
	}
	// htmlName: "text", aliases: "rut", "run", "dni"
	r.Base.InitBase(parentID, name, "text", "rut", "run", "dni")
	r.Base.SetPlaceholder("12345678-9")
	return r
}

// HTMLName returns "text".
func (r *rut) HTMLName() string {
	return r.Base.HTMLName()
}

func (r *rut) ValidateField(value string) error {
	err := r.Permitted.Validate(value)
	if err != nil {
		return err
	}

	if len(value) < 3 {
		return fmt.Err("Format", "Invalid")
	}

	hasHyphen := false
	for _, c := range value {
		if c == '-' {
			hasHyphen = true
			break
		}
	}

	if !hasHyphen {
		return fmt.Err("Hyphen", "Missing")
	}

	partsStr := fmt.Convert(value).Split("-")
	if len(partsStr) != 2 {
		return fmt.Err("Format", "Invalid")
	}

	if len(partsStr[0]) == 0 || partsStr[0][0] == '0' {
		return fmt.Err("Format", "Invalid")
	}

	numRut, cvtErr := fmt.Convert(partsStr[0]).Int()
	if cvtErr != nil {
		return fmt.Err("Format", "Invalid")
	}

	dvStr := fmt.Convert(partsStr[1]).ToLower().String()
	dvCalc := r.dvRut(numRut)

	if dvCalc != dvStr {
		return fmt.Err("Digit", "Invalid")
	}

	return nil
}

// dvRut calculates the verification digit for a given RUT number.
func (r *rut) dvRut(rut int) string {
	sum := 0
	factor := 2
	for ; rut != 0; rut /= 10 {
		sum += rut % 10 * factor
		if factor == 7 {
			factor = 2
		} else {
			factor++
		}
	}

	val := 11 - (sum % 11)
	if val == 11 {
		return "0"
	} else if val == 10 {
		return "k"
	}
	return fmt.Convert(val).String()
}

// RenderHTML delegates to Base.RenderInput.
func (r *rut) RenderHTML() string {
	return r.Base.RenderInput()
}

// Clone creates a new rut input with the given parentID and name.
func (r *rut) Clone(parentID, name string) Input {
	return Rut(parentID, name)
}

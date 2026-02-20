package input

import (
	"github.com/tinywasm/fmt"
)

// date represents a date input field.
type date struct {
	Base
	Permitted Permitted
}

// Date creates a new date input instance.
func Date(parentID, name string) Input {
	d := &date{
		Permitted: Permitted{
			Numbers:    true,
			Characters: []rune{'-'},
			Minimum:    10,
			Maximum:    10,
		},
	}
	// htmlName: "date", aliases: "fecha"
	d.Base.InitBase(parentID, name, "date", "fecha")
	return d
}

// HTMLName returns "date".
func (d *date) HTMLName() string {
	return d.Base.HTMLName()
}

// ValidateField validates the value format for YYYY-MM-DD.
func (d *date) ValidateField(value string) error {
	err := d.Permitted.Validate(value)
	if err != nil {
		return err
	}

	if len(value) != 10 {
		return fmt.Err("Format", "Invalid", "2006-01-02")
	}

	// Format check
	for i, char := range value {
		if i == 4 || i == 7 {
			if char != '-' {
				return fmt.Err("Format", "Invalid", "2006-01-02")
			}
		} else {
			if char < '0' || char > '9' {
				return fmt.Err("Format", "Invalid", "2006-01-02")
			}
		}
	}

	year, _ := fmt.Convert(value[:4]).Int()
	month, _ := fmt.Convert(value[5:7]).Int()
	day, _ := fmt.Convert(value[8:10]).Int()

	if year < 1000 || year > 9999 {
		return fmt.Err("Date", "Invalid", "year")
	}
	if month < 1 || month > 12 {
		return fmt.Err("Month", "Invalid")
	}

	if day < 1 {
		return fmt.Err("Date", "Invalid", "day")
	}

	monthDays := d.monthDays(year)[month]
	if day > monthDays {
		return fmt.Err("Date", "Invalid")
	}

	return nil
}

// RenderHTML delegates to Base.RenderInput.
func (d *date) RenderHTML() string {
	return d.Base.RenderInput()
}

// Clone creates a new date input with the given parentID and name.
func (d *date) Clone(parentID, name string) Input {
	return Date(parentID, name)
}

func (d *date) monthDays(year int) map[int]int {
	febDays := 28
	if d.isLeap(year) {
		febDays = 29
	}

	return map[int]int{
		1:  31,
		2:  febDays,
		3:  31,
		4:  30,
		5:  31,
		6:  30,
		7:  31,
		8:  31,
		9:  30,
		10: 31,
		11: 30,
		12: 31,
	}
}

func (d *date) isLeap(year int) bool {
	return year%4 == 0 && year%100 != 0 || year%400 == 0
}

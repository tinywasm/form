package input

import "github.com/tinywasm/fmt"

// date represents a date input field.
type date struct{ Base }

// Date creates a new date input instance.
func Date(parentID, name string) Input {
	d := &date{}
	d.Numbers = true
	d.Characters = []rune{'-'}
	d.Minimum = 10
	d.Maximum = 10
	d.InitBase(parentID, name, "date", "fecha")
	return d
}

// ValidateField validates YYYY-MM-DD format with leap year and day range checks.
func (d *date) ValidateField(value string) error {
	if err := d.Permitted.Validate(value); err != nil {
		return err
	}
	if len(value) != 10 {
		return fmt.Err("Format", "Invalid", "2006-01-02")
	}
	for i, char := range value {
		if i == 4 || i == 7 {
			if char != '-' {
				return fmt.Err("Format", "Invalid", "2006-01-02")
			}
		} else if char < '0' || char > '9' {
			return fmt.Err("Format", "Invalid", "2006-01-02")
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
	if day < 1 || day > d.monthDays(year)[month] {
		return fmt.Err("Date", "Invalid")
	}
	return nil
}

func (d *date) monthDays(year int) [13]int {
	feb := 28
	if year%4 == 0 && year%100 != 0 || year%400 == 0 {
		feb = 29
	}
	return [13]int{0, 31, feb, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
}

// Clone creates a new date input with the given parentID and name.
func (d *date) Clone(parentID, name string) Input { return Date(parentID, name) }

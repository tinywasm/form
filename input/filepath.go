package input

import "github.com/tinywasm/fmt"

type filepath struct{ Base }

// Filepath creates a new filepath input instance.
func Filepath(parentID, name string) Input {
	fp := &filepath{}
	fp.Letters = true
	fp.Numbers = true
	fp.Extra = []rune{'.', '\\', '/', '-', '_'}
	fp.Minimum = 1
	fp.Maximum = 200
	fp.InitBase(parentID, name, "text", "path", "dir", "file")
	return fp
}

// Validate validates the path — no whitespace, no leading backslash.
func (fp *filepath) Validate(value string) error {
	if err := fp.Permitted.Validate(fp.name, value); err != nil {
		return err
	}
	if fmt.Contains(value, " ") {
		return fmt.Err("WhiteSpace", "NotAllowed")
	}
	if len(value) > 0 && value[0] == '\\' {
		return fmt.Err("DoNotStartWith", "\\")
	}
	return nil
}

// Clone creates a new filepath input with the given parentID and name.
func (fp *filepath) Clone(parentID, name string) fmt.Widget { return Filepath(parentID, name) }

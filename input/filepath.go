package input

import "github.com/tinywasm/fmt"

// filepath represents a file or directory path input field.
type filepath struct {
	Base
	Permitted Permitted
}

// Filepath creates a new filepath input instance.
func Filepath(parentID, name string) Input {
	fp := &filepath{
		Permitted: Permitted{
			Letters:    true,
			Numbers:    true,
			Characters: []rune{'.', '\\', '/', '-', '_'},
			Minimum:    1,
			Maximum:    200,
		},
	}
	// htmlName: "text", aliases: "path", "dir", "file"
	fp.Base.InitBase(parentID, name, "text", "path", "dir", "file")
	return fp
}

// HTMLName returns "text".
func (fp *filepath) HTMLName() string {
	return fp.Base.HTMLName()
}

// ValidateField validates the path.
func (fp *filepath) ValidateField(value string) error {
	err := fp.Permitted.Validate(value)
	if err != nil {
		return err
	}

	if fmt.Contains(value, " ") {
		return fmt.Err("WhiteSpace", "NotAllowed")
	}

	// Must not start with \
	if len(value) > 0 && value[0] == '\\' {
		return fmt.Err("DoNotStartWith", "\\")
	}

	return nil
}

// RenderHTML delegates to Base.RenderInput.
func (fp *filepath) RenderHTML() string {
	return fp.Base.RenderInput()
}

// Clone creates a new filepath input with the given parentID and name.
func (fp *filepath) Clone(parentID, name string) Input {
	return Filepath(parentID, name)
}

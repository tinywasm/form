package input

import (
	"github.com/tinywasm/fmt"
)

type Permitted struct {
	Letters         bool
	Tilde           bool
	Numbers         bool
	BreakLine       bool     // line breaks allowed
	WhiteSpaces     bool     // white spaces allowed
	Tabulation      bool     // tabulation allowed
	TextNotAllowed  []string // text not allowed eg: "hola" not allowed
	Characters      []rune   // other special characters eg: '\','/','@'
	Minimum         int      // min characters eg 2 "lo" ok default 0 no defined
	Maximum         int      // max characters eg 1 "l" ok default 0 no defined}
	ExtraValidation func(string) error
	StartWith       *Permitted // characters allowed at the beginning
}

const tabulation = '	'
const white_space = ' '
const break_line = '\n'

var valid_letters = map[rune]bool{
	'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true, 'g': true, 'h': true, 'i': true,
	'j': true, 'k': true, 'l': true, 'm': true, 'n': true, 'o': true, 'p': true, 'q': true, 'r': true,
	's': true, 't': true, 'u': true, 'v': true, 'w': true, 'x': true, 'y': true, 'z': true,
	'ñ': true,

	'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'F': true, 'G': true, 'H': true, 'I': true,
	'J': true, 'K': true, 'L': true, 'M': true, 'N': true, 'O': true, 'P': true, 'Q': true, 'R': true,
	'S': true, 'T': true, 'U': true, 'V': true, 'W': true, 'X': true, 'Y': true, 'Z': true,
	'Ñ': true,
}

var valid_tilde = map[rune]bool{
	'á': true, 'é': true, 'í': true, 'ó': true, 'ú': true,
}

var valid_number = map[rune]bool{
	'0': true, '1': true, '2': true, '3': true, '4': true, '5': true, '6': true, '7': true, '8': true, '9': true,
}

func (h Permitted) Validate(text string) (err error) {

	if h.Minimum != 0 {
		if len(text) < h.Minimum {
			return fmt.Err("MinSize", h.Minimum, fmt.D.Chars)
		}
	}

	if h.Maximum != 0 {
		if len(text) > h.Maximum {
			return fmt.Err(fmt.D.Maximum, h.Maximum, fmt.D.Chars)
		}
	}

	if len(h.TextNotAllowed) != 0 {
		for _, notAllowed := range h.TextNotAllowed {
			if fmt.Contains(text, notAllowed) {
				return fmt.Err(fmt.D.Not, fmt.D.Allowed, ':', h.TextNotAllowed)
			}
		}
	}

	for _, char := range text {
		if char == tabulation && h.Tabulation {
			continue
		}

		if char == white_space && h.WhiteSpaces {
			continue
		}

		if char == break_line && h.BreakLine {
			continue
		}

		if h.Letters {
			if !valid_letters[char] {
				return fmt.Err(string(char), fmt.D.Not, fmt.D.Letters)
			} else {
				err = nil
				continue
			}
		}

		if h.Tilde {
			if !valid_tilde[char] {
				return fmt.Err(string(char), "TildeNotAllowed")
			} else {
				err = nil
				continue
			}
		}

		if h.Numbers {
			if !valid_number[char] {
				if char == ' ' {
					return fmt.Err(fmt.D.Space, fmt.D.Not, fmt.D.Allowed)
				} else {
					return fmt.Err(string(char), fmt.D.Not, fmt.D.Number)
				}
			} else {
				err = nil
				continue
			}
		}

		if len(h.Characters) != 0 {
			var found bool
			for _, c := range h.Characters {
				if c == char {
					found = true
					break
				}
			}

			if found {
				err = nil
				continue
			} else {
				if char == white_space {
					return fmt.Err(fmt.D.Space, fmt.D.Not, fmt.D.Allowed)
				} else if valid_tilde[char] {
					return fmt.Err(string(char), "TildeNotAllowed")
				} else if char == tabulation {
					return fmt.Err(fmt.D.Tab, fmt.D.Text, fmt.D.Not, fmt.D.Allowed)
				} else if char == break_line {
					return fmt.Err("Newline", fmt.D.Not, fmt.D.Allowed)
				}
				return fmt.Err(fmt.D.Character, string(char), fmt.D.Not, fmt.D.Allowed)
			}
		}

		if err != nil {
			return err
		}
	}

	return err
}

func (p Permitted) MinMaxAllowedChars() (min, max int) {
	return p.Minimum, p.Maximum
}

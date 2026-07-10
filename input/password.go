package input







type password struct{ Base }

// Password creates a new Password input instance.
//ormc:storage text
func Password() Input {
	p := &password{}
	p.Letters = true
	p.Numbers = true
	p.Tilde = true
	p.Extra = []rune{'!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '-', '_', '=', '+'}
	p.Minimum = 5
	p.Maximum = 50
	p.InitBase("", "", "password")
	return p
}

// Clone creates a new Password input with the given parentID and name.
func (p *password) Clone(parentID, name string) Input {
	c := *p
	c.InitBase(parentID, name, "password")
	return &c
}

func (p *password) setTilde(v bool) { p.Tilde = v }

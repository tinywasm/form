package form

// Global storage for forms
var forms = make([]*Form, 0)

// Global class configuration
var globalClass string

func SetGlobalClass(classes ...string) {
	for _, c := range classes {
		if globalClass != "" {
			globalClass += " "
		}
		globalClass += c
	}
}

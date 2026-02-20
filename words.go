package form

import "github.com/tinywasm/fmt"

// init registers generic form UI words into the shared fmt dictionary.
// These words are used by form rendering (e.g. submit button) and
// are available to all packages that import tinywasm/form.
func init() {
	fmt.RegisterWords([]fmt.DictEntry{
		{EN: "Submit",   ES: "Enviar",   FR: "Soumettre", DE: "Absenden",    ZH: "提交",  HI: "जमा करें",  AR: "إرسال",    PT: "Enviar",    RU: "Отправить"},
		{EN: "Optional", ES: "Opcional", FR: "Optionnel", DE: "Optional",    ZH: "可选",  HI: "वैकल्पिक", AR: "اختياري",  PT: "Opcional",  RU: "Необязательно"},
	})
}

# Struct Tags

Tags configure inputs declaratively on struct fields.

## Supported Tags

| Tag | Format | Implemented |
|-----|--------|-------------|
| `options` | `"key1:text1,key2:text2"` | ✅ |
| `placeholder` | `"string"` | ✅ (optional — see Auto-translation) |
| `title` | `"string"` | ✅ (optional — see Auto-translation) |
| `validate` | `"false"` | ✅ (skips validation) |
| `label` | `"string"` | 🔲 reserved |

## Auto-translation (No Tag Required)

When `placeholder` or `title` tags are omitted, `form.New()` auto-generates
the placeholder/title by calling `fmt.Translate(fieldName)`:

- If the field name is registered in the `tinywasm/fmt` dictionary → **translated value** (e.g. `"Email"` → `"Correo electrónico"` in ES)
- If the field name is NOT in the dictionary → **field name as-is** (pass-through)

The active language is controlled by `fmt.OutLang(lang)` globally.

To register translations for your domain words, call `fmt.RegisterWords()` in an `init()`:

```go
// In your package's words.go or init file:
func init() {
    fmt.RegisterWords([]fmt.DictEntry{
        {EN: "Name", ES: "Nombre", FR: "Nom", DE: "Name"},
        {EN: "Phone", ES: "Teléfono", FR: "Téléphone"},
    })
}
```

**Result**: Structs need zero tags for standard field names — less code, smaller WASM binaries.

## Example

```go
// Most fields need no tags — auto-translated from field name
type Product struct {
    Name     string                                      // placeholder: "Name" / "Nombre" (ES)
    Category string `options:"food:Food,tech:Technology"` // options still need tags
    Internal string `validate:"false"`                    // skip validation still needs tag
}

// Override only when translation is insufficient
type Order struct {
    Ref string `placeholder:"Order reference #"`
}
```

## Options Format

- Pairs separated by `,`
- Key and display text separated by `:`
- If no `:`, key and text are the same value

```
"admin:Administrator,user:Regular User,guest:Guest"
→ [{Key:"admin", Value:"Administrator"}, {Key:"user", Value:"Regular User"}, {Key:"guest", Value:"Guest"}]
```

Parsed with `fmt.Convert(tag).TagPairs("options")` from `tinywasm/fmt`.

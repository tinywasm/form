# API Reference

See `README.md` for the consolidated API. This file contains additional detail.

## `form.New` — Field Matching Detail

```go
f, err := form.New("content", &MyStruct{})
// → f.GetID() == "content.mystruct"
// → f.RenderHTML() renders all matching fields
// → err != nil if any exported field has no matching input
```

## `(*Form).Validate()` — Validation Detail

- Skips fields tagged `validate:"false"`
- Uses `GetSelectedValue()` to get current value per input
- Returns the **first** error encountered

## `(*Form).SyncValues()` — Binding Detail

Supports these struct field kinds: `string`, `[]string`, any type converted via `fmt.Convert`.

## `input.Permitted` — Validation Engine

```go
type Permitted struct {
    Letters         bool     // A-Z, a-z, Ñ
    Tilde           bool     // Á É Í Ó Ú á é í ó ú
    Numbers         bool     // 0-9
    WhiteSpaces     bool     // space ' '
    BreakLine       bool     // '\n'
    Tabulation      bool     // '\t'
    Characters      []rune   // extra allowed chars e.g. []rune{'@', '.'}
    TextNotAllowed  []string // blacklisted substrings
    Minimum         int      // minimum length (0 = no limit)
    Maximum         int      // maximum length (0 = no limit)
    ExtraValidation func(string) error // custom logic
    StartWith       *Permitted // rules for first character only
}
```

Error messages from `Permitted.Validate()`:
- `"minimum N chars"` — value shorter than Minimum
- `"maximum N chars"` — value longer than Maximum
- `"space not allowed"` — space when WhiteSpaces=false
- `"character X not allowed"` — disallowed character

## Tags Parsing Implementation

Tags are parsed in `form.go` inside `New()` using `tinywasm/fmt`:

```go
conv := fmt.Convert(fieldTag)
// options
opts := conv.TagPairs("options")          // []fmt.KeyValue
// scalar attrs
val, _ := conv.TagValue("placeholder")
val, _ = conv.TagValue("title")
val, _ = conv.TagValue("validate")       // "false" → skip validation
```

# Input Types

This package contains all input implementations for `tinywasm/form`.

## Default Inputs

| Type | HTML | Aliases | Validation |
|------|------|---------|------------|
| `Text` | `<input type="text">` | `name`, `fullname`, `username` | Letters, Numbers, Min: 2, Max: 100 |
| `Email` | `<input type="email">` | `mail`, `correo` | Letters, Numbers, `@._-`, Min: 5, Max: 100 |
| `Password` | `<input type="password">` | `pass`, `clave`, `pwd` | All chars, Min: 5, Max: 50 |

## Composition Pattern

Reuse base inputs to create semantic ones:

```go
func Gender(parentID, name string) input.Input {
    g := input.Radio(parentID, name).(*radio)
    g.Base.SetOptions(fmt.KeyValue{Key: "m", Value: "Male"}, ...)
    return g
}
```

## Matching Logic

The form engine matches fields in this order:
1. `LowerCase(FieldName)` vs `htmlName` or `aliases`
2. `LowerCase(StructName.FieldName)` vs `aliases`

## Registering Custom Inputs

If you create a new input type, register it to make it available:

```go
form.RegisterInput(MyCustomInput("", ""))
```

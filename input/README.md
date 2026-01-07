# Input Types

This package contains all input implementations for `tinywasm/form`.

## Default Inputs

| Type | HTML | Aliases | Validation |
|------|------|---------|------------|
| `Text` | `<input type="text">` | `name`, `fullname`, `username` | Letters, Numbers, Min: 2, Max: 100 |
| `Email` | `<input type="email">` | `mail`, `correo` | Letters, Numbers, `@._-`, Min: 5, Max: 100 |
| `Password` | `<input type="password">` | `pass`, `clave`, `pwd` | All chars, Min: 5, Max: 50 |

## Creating Custom Inputs

Implement the `Input` interface:

```go
type Input interface {
    dom.Component // ID(), RenderHTML()
    HtmlName() string
    ValidateField(value string) error
}
```

Embed `Base` and call `InitBase(id, name, htmlName, ...aliases)`:

```go
type myInput struct {
    input.Base
    input.Permitted
}

func MyInput(parentID, name string) input.Input {
    m := &myInput{...}
    m.Base.InitBase(parentID+"."+name, name, "text", "myalias")
    return m
}
```

Register with `form.RegisterInput(MyInput("", ""))`.

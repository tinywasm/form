# Struct Tags

The `tinywasm/form` library uses struct tags to configure inputs declaratively.

## `options` Tag

Defines options for multi-value inputs (select, radio, checkbox).

**Format**: `options:"key1:text1,key2:text2"`

### Example
```go
type User struct {
    Role string `options:"admin:Administrador,user:Usuario"`
}
```

### Parsing

Options are parsed by [ParseOptionsTag](../tags.go).

- Each option is separated by `,`.
- Key and text are separated by `:`.
- If no `:` is found, key and text are the same.

## Future Tags

Reserved for future use:
- `validate`: Custom validation rules.
- `label`: Custom display label.

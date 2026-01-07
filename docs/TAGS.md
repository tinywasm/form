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

Options are parsed using `fmt.Convert().TagPairs()` from the `tinywasm/fmt` library.

- Each pair is separated by `,`.
- Key and text are separated by `:`.
- If no `:` is found, key and text are the same.

See [tags.go](../tags.go) for the implementation.

## Other Attributes
Metadata tags like `placeholder`, `title`, etc., are automatically extracted if present.

## Future Tags
Reserved for future use:
- `validate`: Custom validation rules.
- `label`: Custom display label.

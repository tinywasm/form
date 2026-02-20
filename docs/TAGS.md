# Struct Tags

Tags configure inputs declaratively on struct fields.

## Supported Tags

| Tag | Format | Implemented |
|-----|--------|-------------|
| `options` | `"key1:text1,key2:text2"` | ✅ |
| `placeholder` | `"string"` | ✅ |
| `title` | `"string"` | ✅ |
| `validate` | `"false"` | ✅ (skips validation) |
| `label` | `"string"` | 🔲 reserved |

## Example

```go
type Product struct {
    Name     string `placeholder:"Product name" title:"Enter product name"`
    Category string `options:"food:Food,tech:Technology,other:Other"`
    Internal string `validate:"false"`
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

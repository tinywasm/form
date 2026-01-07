# Standard Input Types

Default inputs with aliases and validation rules.

## Registered Types

| Input | HTML | Aliases | Validation |
|-------|------|---------|------------|
| `Text` | `text` | `name`, `fullname`, `username` | Letters, Numbers, `. , ( )`, Min: 2, Max: 100 |
| `Email` | `email` | `mail`, `correo` | Letters, Numbers, `@ . _ -`, Min: 5, Max: 100 |
| `Password` | `password` | `pass`, `clave`, `pwd` | All chars, Min: 5, Max: 50 |

## Extensibility

Register custom inputs:

```go
form.RegisterInput(myCustomInput)
```

See [input/README.md](../input/README.md) for creating custom inputs.

## Validation (Permitted)

Whitelist approach:
- `Letters`: A-Z, a-z
- `Numbers`: 0-9
- `Characters`: Specific runes
- `Min/Max`: Length constraints
- `ExtraValidation`: Custom function

See [input/permitted.go](../input/permitted.go).

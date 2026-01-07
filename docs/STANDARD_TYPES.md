# Standard Input Types

Default inputs with aliases and validation rules.

## Registered Types

| Input | HTML | Aliases | Validation |
|-------|------|---------|------------|
| `Text` | `text` | `name`, `fullname`, `username` | Letters, Numbers, `. , ( )`, Spaces, Min: 2, Max: 100 |
| `Email` | `email` | `mail`, `correo` | Letters, Numbers, `@ . _ -`, Min: 5, Max: 100 |
| `Password` | `password` | `pass`, `clave`, `pwd` | All chars allowed, Min: 5, Max: 50 |
| `Address` | `text` | `address`, `direccion`, `dir`, `location` | Letters, Numbers, `# - / . , ( )`, Spaces, Min: 5, Max: 200 |
| `Radio` | `radio` | - | - |
| `Select` | `select` | `role`, `tipo` | - |
| `Gender` | `radio` | `gender`, `sexo` | `m:Male, f:Female` |

## Matching Logic
Matching uses the `HTMLName()` of the input or its predefined aliases. The comparison is case-insensitive.

## Validation (Permitted)
The validation engine uses a whitelist approach in [permitted.go](../input/permitted.go).

Current Rules:
- `Letters`: A-Z, a-z, Ñ
- `Numbers`: 0-9
- `WhiteSpaces`: Allows spaces (fixed in latest refactor)
- `Tilde`: ÁÉÍÓÚ
- `Characters`: List of specific allowed runes
- `Min/Max`: Length constraints
- `ExtraValidation`: User-defined callback

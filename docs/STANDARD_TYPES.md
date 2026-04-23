# Standard Input Types

18 built-in types in `input/`. Each provides `Xxx() fmt.Widget` and `Clone(parentID, name) fmt.Widget`.

| Input | HTML type | Constructor | Key Aliases |
|-------|-----------|-------------|-------------|
| `Address` | `text` | `input.Address()` | `address`, `addr` |
| `Checkbox` | `checkbox` | `input.Checkbox()` | `check`, `boolean`, `bool` |
| `Date` | `date` | `input.Date()` | `fecha` |
| `Datalist` | `text` | `input.Datalist()` | `list`, `options` |
| `Email` | `email` | `input.Email()` | `mail`, `correo` |
| `Filepath` | `text` | `input.Filepath()` | `path`, `dir`, `file` |
| `Gender` | `radio` | `input.Gender()` | `gender`, `sex` |
| `Hour` | `time` | `input.Hour()` | `hour` |
| `IP` | `text` | `input.IP()` | `ip` |
| `Number` | `number` | `input.Number()` | `num`, `amount`, `price`, `age` |
| `Password` | `password` | `input.Password()` | `pass`, `clave`, `pwd` |
| `Phone` | `tel` | `input.Phone()` | `phone`, `mobile`, `cell` |
| `Radio` | `radio` | `input.Radio()` | -- |
| `Rut` | `text` | `input.Rut()` | `rut`, `run`, `dni` |
| `Search` | `search` | `input.Search()` | `search`, `find`, `q` |
| `Select` | `select` | `input.Select()` | -- |
| `Text` | `text` | `input.Text()` | `name`, `fullname`, `username` |
| `Textarea` | `textarea` | `input.Textarea()` | `description`, `details`, `comments` |

See [input/README.md](../input/README.md) for detailed validation rules and rendering behavior per type.

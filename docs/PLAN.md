---
PLAN: "`tinywasm/form`: input.IP() doesn't validate real IP shape (accepts pure letters)"
---
> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.
> Repo rules: `AGENTS.md` at this repo's root — read it first (especially
> "Construction Harness" and the "No Go stdlib" / `tinywasm/fmt`-only rule).
> No master plan: scoped entirely to this repo, one file (`input/ip.go`) plus
> its test. Nothing else needs to change — `input.IP()`'s public shape
> (`func IP() Input`, `Validate(value string) error`) does not change.

## Context (zero-context summary)

`input.IP()` (`input/ip.go`) is meant to validate that a form field holds a
real IPv4 or IPv6 address. It does not. Found live, in a real app
(`tinywasm/layout`'s CRUD demo, `platformd/web/client.go`), by typing
`"uiuiuiu"` into a field configured with `input.IP()`: it saved without error.

**Root cause** — `Base.Permitted` (`model/permitted.go`) is a **character-
class whitelist**, not a grammar/shape validator: `Letters: true` means "every
character must be A-Z/a-z" (ANY letter, not just hex a-f — see
`Permitted.isAllowed`, `model/permitted.go:106-117`), and `Numbers: true`
means "every character must be 0-9". `IP()`'s constructor (`input/ip.go:13-22`)
sets `Numbers: true, Letters: true, Extra: []rune{'.', ':'}, Minimum: 7,
Maximum: 39` — so ANY string of 7-39 letters/digits/dots/colons passes the
whitelist. `ip.Validate` (`input/ip.go:25-46`) then only adds two checks: not
literally `"0.0.0.0"`, and not mixing `.` and `:` in the same value. Neither
requires the value to actually look like four numeric octets or hex groups.
`"uiuiuiu"` (7 letters, no dot, no colon) satisfies every existing check.

This is a real gap in a shared, public widget — not a demo-specific issue —
so it must be fixed in `tinywasm/form` (upstream), not worked around in a
consumer's model (the demo can't override `Letters`/`Numbers`; they're set
inside the constructor).

**Related, NOT in scope of this plan** (flagging for a separate look, do not
fix here): `input/email.go` and `input/phone.go`/`input/decimal.go`/
`input/number.go` don't override `Validate` at all — they rely purely on
`Base.Validate` (the same character-whitelist-only check), so e.g. `Email()`
accepts `"abc"` (5 letters, no `@`) as "valid". `input/date.go` and
`input/rut.go` already do real structural validation and are the reference
pattern this plan follows (see Stage 1).

## Stage 1 — real IPv4/IPv6 shape validation in `input/ip.go`

Replace `ip.Validate` (currently `input/ip.go:25-46`) with a version that
requires the value to actually be dot-separated (IPv4) or colon-separated
(IPv6) — not just built from permitted characters. Follow `input/date.go`'s
idiom: manual char-by-char parsing via `fmt.Convert(...).Int()`, no stdlib
`regexp`/`strings`/`strconv` (`AGENTS.md`: "No Go stdlib: use
`github.com/tinywasm/fmt`").

Rules to implement:

1. Keep the existing `Permitted.Validate` char-whitelist call first (cheap
   reject of anything with a truly stray character) and the existing
   `"0.0.0.0"` literal rejection (a product rule — that IPv4 numerically
   parses as valid four-octet form, so it must stay an explicit check, not
   rely on shape validation to catch it).
2. Classify by presence of `.` / `:`:
   - **Both present** → invalid (existing rule, keep as-is). Known, accepted
     limitation: this rejects IPv4-mapped IPv6 (`::ffff:192.168.1.1`) —
     out of scope for this fix (see Anti-footguns).
   - **Neither present** → invalid (NEW — today `"uiuiuiu"` falls exactly
     here and wrongly passes).
   - **Only `.`** → validate as IPv4: split on `.`, require exactly 4 parts;
     each part 1-3 ASCII digits (reject if any rune isn't `0-9`, reject empty
     parts); parse each with `fmt.Convert(part).Int()` and require `0-255`.
   - **Only `:`** → validate as IPv6: split on `:`; allow ONE empty part from
     `::` compression (an empty part from a leading/trailing/double colon);
     each non-empty part must be 1-4 hex characters (`0-9`, `a-f`, `A-F` —
     manual range check, same style as `date.go`'s digit check at
     `input/date.go:36`); total groups (counting a `::` compression as at
     least one zero-group) must be ≤ 8, and exactly 8 when there is no `::`.
3. Return the existing `fmt.Err("Format", "Invalid")` shape for every
   rejection (matches this package's existing error convention, e.g.
   `input/date.go:29,34,37`).

Sketch (adjust to compile — this is the shape, not verbatim final code):

```go
func (i *ip) Validate(value string) error {
	if value == "0.0.0.0" {
		return fmt.Err("Format", "Invalid")
	}
	if err := i.Permitted.Validate(i.name, value); err != nil {
		return err
	}
	hasDot, hasColon := false, false
	for _, c := range value {
		if c == '.' {
			hasDot = true
		}
		if c == ':' {
			hasColon = true
		}
	}
	switch {
	case hasDot && hasColon:
		return fmt.Err("Format", "Invalid")
	case !hasDot && !hasColon:
		return fmt.Err("Format", "Invalid")
	case hasDot:
		return validateIPv4(value)
	default:
		return validateIPv6(value)
	}
}
```

with `validateIPv4`/`validateIPv6` as unexported helpers in the same file
(package-private — `AGENTS.md`'s "minimal public surface").

## Stage 2 — tests

New `tests/ip_test.go` (no existing `ip` test file — confirmed via
`grep -rln "input.IP" tests/`, empty). Follow the shape of
`tests/load_test.go`/other `input_test.go` files already in `tests/`. Cover,
at minimum:

- The bug as reported: `Validate("uiuiuiu")` → error (currently nil).
- Valid IPv4: `"192.168.1.1"`, `"0.0.0.1"` → nil.
- Invalid IPv4: `"999.1.1.1"` (out of range), `"1.2.3"` (only 3 groups),
  `"1.2.3.4.5"` (5 groups), `"1.2.3.a"` (non-digit group) → error.
- `"0.0.0.0"` → error (existing rule, keep covered).
- Valid IPv6: `"::1"`, `"fe80::1"`, `"2001:db8::8a2e:370:7334"` → nil.
- Invalid IPv6: `"gggg::1"` (non-hex group), too many groups → error.
- Mixed `.`+`:` (e.g. `"::ffff:1.2.3.4"`) → error (documented limitation,
  not a regression — assert the CURRENT behavior, don't try to support it).
- Neither `.` nor `:` (`"uiuiuiu"`, `"12345"`) → error.
- `gotest` (never `go test`) green: vet/race/tests/wasm/coverage.

## Anti-footguns (do NOT do)

- **Do NOT attempt full RFC 4291 IPv6 support** (IPv4-mapped addresses,
  zone IDs like `%eth0`, multiple `::` detection beyond "at most one").
  This is a device-management form field, not a network stack — "obviously
  not garbage" is the bar, not full spec compliance. If asked to go further,
  stop and ask rather than expanding silently.
- **Do NOT touch `input/email.go`, `input/phone.go`, `input/decimal.go`,
  `input/number.go`** — their similar charset-only gap is a SEPARATE,
  unverified finding (see Context) and explicitly out of scope here.
- **Do NOT change `IP()`'s public constructor signature or the `Minimum`/
  `Maximum`/`Letters`/`Numbers`/`Extra` field values** — only `Validate`'s
  body changes. Consumers that already call `input.IP()` must keep compiling
  and behaving the same for already-valid values.
- **Do NOT add `regexp`, `strings`, or `strconv`** — this package is
  `tinywasm/fmt`-only per `AGENTS.md`; `date.go`/`rut.go` are the reference
  idiom for manual parsing.
- Never run `gopush` or `codejob` from this plan — this repo's maintainer
  dispatches it themselves (see the note at the top: no master plan, single
  repo, single file).

## Stages table

| # | Stage | Files | Done |
|---|---|---|---|
| 1 | Real IPv4/IPv6 shape validation | `input/ip.go` | ☐ |
| 2 | Tests | `tests/ip_test.go` | ☐ |

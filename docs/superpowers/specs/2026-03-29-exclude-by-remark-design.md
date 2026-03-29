# Exclude by Remark

Exclude VLESS outbounds from output based on wildcard patterns matched against their remark field.

## Motivation

Subscriptions may contain outbounds that the user does not want in the output (e.g., servers in specific countries). Currently there is no way to filter by remark — the user gets all matching VLESS outbounds.

## CLI Flag

`--exclude-by-remark <pattern>` — can be specified multiple times. Each value is a glob pattern (per `filepath.Match` semantics: `*`, `?`, `[...]`). Matching is case-insensitive.

If any pattern is syntactically invalid (e.g., `[unclosed`), the program prints an error to stderr and exits with code 1 before fetching the subscription.

Example:

```
goxsub --exclude-by-remark='*Russia*' --exclude-by-remark='*China*' <url>
```

## Library

New file: `xray/filter.go`.

```go
func FilterByRemark(proxies []VLESSProxy, patterns []string) []VLESSProxy
```

Returns proxies whose `Remarks` do not match any of the given patterns. Matching is case-insensitive using `strings.ToLower` on both sides, with `filepath.Match` for glob evaluation.

Precondition: all patterns must be valid (the caller validates them before calling this function).

## Integration

In `cmd/goxsub/main.go`, after `ExtractVLESSOutbounds` and before output formatting:

```
proxies = xray.FilterByRemark(proxies, excludePatterns)
```

## Testing

Tests in `xray/filter_test.go`:

- Empty patterns list returns all proxies unchanged.
- A single pattern excludes matching proxies.
- Multiple patterns: proxy excluded if it matches any one of them.
- Case-insensitive matching: `*russia*` matches `RUSSIA`, `Russia`, `russia`.
- Glob specials: `?`, `[...]` work as expected.
- No matches returns empty slice.
- All excluded returns empty slice.

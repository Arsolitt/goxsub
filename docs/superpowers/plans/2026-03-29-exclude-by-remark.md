# Exclude by Remark Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `--exclude-by-remark` CLI flag that filters out VLESS outbounds whose remark matches a glob pattern.

**Architecture:** Library function `xray.FilterByRemark` in a new file `xray/filter.go`. CLI validates patterns and calls the filter between extraction and output formatting.

**Tech Stack:** Go stdlib (`flag`, `path/filepath`, `strings`)

---

## File Structure

| File | Action | Responsibility |
|------|--------|---------------|
| `xray/filter.go` | Create | `FilterByRemark` function |
| `xray/filter_test.go` | Create | Tests for `FilterByRemark` |
| `cmd/goxsub/main.go` | Modify | Add `--exclude-by-remark` flag, validate patterns, call filter |

---

### Task 1: FilterByRemark — Write Failing Tests

**Files:**
- Create: `xray/filter_test.go`

- [ ] **Step 1: Write the failing tests**

```go
package xray

import "testing"

func makeProxy(remarks string) VLESSProxy {
	return VLESSProxy{
		Remarks: remarks,
		Outbound: Outbound{
			Protocol: "vless",
			Settings: OutboundSettings{Vnext: []VNext{{Address: "a.com", Port: 443, Users: []User{{ID: "u"}}}}},
		},
	}
}

func TestFilterByRemark_EmptyPatterns(t *testing.T) {
	proxies := []VLESSProxy{makeProxy("Russia Server"), makeProxy("NL Server")}
	result := FilterByRemark(proxies, nil)
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestFilterByRemark_SinglePattern(t *testing.T) {
	proxies := []VLESSProxy{makeProxy("Russia Server"), makeProxy("NL Server")}
	result := FilterByRemark(proxies, []string{"*Russia*"})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Remarks != "NL Server" {
		t.Errorf("expected 'NL Server', got %q", result[0].Remarks)
	}
}

func TestFilterByRemark_MultiplePatterns(t *testing.T) {
	proxies := []VLESSProxy{makeProxy("Russia Server"), makeProxy("China Node"), makeProxy("NL Server")}
	result := FilterByRemark(proxies, []string{"*Russia*", "*China*"})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Remarks != "NL Server" {
		t.Errorf("expected 'NL Server', got %q", result[0].Remarks)
	}
}

func TestFilterByRemark_CaseInsensitive(t *testing.T) {
	proxies := []VLESSProxy{makeProxy("RUSSIA"), makeProxy("russia"), makeProxy("Russia"), makeProxy("NL")}
	result := FilterByRemark(proxies, []string{"*russia*"})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Remarks != "NL" {
		t.Errorf("expected 'NL', got %q", result[0].Remarks)
	}
}

func TestFilterByRemark_GlobSpecials(t *testing.T) {
	proxies := []VLESSProxy{makeProxy("A1"), makeProxy("A2"), makeProxy("B1")}
	result := FilterByRemark(proxies, []string{"A?"})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Remarks != "B1" {
		t.Errorf("expected 'B1', got %q", result[0].Remarks)
	}
}

func TestFilterByRemark_NoMatches(t *testing.T) {
	proxies := []VLESSProxy{makeProxy("NL Server"), makeProxy("DE Server")}
	result := FilterByRemark(proxies, []string{"*JP*"})
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestFilterByRemark_AllExcluded(t *testing.T) {
	proxies := []VLESSProxy{makeProxy("RUSSIA"), makeProxy("Russia")}
	result := FilterByRemark(proxies, []string{"*Russia*"})
	if len(result) != 0 {
		t.Fatalf("expected 0, got %d", len(result))
	}
}

func TestFilterByRemark_NilInput(t *testing.T) {
	result := FilterByRemark(nil, []string{"*Russia*"})
	if len(result) != 0 {
		t.Fatalf("expected 0, got %d", len(result))
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./xray/ -run TestFilterByRemark -v`
Expected: compilation error — `FilterByRemark` not defined

- [ ] **Step 3: Commit failing tests**

```bash
git add xray/filter_test.go
git commit -m "test: add failing tests for FilterByRemark"
```

---

### Task 2: FilterByRemark — Implement

**Files:**
- Create: `xray/filter.go`

- [ ] **Step 1: Write the implementation**

```go
package xray

import (
	"path/filepath"
	"strings"
)

func FilterByRemark(proxies []VLESSProxy, patterns []string) []VLESSProxy {
	if len(patterns) == 0 || len(proxies) == 0 {
		return proxies
	}

	var result []VLESSProxy
	for _, p := range proxies {
		excluded := false
		remark := strings.ToLower(p.Remarks)
		for _, pattern := range patterns {
			matched, _ := filepath.Match(strings.ToLower(pattern), remark)
			if matched {
				excluded = true
				break
			}
		}
		if !excluded {
			result = append(result, p)
		}
	}
	return result
}
```

- [ ] **Step 2: Run tests to verify they pass**

Run: `go test ./xray/ -run TestFilterByRemark -v`
Expected: all 8 tests PASS

- [ ] **Step 3: Run full test suite**

Run: `go test ./...`
Expected: all tests PASS

- [ ] **Step 4: Lint**

Run: `golangci-lint run --fix`
Expected: no errors

- [ ] **Step 5: Commit**

```bash
git add xray/filter.go
git commit -m "feat: add FilterByRemark for glob-based remark exclusion"
```

---

### Task 3: CLI Integration

**Files:**
- Modify: `cmd/goxsub/main.go`

- [ ] **Step 1: Add the flag and filter logic**

In `cmd/goxsub/main.go`, add imports `"fmt"`, `"path/filepath"`, `"strings"` are already imported. Add `"fmt"` is already there. Need to verify all imports.

Replace the flag parsing section (lines 18-21) and add filter logic after line 71. The full modified `run()` function:

```go
func run() int {
	format := flag.String("format", "uri", "output format: uri, podkop")
	podkopSection := flag.String("podkop-section", "main", "podkop uci section name")
	var excludePatterns stringSlice
	flag.Var(&excludePatterns, "exclude-by-remark", "exclude outbounds by remark glob pattern (case-insensitive, repeatable)")
	flag.Parse()

	if *format != "podkop" && flag.Lookup("podkop-section").DefValue != *podkopSection {
		fmt.Fprintf(os.Stderr, "error: --podkop-section can only be used with --format podkop\n")
		return 1
	}

	for _, p := range excludePatterns {
		if _, err := filepath.Match(p, ""); err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid glob pattern %q: %v\n", p, err)
			return 1
		}
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: goxsub [flags] <subscription-url>\n")
		fmt.Fprintf(os.Stderr, "flags:\n")
		flag.PrintDefaults()
		return 1
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		args[0],
		nil,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "error: HTTP %d\n", resp.StatusCode)
		return 1
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	subs, err := xray.ParseSubscription(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	proxies := xray.ExtractVLESSOutbounds(subs)
	proxies = xray.FilterByRemark(proxies, excludePatterns)

	switch *format {
	case "podkop":
		output, err := xray.FormatPodkop(proxies, *podkopSection)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		fmt.Println(output)
	default:
		for _, p := range proxies {
			uri, err := xray.ToVLESSURI(p.Outbound, p.Remarks)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				return 1
			}
			fmt.Println(uri)
		}
	}

	return 0
}
```

Also add the `stringSlice` type and the `"path/filepath"` import. Full imports:

```go
import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Arsolitt/goxsub/xray"
)

type stringSlice []string

func (s *stringSlice) String() string { return strings.Join(*s, ", ") }

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}
```

- [ ] **Step 2: Build**

Run: `go build -o build/goxsub ./cmd/goxsub/`
Expected: build succeeds

- [ ] **Step 3: Test help output**

Run: `./build/goxsub -h`
Expected: `--exclude-by-remark` appears in flags listing

- [ ] **Step 4: Test invalid pattern**

Run: `./build/goxsub --exclude-by-remark='[unclosed' http://example.com`
Expected: stderr contains "invalid glob pattern", exit code 1

- [ ] **Step 5: Run full test suite**

Run: `go test ./...`
Expected: all tests PASS

- [ ] **Step 6: Lint**

Run: `golangci-lint run --fix`
Expected: no errors

- [ ] **Step 7: Commit**

```bash
git add cmd/goxsub/main.go
git commit -m "feat: add --exclude-by-remark flag for glob-based remark exclusion"
```

package proxy

import (
	"path/filepath"
	"strings"
)

// FilterByRemark returns proxies whose Remarks field does not match any of the given glob patterns.
// Matching is case-insensitive. If patterns is empty, all proxies are returned unchanged.
func FilterByRemark(proxies []Proxy, patterns []string) []Proxy {
	if len(patterns) == 0 || len(proxies) == 0 {
		return proxies
	}

	var result []Proxy
	for _, p := range proxies {
		excluded := false
		remark := strings.ToLower(p.Remarks())
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

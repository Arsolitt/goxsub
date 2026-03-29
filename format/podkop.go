package format

import (
	"fmt"
	"strings"

	"github.com/Arsolitt/goxsub/protocol"
	"github.com/Arsolitt/goxsub/proxy"
)

// Podkop converts proxies to uci shell commands for podkop OpenWrt package.
func Podkop(proxies []proxy.Proxy, section string) (string, error) {
	var lines []string
	lines = append(lines, fmt.Sprintf("uci del podkop.%s.urltest_proxy_links", section))
	for _, p := range proxies {
		uri, err := protocol.ToURI(p)
		if err != nil {
			return "", fmt.Errorf("convert proxy to URI: %w", err)
		}
		lines = append(lines, fmt.Sprintf("uci add_list podkop.%s.urltest_proxy_links='%s'", section, uri))
	}
	lines = append(lines, "service podkop restart")
	return strings.Join(lines, "\n"), nil
}

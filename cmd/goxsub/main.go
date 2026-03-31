package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	goxsub "github.com/Arsolitt/goxsub"
)

type stringSlice []string

func (s *stringSlice) String() string { return strings.Join(*s, ", ") }

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	os.Exit(run())
}

//nolint:funlen,gocognit
func run() int {
	formatFlag := flag.String("format", "uri", "output format: uri, podkop, singbox")
	podkopSection := flag.String("podkop-section", "main", "podkop uci section name")
	singboxDNSResolver := flag.String("singbox-dns-resolver", "dns-local", "sing-box domain_resolver value")
	singboxOutboundPrefix := flag.String("singbox-outbound-prefix", "", "sing-box outbound tag prefix")
	singboxOutboundSuffix := flag.String("singbox-outbound-suffix", "", "sing-box outbound tag suffix")
	keepRemark := flag.Bool("keep-remark", true, "keep original remark or replace with sequential number")
	var excludePatterns stringSlice
	flag.Var(
		&excludePatterns,
		"exclude-by-remark",
		"exclude outbounds by remark glob pattern (case-insensitive, repeatable)",
	)
	flag.Parse()

	if *formatFlag != "podkop" && flag.Lookup("podkop-section").DefValue != *podkopSection {
		fmt.Fprintf(os.Stderr, "error: --podkop-section can only be used with --format podkop\n")
		return 1
	}

	if *formatFlag != "singbox" && (flag.Lookup("singbox-dns-resolver").DefValue != *singboxDNSResolver ||
		*singboxOutboundPrefix != "" || *singboxOutboundSuffix != "") {
		fmt.Fprintf(os.Stderr, "error: --singbox-* flags can only be used with --format singbox\n")
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

	subs, err := goxsub.ParseSubscription(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	proxies := goxsub.ExtractProxies(subs)
	proxies = goxsub.FilterByRemark(proxies, excludePatterns)

	if !*keepRemark {
		for i, p := range proxies {
			if vp, ok := p.(*goxsub.VLESSProxy); ok {
				vp.SetRemarks(strconv.Itoa(i + 1))
			}
		}
	}

	switch *formatFlag {
	case "podkop":
		lines, err := goxsub.Podkop(proxies, *podkopSection)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		for _, line := range lines {
			fmt.Println(line)
		}
	case "singbox":
		cfg := goxsub.SingboxConfig{
			OutboundPrefix: *singboxOutboundPrefix,
			OutboundSuffix: *singboxOutboundSuffix,
			KeepRemark:     *keepRemark,
			DNSResolver:    *singboxDNSResolver,
		}
		lines, err := goxsub.Singbox(proxies, cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		for _, line := range lines {
			fmt.Println(line + ",")
		}
	default:
		for _, p := range proxies {
			uri, err := goxsub.ToURI(p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				return 1
			}
			fmt.Println(uri)
		}
	}

	return 0
}

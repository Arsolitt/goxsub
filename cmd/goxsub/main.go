package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Arsolitt/goxsub/xray"
)

const expectedArgs = 2

func main() {
	os.Exit(run())
}

func run() int {
	if len(os.Args) != expectedArgs {
		fmt.Fprintf(os.Stderr, "usage: goxsub <subscription-url>\n")
		return 1
	}

	req, err := http.NewRequestWithContext( //nolint:gosec // G704: CLI tool intentionally fetches user-provided URL
		context.Background(),
		http.MethodGet,
		os.Args[1],
		nil,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	resp, err := http.DefaultClient.Do(req) //nolint:gosec // G704: CLI tool intentionally fetches user-provided URL
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
	for _, p := range proxies {
		uri, err := xray.ToVLESSURI(p.Outbound, p.Remarks)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		fmt.Println(uri)
	}

	return 0
}

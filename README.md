<div align="center">

# goxsub

[![Go Reference](https://pkg.go.dev/badge/github.com/Arsolitt/goxsub.svg)](https://pkg.go.dev/github.com/Arsolitt/goxsub)
[![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: GPL-3.0](https://img.shields.io/badge/License-GPL--3.0-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![No Dependencies](https://img.shields.io/badge/dependencies-none-brightgreen)](https://github.com/Arsolitt/goxsub)
[![Codecov](https://img.shields.io/codecov/c/github/Arsolitt/goxsub?token=YOUR_CODECOV_TOKEN&logo=codecov&logoColor=white)](https://codecov.io/gh/Arsolitt/goxsub)
[![golangci-lint](https://img.shields.io/badge/linting-golangci--lint-00B4D8?logo=go&logoColor=white)](https://golangci-lint.run/)

**Go library and CLI for parsing xray core JSON subscriptions and converting vless outbounds to vless:// URIs, sing-box outbounds, or podkop UCI commands.**

Standard library only. No third-party dependencies.

</div>

## Supported

- **Protocol**: VLESS
- **Transports**: TCP, WebSocket, gRPC, HTTP/2, mKCP
- **Security**: none, TLS, REALITY

## Quick Start

```shell
go install github.com/Arsolitt/goxsub/cmd/goxsub@latest
goxsub https://example.com/subscription
```

## CLI Usage

```shell
goxsub [flags] <subscription-url>
```

| Flag | Default | Description |
|------|---------|-------------|
| `-format` | `uri` | Output format: `uri`, `podkop`, or `singbox` |
| `-podkop-section` | `main` | Podkop UCI section name |
| `-keep-remark` | `true` | Keep original remark or replace with sequential number |
| `-singbox-dns-resolver` | `dns-local` | sing-box domain_resolver value |
| `-singbox-outbound-prefix` | — | sing-box outbound tag prefix |
| `-singbox-outbound-suffix` | — | sing-box outbound tag suffix |
| `-exclude-by-remark` | — | Exclude proxies by remark glob, case-insensitive (repeatable) |

## Library Usage

```go
import goxsub "github.com/Arsolitt/goxsub"

subs, _ := goxsub.ParseSubscription(data)
proxies := goxsub.ExtractProxies(subs)
filtered := goxsub.FilterByRemark(proxies, []string{"blocklist*"})
uri, _ := goxsub.ToURI(filtered[0])
```

## Architecture

```
api.go         Re-exports all public types and functions
cmd/goxsub/    CLI binary
sub/           JSON subscription parsing and types
proxy/         Proxy extraction, filtering, and Proxy interface
protocol/      URI conversion (vless://)
format/        Output formatters (podkop UCI commands, sing-box outbound JSON)
```

Data flow: `JSON → ParseSubscription → ExtractProxies → FilterByRemark → ToURI/Podkop/Singbox`

## Using with AI Assistants

It is recommended to copy the following prompt and send it to an AI assistant — this can significantly improve the quality of generated configurations:

```
https://github.com/Arsolitt/goxsub/blob/main/llms-full.txt This link is the full documentation of goxsub.

【Role Setting】
You are an expert proficient in Xray-core proxy configuration and goxsub library usage.

【Task Requirements】
1. Knowledge Base: Please read and deeply understand the content of this link, and use it as the sole basis for answering questions and writing configurations.
2. No Hallucinations: Absolutely do not fabricate fields that do not exist in the documentation. If the documentation does not mention it, please tell me directly "Documentation does not mention".
3. Default Format: Output vless:// URIs by default (unless I explicitly request a different format), and add key comments.
4. Exception Handling: If you cannot access this link, please inform me clearly and prompt me to manually download the documentation and upload it to you.
```

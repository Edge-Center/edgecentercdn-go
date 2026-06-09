# edgecentercdn-go

[![CI](https://github.com/Edge-Center/edgecentercdn-go/actions/workflows/ci.yml/badge.svg)](https://github.com/Edge-Center/edgecentercdn-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/Edge-Center/edgecentercdn-go.svg)](https://pkg.go.dev/github.com/Edge-Center/edgecentercdn-go)

`edgecentercdn-go` is the EdgeCenter CDN API SDK for the Go programming language.
It ships both a library for building integrations (including the Terraform
provider) and an `edge-cli` command-line client.

## Installation

```sh
go get github.com/Edge-Center/edgecentercdn-go
```

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	edgecentercdn "github.com/Edge-Center/edgecentercdn-go"
	"github.com/Edge-Center/edgecentercdn-go/edgecenter/provider"
	"github.com/Edge-Center/edgecentercdn-go/resources"
)

func main() {
	client := provider.NewClient(
		"https://api.edgecenter.ru",
		provider.WithSignerFunc(func(req *http.Request) error {
			for k, v := range provider.AuthenticatedHeaders("<API-KEY>") {
				req.Header.Set(k, v)
			}
			return nil
		}),
	)

	svc := edgecentercdn.NewService(client)

	list, err := svc.Resources().List(context.Background(), &resources.ListFilterRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d resources\n", len(list))
}
```

API errors are typed: use `errors.Is` with the sentinels in `edgecenter`
(`ErrNotFound`, `ErrConflict`, `ErrRateLimit`, ...) and `errors.As` with
`*edgecenter.APIError` to inspect the status code and field-level details.

## CLI

```sh
go build -o bin/edge-cli ./cmd      # or: make build-cli

export EC_CDN_API_KEY=<your-api-key>
export EC_CDN_API_URL=https://api.edgecenter.ru

bin/edge-cli resource list
bin/edge-cli origins get --id 123
bin/edge-cli purge --resource 123 --path /static/app.js
```

Pre-built binaries for macOS, Linux, Windows and FreeBSD are attached to each
[GitHub release](https://github.com/Edge-Center/edgecentercdn-go/releases).

## Development

Common tasks are wrapped in the `Makefile`:

| Target | What it does |
|--------|--------------|
| `make build` | Build all packages |
| `make test` | Run unit tests with the race detector |
| `make cover` | Generate `coverage.html` |
| `make cover-check` | Fail if SDK coverage drops below 70% |
| `make lint` | Run `golangci-lint` |
| `make fmt` | Format with `gofmt` + `goimports` |
| `make vet` | Run `go vet` |
| `make tidy` | `go mod tidy` |
| `make snapshot` | Build release artifacts locally via GoReleaser |
| `make install-tools` | Install the dev tools used above |

CI (`.github/workflows/ci.yml`) runs lint, the test suite across Go 1.25 and
1.26, the coverage gate, and a build check on every pull request.

## Testing

Each SDK service is covered by unit tests that run against an in-process
`httptest` server, so no network or credentials are required:

```sh
make test
```

The coverage gate (`make cover-check`) protects the SDK packages. CLI commands
are covered separately by the end-to-end test track.

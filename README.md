# Tempest Developer SDK for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/tempestdx/sdk-go)](https://pkg.go.dev/github.com/tempestdx/sdk-go)
[![Test Status](https://github.com/tempestdx/sdk-go/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/tempestdx/sdk-go/actions/workflows/go.yml?query=branch%3Amain)

The official [Tempest][tempest] SDK for Go.

## Requirements

- Go 1.23 or later

## Installation

Make sure your project is using Go Modules (it will have a `go.mod` file in its
root if it already is):

```sh
go mod init
```

Then, reference sdk-go in a Go program with `import`:

```go
import (
    "github.com/tempestdx/sdk-go/app"
)
```

For more information on how to use the Tempest SDK for Go, see our
[Hello World][hello-world] guide.

## Documentation

For details on all the functionality in this SDK, see our
[Go documentation][goref].

## Support

New features and bug fixes are released on the latest version of the Tempest SDK
library. If you're using an older major version, we recommend updating to the
latest version to access new features, benefit from recent bug fixes, and ensure
you have the latest security patches. Older major versions of the SDK will
continue to be available for use, but will not receive any further updates.

## Development

Pull requests from the community are welcome. If you submit one, please keep the
following guidelines in mind:

1. Code must be `go fmt` compliant.
2. All types, structs and funcs should be documented.
3. Ensure that `go test` succeeds.

## Test

The test suite needs testify's `require` package to run:

    github.com/stretchr/testify/require

Before running any tests, make sure to grab all of the package's dependencies:

    go get -t -v ./...

Run all tests:

    go test -race -v ./...

Run tests for one package:

    go test -v ./app/...

Run a single test:

    go test -v ./app/... -run TestHealthCheck

To share any requests, bugs or comments, please [open an issue][issues] or
[submit a pull request][pulls].

[goref]: https://pkg.go.dev/github.com/tempestdx/sdk-go
[issues]: https://github.com/tempestdx/sdk-go/issues/new
[pulls]: https://github.com/tempestdx/sdk-go/pulls
[tempest]: https://tempestdx.com/
[hello-world]: https://docs.tempestdx.com/developer/guides/hello-world

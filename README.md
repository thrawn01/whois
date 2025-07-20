# Whois

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![GoDoc](https://pkg.go.dev/badge/github.com/thrawn01/whois.svg)](https://pkg.go.dev/github.com/thrawn01/whois)
[![Go Report Card](https://goreportcard.com/badge/github.com/thrawn01/whois)](https://goreportcard.com/report/github.com/thrawn01/whois)
[![Build Status](https://github.com/thrawn01/whois/actions/workflows/gotest.yaml/badge.svg)](https://github.com/thrawn01/whois/actions/workflows/gotest.yaml)

Whois is a simple Go module for domain and ip whois information query.

> **Note**: This is a fork of [likexian/whois](https://github.com/likexian/whois) with the following improvements:
> - Applied all unmerged pull requests from the original repository
> - Removed telemetry/update checking that reports usage to external servers
> - Added context support for all queries (PR #54)
> - Fixed unreachable referral server handling (PR #58)
> - Updated dependencies (PR #57)

## Overview

All of domain, IP include IPv4 and IPv6, ASN are supported.

You can directly using the binary distributions whois, follow [whois release tool](cmd/whois).

Or you can do development by using this golang module as below.

## Installation

### As a library

```shell
go get -u github.com/thrawn01/whois
```

### As a command-line tool

```shell
go install github.com/thrawn01/whois/cmd/whois@latest
```

## Importing

```go
import (
    "github.com/thrawn01/whois"
)
```

## Documentation

Visit the docs on [GoDoc](https://pkg.go.dev/github.com/thrawn01/whois)

## Example

### whois query for domain

```go
result, err := whois.Whois("likexian.com")
if err == nil {
    fmt.Println(result)
}
```

### whois query for IPv6

```go
result, err := whois.Whois("2001:dc7::1")
if err == nil {
    fmt.Println(result)
}
```

### whois query for IPv4

```go
result, err := whois.Whois("1.1.1.1")
if err == nil {
    fmt.Println(result)
}
```

### whois query for ASN

```go
// or whois.Whois("AS60614")
result, err := whois.Whois("60614")
if err == nil {
    fmt.Println(result)
}
```
## CLI Usage

### whois query for domain

```shell
whois likexian.com
```

### whois query for IPv6

```shell
whois 2001:dc7::1
```

### whois query for IPv4

```shell
whois 1.1.1.1
```

### whois query for ASN

```shell
# or whois as60614
whois 60614
```

### whois query output as json

```shell
whois -j likexian.com
```

## Whois information parsing

Please refer to [whois-parser](https://github.com/likexian/whois-parser)

## License

- Copyright 2014-2024 [Li Kexian](https://www.likexian.com/)
- Copyright 2025 Derrick Wippler

Licensed under the Apache License 2.0

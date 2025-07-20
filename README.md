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

## MCP Server Integration

This package includes an MCP (Model Context Protocol) server that allows Claude Code and other MCP-compatible clients to perform whois queries.

### Installation

You can run the MCP server directly using Go:

```shell
go run github.com/thrawn01/whois/cmd/mcp-whois@latest
```

Or install it locally:

```shell
go install github.com/thrawn01/whois/cmd/mcp-whois@latest
```

### Configuration for Claude Code

Add the following to your Claude Code configuration:

```json
{
  "mcpServers": {
    "whois": {
      "command": "go",
      "args": ["run", "github.com/thrawn01/whois/cmd/mcp-whois@latest"]
    }
  }
}
```

Or if you've installed it locally:

```json
{
  "mcpServers": {
    "whois": {
      "command": "mcp-whois"
    }
  }
}
```

### Available Tools

#### whois_lookup
Performs whois lookups for domains, IP addresses, and ASNs.

**Parameters:**
- `query` (required): Domain, IP address, or ASN to lookup
- `server` (optional): Specific whois server to use
- `timeout` (optional): Query timeout in seconds (default: 30)
- `disable_referral` (optional): Disable referral server queries
- `parse_json` (optional): Return parsed JSON format instead of raw whois data

**Example usage in Claude:**
```
Please lookup whois information for example.com
```

### Available Resources

The MCP server also provides resource URIs in the format `whois://[query]` for direct access to whois data.

**Example:**
```
Access the resource whois://example.com
```

### Testing the MCP Server

The MCP server includes comprehensive tests. To run them:

```shell
# Run all tests
go test ./cmd/mcp-whois

# Run tests with verbose output
go test -v ./cmd/mcp-whois

# Run benchmarks
go test -bench=. ./cmd/mcp-whois
```

The test suite includes:
- Unit tests for the whois_lookup tool
- Integration tests with mock whois servers
- Error handling and timeout tests  
- JSON parsing validation tests
- Parameter validation tests

## License

- Copyright 2014-2024 [Li Kexian](https://www.likexian.com/)
- Copyright 2025 Derrick Wippler

Licensed under the Apache License 2.0

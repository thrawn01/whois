/*
 * Copyright 2025 Derrick Wippler
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * MCP server for whois queries
 */

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	whoisparser "github.com/likexian/whois-parser"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/thrawn01/whois"
)

// WhoisLookupParams defines the parameters for whois lookup tool
type WhoisLookupParams struct {
	Query           string `json:"query" jsonschema:"the domain, IP address, or ASN to lookup"`
	Server          string `json:"server,omitempty" jsonschema:"specific whois server to use (optional)"`
	Timeout         int    `json:"timeout,omitempty" jsonschema:"query timeout in seconds (default: 30)"`
	DisableReferral bool   `json:"disable_referral,omitempty" jsonschema:"disable referral server queries"`
	ParseJSON       bool   `json:"parse_json,omitempty" jsonschema:"return parsed JSON format instead of raw whois data"`
}

func main() {
	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "whois-mcp",
		Version: "1.0.0",
	}, nil)

	// Add whois lookup tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "whois_lookup",
		Description: "Perform whois lookups for domains, IP addresses, and ASNs",
	}, handleWhoisLookup)

	// Run the server
	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		log.Fatal(err)
	}
}

// handleWhoisLookup handles the whois_lookup tool requests
func handleWhoisLookup(_ context.Context, _ *mcp.ServerSession, params *mcp.CallToolParamsFor[WhoisLookupParams]) (*mcp.CallToolResultFor[any], error) {
	args := params.Arguments

	// Create whois client
	client := whois.NewClient()

	// Set timeout if specified
	if args.Timeout > 0 {
		client.SetTimeout(time.Duration(args.Timeout) * time.Second)
	}

	// Set disable referral if specified
	if args.DisableReferral {
		client.SetDisableReferral(true)
	}

	// Perform whois query
	var result string
	var err error
	if args.Server != "" {
		result, err = client.Whois(args.Query, args.Server)
	} else {
		result, err = client.Whois(args.Query)
	}

	if err != nil {
		return nil, fmt.Errorf("whois query failed: %w", err)
	}

	// Return parsed JSON if requested
	if args.ParseJSON {
		parsed, parseErr := whoisparser.Parse(result)
		if parseErr != nil {
			// If parsing fails, return raw result with a note
			return &mcp.CallToolResultFor[any]{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("Failed to parse whois data as JSON: %v\n\nRaw whois data:\n%s", parseErr, result),
					},
				},
			}, nil
		}

		jsonData, err := json.MarshalIndent(parsed, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal parsed data: %w", err)
		}

		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: string(jsonData),
				},
			},
		}, nil
	}

	// Return raw whois data
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: result,
			},
		},
	}, nil
}

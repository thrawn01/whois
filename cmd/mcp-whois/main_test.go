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
 * Tests for MCP server for whois queries
 */

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestHandleWhoisLookup(t *testing.T) {
	// Create a mock whois server for all tests
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create mock server: %v", err)
	}
	defer func() { _ = listener.Close() }()

	address := listener.Addr().(*net.TCPAddr)
	mockServer := fmt.Sprintf("127.0.0.1:%d", address.Port)

	// Start mock server
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer func() { _ = c.Close() }()

				// Read the query
				buffer := make([]byte, 1024)
				n, err := c.Read(buffer)
				if err != nil {
					return
				}

				query := strings.TrimSpace(string(buffer[:n]))

				// Handle different query types
				var response string
				if query == "" {
					// Empty query - close without response
					return
				} else if strings.Contains(query, "invalid") {
					// Invalid domain - return error-like response
					response = "No match for domain"
				} else {
					// Valid query - return mock whois data
					response = fmt.Sprintf(`Domain Name: %s
Registry Domain ID: 123456_DOMAIN_COM-VRSN
Registrar WHOIS Server: whois.example.com
Creation Date: 2024-01-01T00:00:00Z
Registry Expiry Date: 2025-01-01T00:00:00Z
Registrar: Test Registrar
`, strings.ToUpper(query))
				}

				_, _ = c.Write([]byte(response))
			}(conn)
		}
	}()

	// Small delay to ensure server is ready
	time.Sleep(10 * time.Millisecond)

	tests := []struct {
		name          string
		params        WhoisLookupParams
		expectError   bool
		expectContent bool
	}{
		{
			name: "validDomain",
			params: WhoisLookupParams{
				Query:  "example.com",
				Server: mockServer,
			},
			expectError:   false,
			expectContent: true,
		},
		{
			name: "validDomainWithTimeout",
			params: WhoisLookupParams{
				Query:   "example.com",
				Server:  mockServer,
				Timeout: 10,
			},
			expectError:   false,
			expectContent: true,
		},
		{
			name: "validDomainDisableReferral",
			params: WhoisLookupParams{
				Query:           "example.com",
				Server:          mockServer,
				DisableReferral: true,
			},
			expectError:   false,
			expectContent: true,
		},
		{
			name: "validDomainParseJSON",
			params: WhoisLookupParams{
				Query:     "example.com",
				Server:    mockServer,
				ParseJSON: true,
			},
			expectError:   false,
			expectContent: true,
		},
		{
			name: "invalidDomain",
			params: WhoisLookupParams{
				Query:  "this-is-not-a-valid-domain-name-that-should-fail.invalid",
				Server: mockServer,
			},
			expectError:   false, // Mock server will return "No match" but not error
			expectContent: true,
		},
		{
			name: "emptyQuery",
			params: WhoisLookupParams{
				Query:  "",
				Server: mockServer,
			},
			expectError:   true,
			expectContent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			session := &mcp.ServerSession{}

			params := &mcp.CallToolParamsFor[WhoisLookupParams]{
				Arguments: tt.params,
			}

			result, err := handleWhoisLookup(ctx, session, params)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !tt.expectContent {
				return
			}

			if result == nil {
				t.Errorf("expected result but got nil")
				return
			}

			if len(result.Content) == 0 {
				t.Errorf("expected content but got empty slice")
				return
			}

			// Check that we have text content
			textContent, ok := result.Content[0].(*mcp.TextContent)
			if !ok {
				t.Errorf("expected TextContent but got %T", result.Content[0])
				return
			}

			if textContent.Text == "" {
				t.Errorf("expected non-empty text content")
			}

			// If ParseJSON was requested, verify it's valid JSON
			if tt.params.ParseJSON {
				var parsed interface{}
				if err := json.Unmarshal([]byte(textContent.Text), &parsed); err != nil {
					// If it fails to parse, it should contain an error message
					if !strings.Contains(textContent.Text, "Failed to parse whois data as JSON") {
						t.Errorf("expected JSON or parse error message, got: %s", textContent.Text)
					}
				}
			}
		})
	}
}

func TestHandleWhoisLookupWithMockServer(t *testing.T) {
	// Create a mock whois server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create mock server: %v", err)
	}
	defer func() { _ = listener.Close() }()

	// Get the port number
	address := listener.Addr().(*net.TCPAddr)
	port := fmt.Sprintf("%d", address.Port)

	// Start mock server
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer func() { _ = c.Close() }()

				// Read the query
				buffer := make([]byte, 1024)
				n, err := c.Read(buffer)
				if err != nil {
					return
				}

				query := strings.TrimSpace(string(buffer[:n]))

				// Send mock response
				mockResponse := fmt.Sprintf(`Domain Name: %s
Registry Domain ID: 123456
Registrar WHOIS Server: whois.example.com
Creation Date: 2024-01-01T00:00:00Z
Registry Expiry Date: 2025-01-01T00:00:00Z
`, strings.TrimSpace(query))

				_, _ = c.Write([]byte(mockResponse))
			}(conn)
		}
	}()

	// Test with mock server
	ctx := context.Background()
	session := &mcp.ServerSession{}

	params := &mcp.CallToolParamsFor[WhoisLookupParams]{
		Arguments: WhoisLookupParams{
			Query:  "test.example.com",
			Server: "127.0.0.1:" + port,
		},
	}

	result, err := handleWhoisLookup(ctx, session, params)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if result == nil || len(result.Content) == 0 {
		t.Errorf("expected result with content")
		return
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Errorf("expected TextContent")
		return
	}

	if !strings.Contains(textContent.Text, "test.example.com") {
		t.Errorf("expected response to contain query domain, got: %s", textContent.Text)
	}
}

func TestHandleWhoisLookupTimeout(t *testing.T) {
	// Create a mock server that never responds
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create mock server: %v", err)
	}
	defer func() { _ = listener.Close() }()

	address := listener.Addr().(*net.TCPAddr)
	port := fmt.Sprintf("%d", address.Port)

	// Start mock server that accepts but never responds
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			// Don't close connection, let it timeout
			go func(c net.Conn) {
				time.Sleep(5 * time.Second)
				_ = c.Close()
			}(conn)
		}
	}()

	ctx := context.Background()
	session := &mcp.ServerSession{}

	params := &mcp.CallToolParamsFor[WhoisLookupParams]{
		Arguments: WhoisLookupParams{
			Query:   "timeout.test.com",
			Server:  "127.0.0.1:" + port,
			Timeout: 1, // 1 second timeout
		},
	}

	result, err := handleWhoisLookup(ctx, session, params)
	if err == nil {
		t.Errorf("expected timeout error but got none")
		return
	}

	if result != nil {
		t.Errorf("expected nil result on error")
	}

	// Check that the error message mentions timeout or connection
	errMsg := err.Error()
	if !strings.Contains(errMsg, "timeout") && !strings.Contains(errMsg, "connection") {
		t.Errorf("expected timeout or connection error, got: %s", errMsg)
	}
}

func TestWhoisLookupParamsValidation(t *testing.T) {
	tests := []struct {
		name   string
		params WhoisLookupParams
		valid  bool
	}{
		{
			name: "validMinimal",
			params: WhoisLookupParams{
				Query: "example.com",
			},
			valid: true,
		},
		{
			name: "validWithAllOptions",
			params: WhoisLookupParams{
				Query:           "example.com",
				Server:          "whois.example.com",
				Timeout:         30,
				DisableReferral: true,
				ParseJSON:       true,
			},
			valid: true,
		},
		{
			name: "invalidEmptyQuery",
			params: WhoisLookupParams{
				Query: "",
			},
			valid: false,
		},
		{
			name: "validIPv4",
			params: WhoisLookupParams{
				Query: "8.8.8.8",
			},
			valid: true,
		},
		{
			name: "validIPv6",
			params: WhoisLookupParams{
				Query: "2001:4860:4860::8888",
			},
			valid: true,
		},
		{
			name: "validASN",
			params: WhoisLookupParams{
				Query: "AS15169",
			},
			valid: true,
		},
		{
			name: "validASNNumeric",
			params: WhoisLookupParams{
				Query: "15169",
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling/unmarshaling
			data, err := json.Marshal(tt.params)
			if err != nil {
				t.Errorf("failed to marshal params: %v", err)
				return
			}

			var unmarshaled WhoisLookupParams
			if err := json.Unmarshal(data, &unmarshaled); err != nil {
				t.Errorf("failed to unmarshal params: %v", err)
				return
			}

			// Basic validation - query should not be empty for valid cases
			if tt.valid && unmarshaled.Query == "" {
				t.Errorf("expected non-empty query for valid case")
			}

			// For our current test cases, invalid means empty query
			// This might need adjustment if we add more validation rules
			if !tt.valid && unmarshaled.Query != "" {
				t.Logf("Note: invalid test case with non-empty query: %s", unmarshaled.Query)
			}
		})
	}
}

func TestHandleWhoisLookupJSONParsing(t *testing.T) {
	// Create a mock server that returns parseable whois data
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create mock server: %v", err)
	}
	defer func() { _ = listener.Close() }()

	address := listener.Addr().(*net.TCPAddr)
	port := fmt.Sprintf("%d", address.Port)

	// Start mock server with parseable whois response
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer func() { _ = c.Close() }()

				// Read the query first
				buffer := make([]byte, 1024)
				_, err := c.Read(buffer)
				if err != nil {
					return
				}

				// Send a realistic whois response that should be parseable
				mockResponse := `Domain Name: EXAMPLE.COM
Registry Domain ID: 2336799_DOMAIN_COM-VRSN
Registrar WHOIS Server: whois.iana.org
Registrar URL: http://res-dom.iana.org
Updated Date: 2023-08-14T07:01:31Z
Creation Date: 1995-08-14T04:00:00Z
Registry Expiry Date: 2024-08-13T04:00:00Z
Registrar: RESERVED-Internet Assigned Numbers Authority
Registrar IANA ID: 376
Registrar Abuse Contact Email: abuse@iana.org
Registrar Abuse Contact Phone: +1.3103015200
Domain Status: clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
Domain Status: clientTransferProhibited https://icann.org/epp#clientTransferProhibited
Domain Status: clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
Name Server: A.IANA-SERVERS.NET
Name Server: B.IANA-SERVERS.NET
DNSSEC: signedDelegation
DNSSEC DS Data: 31589 8 1 3490A6806D47F17A34C29E2CE80E8A999FFBE4BE
DNSSEC DS Data: 31589 8 2 CDE0D742D6998AA554A92D890F8184C698CFAC8A26FA59875A990C03E576343C
URL of the ICANN Whois Inaccuracy Complaint Form: https://www.icann.org/wicf/
`
				_, _ = c.Write([]byte(mockResponse))
			}(conn)
		}
	}()

	// Test JSON parsing
	ctx := context.Background()
	session := &mcp.ServerSession{}

	params := &mcp.CallToolParamsFor[WhoisLookupParams]{
		Arguments: WhoisLookupParams{
			Query:     "example.com",
			Server:    "127.0.0.1:" + port,
			ParseJSON: true,
		},
	}

	result, err := handleWhoisLookup(ctx, session, params)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if result == nil || len(result.Content) == 0 {
		t.Errorf("expected result with content")
		return
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Errorf("expected TextContent")
		return
	}

	// Check if the response is valid JSON
	var parsed interface{}
	if err := json.Unmarshal([]byte(textContent.Text), &parsed); err != nil {
		// If parsing fails, check if it's a parse error message
		if !strings.Contains(textContent.Text, "Failed to parse whois data as JSON") {
			t.Errorf("expected valid JSON or parse error message, got parse error: %v", err)
		}
	} else {
		// If it parsed successfully, verify it contains expected fields
		parsedMap, ok := parsed.(map[string]interface{})
		if !ok {
			t.Errorf("expected parsed JSON to be an object")
			return
		}

		// Check for common whois fields that should be present
		if _, exists := parsedMap["domain"]; !exists {
			t.Errorf("expected parsed JSON to contain domain field")
		}
	}
}

// Benchmark tests
func BenchmarkHandleWhoisLookup(b *testing.B) {
	ctx := context.Background()
	session := &mcp.ServerSession{}

	params := &mcp.CallToolParamsFor[WhoisLookupParams]{
		Arguments: WhoisLookupParams{
			Query: "example.com",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = handleWhoisLookup(ctx, session, params)
	}
}

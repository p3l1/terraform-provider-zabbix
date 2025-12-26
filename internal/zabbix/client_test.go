// ABOUTME: Unit tests for the Zabbix API client using mock HTTP responses.
// ABOUTME: Tests cover successful requests, API errors, HTTP errors, and edge cases.

package zabbix

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("http://example.com/api", "test-token")

	if client.URL != "http://example.com/api" {
		t.Errorf("expected URL 'http://example.com/api', got '%s'", client.URL)
	}
	if client.Token != "test-token" {
		t.Errorf("expected Token 'test-token', got '%s'", client.Token)
	}
	if client.HTTPClient.Timeout != DefaultTimeout {
		t.Errorf("expected timeout %v, got %v", DefaultTimeout, client.HTTPClient.Timeout)
	}
}

func TestNewClientWithTimeout(t *testing.T) {
	timeout := 60 * time.Second
	client := NewClientWithTimeout("http://example.com/api", "test-token", timeout)

	if client.HTTPClient.Timeout != timeout {
		t.Errorf("expected timeout %v, got %v", timeout, client.HTTPClient.Timeout)
	}
}

func TestRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json-rpc" {
			t.Errorf("expected Content-Type 'application/json-rpc', got '%s'", r.Header.Get("Content-Type"))
		}

		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.JSONRPC != "2.0" {
			t.Errorf("expected jsonrpc '2.0', got '%s'", req.JSONRPC)
		}
		if req.Method != "host.get" {
			t.Errorf("expected method 'host.get', got '%s'", req.Method)
		}
		if req.Auth != "test-token" {
			t.Errorf("expected auth 'test-token', got '%s'", req.Auth)
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`[]`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	result, err := client.Request("host.get", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var hosts []interface{}
	if err := json.Unmarshal(result, &hosts); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}
}

func TestRequest_NoAuthForAPIInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.Auth != "" {
			t.Errorf("expected no auth for apiinfo.version, got '%s'", req.Auth)
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`"7.0.0"`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	result, err := client.Request("apiinfo.version", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var version string
	if err := json.Unmarshal(result, &version); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}
	if version != "7.0.0" {
		t.Errorf("expected version '7.0.0', got '%s'", version)
	}
}

func TestRequest_WithParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		params, ok := req.Params.(map[string]interface{})
		if !ok {
			t.Fatalf("expected params to be a map, got %T", req.Params)
		}
		if params["output"] != "extend" {
			t.Errorf("expected output 'extend', got '%v'", params["output"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`[{"hostid": "10084"}]`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	params := map[string]interface{}{
		"output": "extend",
	}
	result, err := client.Request("host.get", params)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var hosts []map[string]interface{}
	if err := json.Unmarshal(result, &hosts); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}
	if len(hosts) != 1 || hosts[0]["hostid"] != "10084" {
		t.Errorf("unexpected result: %v", hosts)
	}
}

func TestRequest_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := Response{
			JSONRPC: "2.0",
			Error: &Error{
				Code:    -32602,
				Message: "Invalid params.",
				Data:    "No permissions to referred object or it does not exist!",
			},
			ID: 1,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	_, err := client.Request("host.get", nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Method != "host.get" {
		t.Errorf("expected method 'host.get', got '%s'", apiErr.Method)
	}
	if apiErr.Err.Code != -32602 {
		t.Errorf("expected code -32602, got %d", apiErr.Err.Code)
	}
}

func TestRequest_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	_, err := client.Request("apiinfo.version", nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	httpErr, ok := err.(*HTTPError)
	if !ok {
		t.Fatalf("expected *HTTPError, got %T", err)
	}
	if httpErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status code 500, got %d", httpErr.StatusCode)
	}
}

func TestRequest_ConnectionError(t *testing.T) {
	client := NewClient("http://localhost:1", "test-token")
	_, err := client.Request("apiinfo.version", nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRequest_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("not json"))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	_, err := client.Request("apiinfo.version", nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRequest_IncrementingID(t *testing.T) {
	var receivedIDs []int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}
		receivedIDs = append(receivedIDs, req.ID)

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`"ok"`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	_, _ = client.Request("test", nil)
	_, _ = client.Request("test", nil)
	_, _ = client.Request("test", nil)

	if len(receivedIDs) != 3 {
		t.Fatalf("expected 3 requests, got %d", len(receivedIDs))
	}
	if receivedIDs[0] != 1 || receivedIDs[1] != 2 || receivedIDs[2] != 3 {
		t.Errorf("expected IDs [1, 2, 3], got %v", receivedIDs)
	}
}

func TestRequestWithContext_Cancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.RequestWithContext(ctx, "test", nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("expected context canceled error, got: %v", err)
	}
}

func TestRequest_ResponseIDMismatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`"ok"`),
			ID:      999,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	_, err := client.Request("test", nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "response id mismatch") {
		t.Errorf("expected response id mismatch error, got: %v", err)
	}
}

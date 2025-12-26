// ABOUTME: Unit tests for host group API methods using mock HTTP responses.
// ABOUTME: Tests cover CRUD operations and error handling for host groups.

package zabbix

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateHostGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.Method != "hostgroup.create" {
			t.Errorf("expected method 'hostgroup.create', got '%s'", req.Method)
		}

		params, ok := req.Params.(map[string]interface{})
		if !ok {
			t.Fatalf("expected params to be a map, got %T", req.Params)
		}
		if params["name"] != "Test Group" {
			t.Errorf("expected name 'Test Group', got '%v'", params["name"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"groupids": ["123"]}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	groupID, err := client.CreateHostGroup(context.Background(), "Test Group")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if groupID != "123" {
		t.Errorf("expected groupID '123', got '%s'", groupID)
	}
}

func TestCreateHostGroup_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		_ = json.Unmarshal(body, &req)

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"groupids": []}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	_, err := client.CreateHostGroup(context.Background(), "Test Group")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateHostGroup_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := Response{
			JSONRPC: "2.0",
			Error: &Error{
				Code:    -32602,
				Message: "Invalid params.",
				Data:    "Host group already exists.",
			},
			ID: 1,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	_, err := client.CreateHostGroup(context.Background(), "Test Group")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Method != "hostgroup.create" {
		t.Errorf("expected method 'hostgroup.create', got '%s'", apiErr.Method)
	}
}

func TestGetHostGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.Method != "hostgroup.get" {
			t.Errorf("expected method 'hostgroup.get', got '%s'", req.Method)
		}

		params, ok := req.Params.(map[string]interface{})
		if !ok {
			t.Fatalf("expected params to be a map, got %T", req.Params)
		}

		groupIDs, ok := params["groupids"].([]interface{})
		if !ok || len(groupIDs) != 1 || groupIDs[0] != "123" {
			t.Errorf("expected groupids ['123'], got '%v'", params["groupids"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`[{"groupid": "123", "name": "Test Group", "uuid": "abc-def-123"}]`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	group, err := client.GetHostGroup(context.Background(), "123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group == nil {
		t.Fatal("expected group, got nil")
	}
	if group.GroupID != "123" {
		t.Errorf("expected groupid '123', got '%s'", group.GroupID)
	}
	if group.Name != "Test Group" {
		t.Errorf("expected name 'Test Group', got '%s'", group.Name)
	}
	if group.UUID != "abc-def-123" {
		t.Errorf("expected uuid 'abc-def-123', got '%s'", group.UUID)
	}
}

func TestGetHostGroup_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		_ = json.Unmarshal(body, &req)

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`[]`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	group, err := client.GetHostGroup(context.Background(), "999")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group != nil {
		t.Errorf("expected nil group, got %v", group)
	}
}

func TestGetHostGroupByName_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.Method != "hostgroup.get" {
			t.Errorf("expected method 'hostgroup.get', got '%s'", req.Method)
		}

		params, ok := req.Params.(map[string]interface{})
		if !ok {
			t.Fatalf("expected params to be a map, got %T", req.Params)
		}

		filter, ok := params["filter"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected filter to be a map, got %T", params["filter"])
		}
		if filter["name"] != "Linux servers" {
			t.Errorf("expected filter name 'Linux servers', got '%v'", filter["name"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`[{"groupid": "2", "name": "Linux servers", "uuid": "xyz-123"}]`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	group, err := client.GetHostGroupByName(context.Background(), "Linux servers")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group == nil {
		t.Fatal("expected group, got nil")
	}
	if group.GroupID != "2" {
		t.Errorf("expected groupid '2', got '%s'", group.GroupID)
	}
	if group.Name != "Linux servers" {
		t.Errorf("expected name 'Linux servers', got '%s'", group.Name)
	}
}

func TestGetHostGroupByName_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		_ = json.Unmarshal(body, &req)

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`[]`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	group, err := client.GetHostGroupByName(context.Background(), "Nonexistent")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group != nil {
		t.Errorf("expected nil group, got %v", group)
	}
}

func TestUpdateHostGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.Method != "hostgroup.update" {
			t.Errorf("expected method 'hostgroup.update', got '%s'", req.Method)
		}

		params, ok := req.Params.(map[string]interface{})
		if !ok {
			t.Fatalf("expected params to be a map, got %T", req.Params)
		}
		if params["groupid"] != "123" {
			t.Errorf("expected groupid '123', got '%v'", params["groupid"])
		}
		if params["name"] != "Updated Group" {
			t.Errorf("expected name 'Updated Group', got '%v'", params["name"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"groupids": ["123"]}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	err := client.UpdateHostGroup(context.Background(), "123", "Updated Group")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateHostGroup_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		_ = json.Unmarshal(body, &req)

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"groupids": []}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	err := client.UpdateHostGroup(context.Background(), "123", "Updated Group")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeleteHostGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.Method != "hostgroup.delete" {
			t.Errorf("expected method 'hostgroup.delete', got '%s'", req.Method)
		}

		params, ok := req.Params.([]interface{})
		if !ok {
			t.Fatalf("expected params to be an array, got %T", req.Params)
		}
		if len(params) != 1 || params[0] != "123" {
			t.Errorf("expected params ['123'], got '%v'", params)
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"groupids": ["123"]}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	err := client.DeleteHostGroup(context.Background(), "123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteHostGroup_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		_ = json.Unmarshal(body, &req)

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"groupids": []}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	err := client.DeleteHostGroup(context.Background(), "123")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeleteHostGroup_APIError(t *testing.T) {
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
	err := client.DeleteHostGroup(context.Background(), "999")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Method != "hostgroup.delete" {
		t.Errorf("expected method 'hostgroup.delete', got '%s'", apiErr.Method)
	}
}

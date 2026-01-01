// ABOUTME: Unit tests for template group API methods using mock HTTP responses.
// ABOUTME: Tests cover CRUD operations and error handling for template groups.

package zabbix

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTemplateGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.Method != "templategroup.create" {
			t.Errorf("expected method 'templategroup.create', got '%s'", req.Method)
		}

		params, ok := req.Params.(map[string]interface{})
		if !ok {
			t.Fatalf("expected params to be a map, got %T", req.Params)
		}
		if params["name"] != "Test Templates" {
			t.Errorf("expected name 'Test Templates', got '%v'", params["name"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"groupids": ["100"]}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	groupID, err := client.CreateTemplateGroup(context.Background(), "Test Templates")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if groupID != "100" {
		t.Errorf("expected groupID '100', got '%s'", groupID)
	}
}

func TestCreateTemplateGroup_EmptyResponse(t *testing.T) {
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
	_, err := client.CreateTemplateGroup(context.Background(), "Test Templates")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateTemplateGroup_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := Response{
			JSONRPC: "2.0",
			Error: &Error{
				Code:    -32602,
				Message: "Invalid params.",
				Data:    "Template group already exists.",
			},
			ID: 1,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	_, err := client.CreateTemplateGroup(context.Background(), "Test Templates")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Method != "templategroup.create" {
		t.Errorf("expected method 'templategroup.create', got '%s'", apiErr.Method)
	}
}

func TestGetTemplateGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.Method != "templategroup.get" {
			t.Errorf("expected method 'templategroup.get', got '%s'", req.Method)
		}

		params, ok := req.Params.(map[string]interface{})
		if !ok {
			t.Fatalf("expected params to be a map, got %T", req.Params)
		}

		groupIDs, ok := params["groupids"].([]interface{})
		if !ok || len(groupIDs) != 1 || groupIDs[0] != "100" {
			t.Errorf("expected groupids ['100'], got '%v'", params["groupids"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`[{"groupid": "100", "name": "Test Templates", "uuid": "abc123"}]`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	group, err := client.GetTemplateGroup(context.Background(), "100")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group == nil {
		t.Fatal("expected group, got nil")
	}
	if group.GroupID != "100" {
		t.Errorf("expected groupid '100', got '%s'", group.GroupID)
	}
	if group.Name != "Test Templates" {
		t.Errorf("expected name 'Test Templates', got '%s'", group.Name)
	}
	if group.UUID != "abc123" {
		t.Errorf("expected uuid 'abc123', got '%s'", group.UUID)
	}
}

func TestGetTemplateGroup_NotFound(t *testing.T) {
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
	group, err := client.GetTemplateGroup(context.Background(), "99999")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group != nil {
		t.Errorf("expected nil group, got %v", group)
	}
}

func TestGetTemplateGroupByName_Success(t *testing.T) {
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

		filter, ok := params["filter"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected filter to be a map, got %T", params["filter"])
		}
		if filter["name"] != "Test Templates" {
			t.Errorf("expected filter name 'Test Templates', got '%v'", filter["name"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`[{"groupid": "100", "name": "Test Templates", "uuid": "abc123"}]`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	group, err := client.GetTemplateGroupByName(context.Background(), "Test Templates")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group == nil {
		t.Fatal("expected group, got nil")
	}
	if group.GroupID != "100" {
		t.Errorf("expected groupid '100', got '%s'", group.GroupID)
	}
}

func TestUpdateTemplateGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.Method != "templategroup.update" {
			t.Errorf("expected method 'templategroup.update', got '%s'", req.Method)
		}

		params, ok := req.Params.(map[string]interface{})
		if !ok {
			t.Fatalf("expected params to be a map, got %T", req.Params)
		}
		if params["groupid"] != "100" {
			t.Errorf("expected groupid '100', got '%v'", params["groupid"])
		}
		if params["name"] != "Updated Templates" {
			t.Errorf("expected name 'Updated Templates', got '%v'", params["name"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"groupids": ["100"]}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	err := client.UpdateTemplateGroup(context.Background(), "100", "Updated Templates")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateTemplateGroup_EmptyResponse(t *testing.T) {
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
	err := client.UpdateTemplateGroup(context.Background(), "100", "Updated Templates")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeleteTemplateGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.Method != "templategroup.delete" {
			t.Errorf("expected method 'templategroup.delete', got '%s'", req.Method)
		}

		params, ok := req.Params.([]interface{})
		if !ok {
			t.Fatalf("expected params to be an array, got %T", req.Params)
		}
		if len(params) != 1 || params[0] != "100" {
			t.Errorf("expected params ['100'], got '%v'", params)
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"groupids": ["100"]}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	err := client.DeleteTemplateGroup(context.Background(), "100")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteTemplateGroup_EmptyResponse(t *testing.T) {
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
	err := client.DeleteTemplateGroup(context.Background(), "100")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeleteTemplateGroup_APIError(t *testing.T) {
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
	err := client.DeleteTemplateGroup(context.Background(), "99999")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Method != "templategroup.delete" {
		t.Errorf("expected method 'templategroup.delete', got '%s'", apiErr.Method)
	}
}

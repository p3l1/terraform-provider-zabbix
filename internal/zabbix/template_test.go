// ABOUTME: Unit tests for template API methods using mock HTTP responses.
// ABOUTME: Tests cover CRUD operations and error handling for templates.

package zabbix

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTemplate_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.Method != "template.create" {
			t.Errorf("expected method 'template.create', got '%s'", req.Method)
		}

		params, ok := req.Params.(map[string]interface{})
		if !ok {
			t.Fatalf("expected params to be a map, got %T", req.Params)
		}
		if params["host"] != "my_template" {
			t.Errorf("expected host 'my_template', got '%v'", params["host"])
		}
		if params["name"] != "My Template" {
			t.Errorf("expected name 'My Template', got '%v'", params["name"])
		}

		groups, ok := params["groups"].([]interface{})
		if !ok || len(groups) != 1 {
			t.Fatalf("expected groups to be array with 1 element, got %v", params["groups"])
		}
		group := groups[0].(map[string]interface{})
		if group["groupid"] != "1" {
			t.Errorf("expected groupid '1', got '%v'", group["groupid"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"templateids": ["10001"]}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	template := &Template{
		Host:   "my_template",
		Name:   "My Template",
		Groups: []TemplateGroupID{{GroupID: "1"}},
	}
	templateID, err := client.CreateTemplate(context.Background(), template)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if templateID != "10001" {
		t.Errorf("expected templateID '10001', got '%s'", templateID)
	}
}

func TestCreateTemplate_WithDescription(t *testing.T) {
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
		if params["description"] != "Template description" {
			t.Errorf("expected description 'Template description', got '%v'", params["description"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"templateids": ["10001"]}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	template := &Template{
		Host:        "my_template",
		Name:        "My Template",
		Description: "Template description",
		Groups:      []TemplateGroupID{{GroupID: "1"}},
	}
	templateID, err := client.CreateTemplate(context.Background(), template)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if templateID != "10001" {
		t.Errorf("expected templateID '10001', got '%s'", templateID)
	}
}

func TestCreateTemplate_WithTags(t *testing.T) {
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

		tags, ok := params["tags"].([]interface{})
		if !ok || len(tags) != 1 {
			t.Fatalf("expected tags to be array with 1 element, got %v", params["tags"])
		}
		tag := tags[0].(map[string]interface{})
		if tag["tag"] != "environment" {
			t.Errorf("expected tag 'environment', got '%v'", tag["tag"])
		}
		if tag["value"] != "production" {
			t.Errorf("expected value 'production', got '%v'", tag["value"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"templateids": ["10001"]}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	template := &Template{
		Host:   "my_template",
		Groups: []TemplateGroupID{{GroupID: "1"}},
		Tags: []TemplateTag{{
			Tag:   "environment",
			Value: "production",
		}},
	}
	templateID, err := client.CreateTemplate(context.Background(), template)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if templateID != "10001" {
		t.Errorf("expected templateID '10001', got '%s'", templateID)
	}
}

func TestCreateTemplate_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		_ = json.Unmarshal(body, &req)

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"templateids": []}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	template := &Template{
		Host:   "my_template",
		Groups: []TemplateGroupID{{GroupID: "1"}},
	}
	_, err := client.CreateTemplate(context.Background(), template)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateTemplate_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := Response{
			JSONRPC: "2.0",
			Error: &Error{
				Code:    -32602,
				Message: "Invalid params.",
				Data:    "Template already exists.",
			},
			ID: 1,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	template := &Template{
		Host:   "my_template",
		Groups: []TemplateGroupID{{GroupID: "1"}},
	}
	_, err := client.CreateTemplate(context.Background(), template)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Method != "template.create" {
		t.Errorf("expected method 'template.create', got '%s'", apiErr.Method)
	}
}

func TestGetTemplate_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.Method != "template.get" {
			t.Errorf("expected method 'template.get', got '%s'", req.Method)
		}

		params, ok := req.Params.(map[string]interface{})
		if !ok {
			t.Fatalf("expected params to be a map, got %T", req.Params)
		}

		templateIDs, ok := params["templateids"].([]interface{})
		if !ok || len(templateIDs) != 1 || templateIDs[0] != "10001" {
			t.Errorf("expected templateids ['10001'], got '%v'", params["templateids"])
		}

		if params["selectGroups"] != "extend" {
			t.Errorf("expected selectGroups 'extend', got '%v'", params["selectGroups"])
		}
		if params["selectTags"] != "extend" {
			t.Errorf("expected selectTags 'extend', got '%v'", params["selectTags"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result: json.RawMessage(`[{
				"templateid": "10001",
				"host": "my_template",
				"name": "My Template",
				"description": "Template description",
				"uuid": "abc123",
				"groups": [{"groupid": "1", "name": "Templates"}],
				"tags": [{"tag": "environment", "value": "production"}]
			}]`),
			ID: req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	template, err := client.GetTemplate(context.Background(), "10001")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if template == nil {
		t.Fatal("expected template, got nil")
	}
	if template.TemplateID != "10001" {
		t.Errorf("expected templateid '10001', got '%s'", template.TemplateID)
	}
	if template.Host != "my_template" {
		t.Errorf("expected host 'my_template', got '%s'", template.Host)
	}
	if template.Name != "My Template" {
		t.Errorf("expected name 'My Template', got '%s'", template.Name)
	}
	if template.Description != "Template description" {
		t.Errorf("expected description 'Template description', got '%s'", template.Description)
	}
	if template.UUID != "abc123" {
		t.Errorf("expected uuid 'abc123', got '%s'", template.UUID)
	}
	if len(template.Groups) != 1 || template.Groups[0].GroupID != "1" {
		t.Errorf("expected groups with groupid '1', got %v", template.Groups)
	}
	if len(template.Tags) != 1 || template.Tags[0].Tag != "environment" {
		t.Errorf("expected tag 'environment', got %v", template.Tags)
	}
}

func TestGetTemplate_NotFound(t *testing.T) {
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
	template, err := client.GetTemplate(context.Background(), "99999")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if template != nil {
		t.Errorf("expected nil template, got %v", template)
	}
}

func TestGetTemplateByHost_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.Method != "template.get" {
			t.Errorf("expected method 'template.get', got '%s'", req.Method)
		}

		params, ok := req.Params.(map[string]interface{})
		if !ok {
			t.Fatalf("expected params to be a map, got %T", req.Params)
		}

		filter, ok := params["filter"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected filter to be a map, got %T", params["filter"])
		}
		if filter["host"] != "my_template" {
			t.Errorf("expected filter host 'my_template', got '%v'", filter["host"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result: json.RawMessage(`[{
				"templateid": "10001",
				"host": "my_template",
				"name": "My Template",
				"description": "",
				"uuid": "abc123",
				"groups": [{"groupid": "1", "name": "Templates"}],
				"tags": []
			}]`),
			ID: req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	template, err := client.GetTemplateByHost(context.Background(), "my_template")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if template == nil {
		t.Fatal("expected template, got nil")
	}
	if template.TemplateID != "10001" {
		t.Errorf("expected templateid '10001', got '%s'", template.TemplateID)
	}
	if template.Host != "my_template" {
		t.Errorf("expected host 'my_template', got '%s'", template.Host)
	}
}

func TestGetTemplateByHost_NotFound(t *testing.T) {
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
	template, err := client.GetTemplateByHost(context.Background(), "nonexistent")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if template != nil {
		t.Errorf("expected nil template, got %v", template)
	}
}

func TestUpdateTemplate_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.Method != "template.update" {
			t.Errorf("expected method 'template.update', got '%s'", req.Method)
		}

		params, ok := req.Params.(map[string]interface{})
		if !ok {
			t.Fatalf("expected params to be a map, got %T", req.Params)
		}
		if params["templateid"] != "10001" {
			t.Errorf("expected templateid '10001', got '%v'", params["templateid"])
		}
		if params["name"] != "Updated Template" {
			t.Errorf("expected name 'Updated Template', got '%v'", params["name"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"templateids": ["10001"]}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	template := &Template{
		TemplateID: "10001",
		Name:       "Updated Template",
	}
	err := client.UpdateTemplate(context.Background(), template)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateTemplate_WithGroups(t *testing.T) {
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

		groups, ok := params["groups"].([]interface{})
		if !ok || len(groups) != 2 {
			t.Fatalf("expected groups to be array with 2 elements, got %v", params["groups"])
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"templateids": ["10001"]}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	template := &Template{
		TemplateID: "10001",
		Groups: []TemplateGroupID{
			{GroupID: "1"},
			{GroupID: "2"},
		},
	}
	err := client.UpdateTemplate(context.Background(), template)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateTemplate_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		_ = json.Unmarshal(body, &req)

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"templateids": []}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	template := &Template{
		TemplateID: "10001",
		Name:       "Updated Template",
	}
	err := client.UpdateTemplate(context.Background(), template)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeleteTemplate_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		if req.Method != "template.delete" {
			t.Errorf("expected method 'template.delete', got '%s'", req.Method)
		}

		params, ok := req.Params.([]interface{})
		if !ok {
			t.Fatalf("expected params to be an array, got %T", req.Params)
		}
		if len(params) != 1 || params[0] != "10001" {
			t.Errorf("expected params ['10001'], got '%v'", params)
		}

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"templateids": ["10001"]}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	err := client.DeleteTemplate(context.Background(), "10001")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteTemplate_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req Request
		_ = json.Unmarshal(body, &req)

		resp := Response{
			JSONRPC: "2.0",
			Result:  json.RawMessage(`{"templateids": []}`),
			ID:      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	err := client.DeleteTemplate(context.Background(), "10001")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeleteTemplate_APIError(t *testing.T) {
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
	err := client.DeleteTemplate(context.Background(), "99999")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Method != "template.delete" {
		t.Errorf("expected method 'template.delete', got '%s'", apiErr.Method)
	}
}

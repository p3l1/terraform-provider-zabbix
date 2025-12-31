// ABOUTME: Provides API methods for managing Zabbix template groups.
// ABOUTME: Implements CRUD operations using the templategroup.* JSON-RPC methods.

package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
)

// TemplateGroup represents a Zabbix template group.
type TemplateGroup struct {
	GroupID string `json:"groupid,omitempty"`
	Name    string `json:"name"`
	UUID    string `json:"uuid,omitempty"`
}

// CreateTemplateGroupParams contains parameters for creating a template group.
type CreateTemplateGroupParams struct {
	Name string `json:"name"`
}

// CreateTemplateGroupResponse contains the response from templategroup.create.
type CreateTemplateGroupResponse struct {
	GroupIDs []string `json:"groupids"`
}

// GetTemplateGroupParams contains parameters for retrieving template groups.
type GetTemplateGroupParams struct {
	GroupIDs []string               `json:"groupids,omitempty"`
	Filter   map[string]interface{} `json:"filter,omitempty"`
	Output   interface{}            `json:"output,omitempty"`
}

// UpdateTemplateGroupParams contains parameters for updating a template group.
type UpdateTemplateGroupParams struct {
	GroupID string `json:"groupid"`
	Name    string `json:"name"`
}

// UpdateTemplateGroupResponse contains the response from templategroup.update.
type UpdateTemplateGroupResponse struct {
	GroupIDs []string `json:"groupids"`
}

// DeleteTemplateGroupResponse contains the response from templategroup.delete.
type DeleteTemplateGroupResponse struct {
	GroupIDs []string `json:"groupids"`
}

// CreateTemplateGroup creates a new template group and returns the created group ID.
func (c *Client) CreateTemplateGroup(ctx context.Context, name string) (string, error) {
	params := CreateTemplateGroupParams{
		Name: name,
	}

	result, err := c.RequestWithContext(ctx, "templategroup.create", params)
	if err != nil {
		return "", err
	}

	var resp CreateTemplateGroupResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to unmarshal templategroup.create response: %w", err)
	}

	if len(resp.GroupIDs) == 0 {
		return "", fmt.Errorf("templategroup.create returned no group IDs")
	}

	return resp.GroupIDs[0], nil
}

// GetTemplateGroup retrieves a template group by ID.
func (c *Client) GetTemplateGroup(ctx context.Context, groupID string) (*TemplateGroup, error) {
	params := GetTemplateGroupParams{
		GroupIDs: []string{groupID},
		Output:   "extend",
	}

	result, err := c.RequestWithContext(ctx, "templategroup.get", params)
	if err != nil {
		return nil, err
	}

	var groups []TemplateGroup
	if err := json.Unmarshal(result, &groups); err != nil {
		return nil, fmt.Errorf("failed to unmarshal templategroup.get response: %w", err)
	}

	if len(groups) == 0 {
		return nil, nil
	}

	return &groups[0], nil
}

// GetTemplateGroupByName retrieves a template group by name.
func (c *Client) GetTemplateGroupByName(ctx context.Context, name string) (*TemplateGroup, error) {
	params := GetTemplateGroupParams{
		Filter: map[string]interface{}{
			"name": name,
		},
		Output: "extend",
	}

	result, err := c.RequestWithContext(ctx, "templategroup.get", params)
	if err != nil {
		return nil, err
	}

	var groups []TemplateGroup
	if err := json.Unmarshal(result, &groups); err != nil {
		return nil, fmt.Errorf("failed to unmarshal templategroup.get response: %w", err)
	}

	if len(groups) == 0 {
		return nil, nil
	}

	return &groups[0], nil
}

// UpdateTemplateGroup updates a template group's name.
func (c *Client) UpdateTemplateGroup(ctx context.Context, groupID, name string) error {
	params := UpdateTemplateGroupParams{
		GroupID: groupID,
		Name:    name,
	}

	result, err := c.RequestWithContext(ctx, "templategroup.update", params)
	if err != nil {
		return err
	}

	var resp UpdateTemplateGroupResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return fmt.Errorf("failed to unmarshal templategroup.update response: %w", err)
	}

	if len(resp.GroupIDs) == 0 {
		return fmt.Errorf("templategroup.update returned no group IDs")
	}

	return nil
}

// DeleteTemplateGroup deletes a template group by ID.
func (c *Client) DeleteTemplateGroup(ctx context.Context, groupID string) error {
	params := []string{groupID}

	result, err := c.RequestWithContext(ctx, "templategroup.delete", params)
	if err != nil {
		return err
	}

	var resp DeleteTemplateGroupResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return fmt.Errorf("failed to unmarshal templategroup.delete response: %w", err)
	}

	if len(resp.GroupIDs) == 0 {
		return fmt.Errorf("templategroup.delete returned no group IDs")
	}

	return nil
}

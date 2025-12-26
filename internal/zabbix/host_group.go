// ABOUTME: Provides API methods for managing Zabbix host groups.
// ABOUTME: Implements CRUD operations using the hostgroup.* JSON-RPC methods.

package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
)

// HostGroup represents a Zabbix host group.
type HostGroup struct {
	GroupID string `json:"groupid,omitempty"`
	Name    string `json:"name"`
	UUID    string `json:"uuid,omitempty"`
}

// CreateHostGroupParams contains parameters for creating a host group.
type CreateHostGroupParams struct {
	Name string `json:"name"`
}

// CreateHostGroupResponse contains the response from hostgroup.create.
type CreateHostGroupResponse struct {
	GroupIDs []string `json:"groupids"`
}

// GetHostGroupParams contains parameters for retrieving host groups.
type GetHostGroupParams struct {
	GroupIDs []string               `json:"groupids,omitempty"`
	Filter   map[string]interface{} `json:"filter,omitempty"`
	Output   interface{}            `json:"output,omitempty"`
}

// UpdateHostGroupParams contains parameters for updating a host group.
type UpdateHostGroupParams struct {
	GroupID string `json:"groupid"`
	Name    string `json:"name"`
}

// UpdateHostGroupResponse contains the response from hostgroup.update.
type UpdateHostGroupResponse struct {
	GroupIDs []string `json:"groupids"`
}

// DeleteHostGroupResponse contains the response from hostgroup.delete.
type DeleteHostGroupResponse struct {
	GroupIDs []string `json:"groupids"`
}

// CreateHostGroup creates a new host group and returns the created group ID.
func (c *Client) CreateHostGroup(ctx context.Context, name string) (string, error) {
	params := CreateHostGroupParams{
		Name: name,
	}

	result, err := c.RequestWithContext(ctx, "hostgroup.create", params)
	if err != nil {
		return "", err
	}

	var resp CreateHostGroupResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to unmarshal hostgroup.create response: %w", err)
	}

	if len(resp.GroupIDs) == 0 {
		return "", fmt.Errorf("hostgroup.create returned no group IDs")
	}

	return resp.GroupIDs[0], nil
}

// GetHostGroup retrieves a host group by ID.
func (c *Client) GetHostGroup(ctx context.Context, groupID string) (*HostGroup, error) {
	params := GetHostGroupParams{
		GroupIDs: []string{groupID},
		Output:   "extend",
	}

	result, err := c.RequestWithContext(ctx, "hostgroup.get", params)
	if err != nil {
		return nil, err
	}

	var groups []HostGroup
	if err := json.Unmarshal(result, &groups); err != nil {
		return nil, fmt.Errorf("failed to unmarshal hostgroup.get response: %w", err)
	}

	if len(groups) == 0 {
		return nil, nil
	}

	return &groups[0], nil
}

// GetHostGroupByName retrieves a host group by name.
func (c *Client) GetHostGroupByName(ctx context.Context, name string) (*HostGroup, error) {
	params := GetHostGroupParams{
		Filter: map[string]interface{}{
			"name": name,
		},
		Output: "extend",
	}

	result, err := c.RequestWithContext(ctx, "hostgroup.get", params)
	if err != nil {
		return nil, err
	}

	var groups []HostGroup
	if err := json.Unmarshal(result, &groups); err != nil {
		return nil, fmt.Errorf("failed to unmarshal hostgroup.get response: %w", err)
	}

	if len(groups) == 0 {
		return nil, nil
	}

	return &groups[0], nil
}

// UpdateHostGroup updates a host group's name.
func (c *Client) UpdateHostGroup(ctx context.Context, groupID, name string) error {
	params := UpdateHostGroupParams{
		GroupID: groupID,
		Name:    name,
	}

	result, err := c.RequestWithContext(ctx, "hostgroup.update", params)
	if err != nil {
		return err
	}

	var resp UpdateHostGroupResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return fmt.Errorf("failed to unmarshal hostgroup.update response: %w", err)
	}

	if len(resp.GroupIDs) == 0 {
		return fmt.Errorf("hostgroup.update returned no group IDs")
	}

	return nil
}

// DeleteHostGroup deletes a host group by ID.
func (c *Client) DeleteHostGroup(ctx context.Context, groupID string) error {
	// hostgroup.delete takes an array of group IDs directly
	params := []string{groupID}

	result, err := c.RequestWithContext(ctx, "hostgroup.delete", params)
	if err != nil {
		return err
	}

	var resp DeleteHostGroupResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return fmt.Errorf("failed to unmarshal hostgroup.delete response: %w", err)
	}

	if len(resp.GroupIDs) == 0 {
		return fmt.Errorf("hostgroup.delete returned no group IDs")
	}

	return nil
}

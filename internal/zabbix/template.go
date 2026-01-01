// ABOUTME: Provides API methods for managing Zabbix templates.
// ABOUTME: Implements CRUD operations using the template.* JSON-RPC methods.

package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
)

// Template represents a Zabbix template.
type Template struct {
	TemplateID  string            `json:"templateid,omitempty"`
	Host        string            `json:"host,omitempty"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	UUID        string            `json:"uuid,omitempty"`
	Groups      []TemplateGroupID `json:"groups,omitempty"`
	Tags        []TemplateTag     `json:"tags,omitempty"`
}

// TemplateGroupID represents a template group reference by ID.
type TemplateGroupID struct {
	GroupID string `json:"groupid"`
	Name    string `json:"name,omitempty"`
}

// TemplateTag represents a template tag.
type TemplateTag struct {
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

// CreateTemplateResponse contains the response from template.create.
type CreateTemplateResponse struct {
	TemplateIDs []string `json:"templateids"`
}

// GetTemplateParams contains parameters for retrieving templates.
type GetTemplateParams struct {
	TemplateIDs  []string               `json:"templateids,omitempty"`
	Filter       map[string]interface{} `json:"filter,omitempty"`
	Output       interface{}            `json:"output,omitempty"`
	SelectGroups interface{}            `json:"selectGroups,omitempty"`
	SelectTags   interface{}            `json:"selectTags,omitempty"`
}

// UpdateTemplateResponse contains the response from template.update.
type UpdateTemplateResponse struct {
	TemplateIDs []string `json:"templateids"`
}

// DeleteTemplateResponse contains the response from template.delete.
type DeleteTemplateResponse struct {
	TemplateIDs []string `json:"templateids"`
}

// CreateTemplate creates a new template and returns the created template ID.
func (c *Client) CreateTemplate(ctx context.Context, template *Template) (string, error) {
	params := map[string]interface{}{
		"host": template.Host,
	}

	if template.Name != "" {
		params["name"] = template.Name
	}

	if template.Description != "" {
		params["description"] = template.Description
	}

	if len(template.Groups) > 0 {
		groups := make([]map[string]string, len(template.Groups))
		for i, g := range template.Groups {
			groups[i] = map[string]string{"groupid": g.GroupID}
		}
		params["groups"] = groups
	}

	if len(template.Tags) > 0 {
		tags := make([]map[string]string, len(template.Tags))
		for i, t := range template.Tags {
			tags[i] = map[string]string{"tag": t.Tag, "value": t.Value}
		}
		params["tags"] = tags
	}

	result, err := c.RequestWithContext(ctx, "template.create", params)
	if err != nil {
		return "", err
	}

	var resp CreateTemplateResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to unmarshal template.create response: %w", err)
	}

	if len(resp.TemplateIDs) == 0 {
		return "", fmt.Errorf("template.create returned no template IDs")
	}

	return resp.TemplateIDs[0], nil
}

// GetTemplate retrieves a template by ID with all related data.
func (c *Client) GetTemplate(ctx context.Context, templateID string) (*Template, error) {
	params := GetTemplateParams{
		TemplateIDs:  []string{templateID},
		Output:       "extend",
		SelectGroups: "extend",
		SelectTags:   "extend",
	}

	result, err := c.RequestWithContext(ctx, "template.get", params)
	if err != nil {
		return nil, err
	}

	var templates []Template
	if err := json.Unmarshal(result, &templates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal template.get response: %w", err)
	}

	if len(templates) == 0 {
		return nil, nil
	}

	return &templates[0], nil
}

// GetTemplateByHost retrieves a template by technical name.
func (c *Client) GetTemplateByHost(ctx context.Context, host string) (*Template, error) {
	params := GetTemplateParams{
		Filter: map[string]interface{}{
			"host": host,
		},
		Output:       "extend",
		SelectGroups: "extend",
		SelectTags:   "extend",
	}

	result, err := c.RequestWithContext(ctx, "template.get", params)
	if err != nil {
		return nil, err
	}

	var templates []Template
	if err := json.Unmarshal(result, &templates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal template.get response: %w", err)
	}

	if len(templates) == 0 {
		return nil, nil
	}

	return &templates[0], nil
}

// UpdateTemplate updates a template.
func (c *Client) UpdateTemplate(ctx context.Context, template *Template) error {
	params := map[string]interface{}{
		"templateid": template.TemplateID,
	}

	if template.Host != "" {
		params["host"] = template.Host
	}

	if template.Name != "" {
		params["name"] = template.Name
	}

	if template.Description != "" {
		params["description"] = template.Description
	}

	if len(template.Groups) > 0 {
		groups := make([]map[string]string, len(template.Groups))
		for i, g := range template.Groups {
			groups[i] = map[string]string{"groupid": g.GroupID}
		}
		params["groups"] = groups
	}

	if template.Tags != nil {
		tags := make([]map[string]string, len(template.Tags))
		for i, t := range template.Tags {
			tags[i] = map[string]string{"tag": t.Tag, "value": t.Value}
		}
		params["tags"] = tags
	}

	result, err := c.RequestWithContext(ctx, "template.update", params)
	if err != nil {
		return err
	}

	var resp UpdateTemplateResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return fmt.Errorf("failed to unmarshal template.update response: %w", err)
	}

	if len(resp.TemplateIDs) == 0 {
		return fmt.Errorf("template.update returned no template IDs")
	}

	return nil
}

// DeleteTemplate deletes a template by ID.
func (c *Client) DeleteTemplate(ctx context.Context, templateID string) error {
	params := []string{templateID}

	result, err := c.RequestWithContext(ctx, "template.delete", params)
	if err != nil {
		return err
	}

	var resp DeleteTemplateResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return fmt.Errorf("failed to unmarshal template.delete response: %w", err)
	}

	if len(resp.TemplateIDs) == 0 {
		return fmt.Errorf("template.delete returned no template IDs")
	}

	return nil
}

// ImportConfigurationParams contains parameters for configuration.import.
type ImportConfigurationParams struct {
	Format string                 `json:"format"`
	Source string                 `json:"source"`
	Rules  map[string]interface{} `json:"rules"`
}

// ImportConfiguration imports configuration from YAML/XML/JSON.
func (c *Client) ImportConfiguration(ctx context.Context, format, source string) error {
	params := ImportConfigurationParams{
		Format: format,
		Source: source,
		Rules: map[string]interface{}{
			"templates": map[string]interface{}{
				"createMissing":  true,
				"updateExisting": true,
			},
			"template_groups": map[string]interface{}{
				"createMissing": true,
			},
			"items": map[string]interface{}{
				"createMissing":  true,
				"updateExisting": true,
			},
			"triggers": map[string]interface{}{
				"createMissing":  true,
				"updateExisting": true,
			},
			"discoveryRules": map[string]interface{}{
				"createMissing":  true,
				"updateExisting": true,
			},
			"valueMaps": map[string]interface{}{
				"createMissing":  true,
				"updateExisting": true,
			},
		},
	}

	_, err := c.RequestWithContext(ctx, "configuration.import", params)
	return err
}

// ExportConfigurationParams contains parameters for configuration.export.
type ExportConfigurationParams struct {
	Format  string                 `json:"format"`
	Options map[string]interface{} `json:"options"`
}

// ExportConfiguration exports a template configuration as YAML/XML/JSON.
func (c *Client) ExportConfiguration(ctx context.Context, format string, templateIDs []string) (string, error) {
	params := ExportConfigurationParams{
		Format: format,
		Options: map[string]interface{}{
			"templates": templateIDs,
		},
	}

	result, err := c.RequestWithContext(ctx, "configuration.export", params)
	if err != nil {
		return "", err
	}

	var exported string
	if err := json.Unmarshal(result, &exported); err != nil {
		return "", fmt.Errorf("failed to unmarshal configuration.export response: %w", err)
	}

	return exported, nil
}

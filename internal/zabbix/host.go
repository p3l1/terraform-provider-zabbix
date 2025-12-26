// ABOUTME: Provides API methods for managing Zabbix hosts.
// ABOUTME: Implements CRUD operations using the host.* JSON-RPC methods.

package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// Host represents a Zabbix host.
type Host struct {
	HostID          string           `json:"hostid,omitempty"`
	Host            string           `json:"host,omitempty"`
	Name            string           `json:"name,omitempty"`
	Status          int              `json:"-"`
	Groups          []HostGroupID    `json:"groups,omitempty"`
	Interfaces      []HostInterface  `json:"interfaces,omitempty"`
	Tags            []HostTag        `json:"tags,omitempty"`
	Templates       []TemplateID     `json:"templates,omitempty"`
	ParentTemplates []ParentTemplate `json:"parentTemplates,omitempty"`
}

// hostJSON is used for JSON marshaling/unmarshaling with string status.
type hostJSON struct {
	HostID          string           `json:"hostid,omitempty"`
	Host            string           `json:"host,omitempty"`
	Name            string           `json:"name,omitempty"`
	Status          string           `json:"status,omitempty"`
	Groups          []HostGroupID    `json:"groups,omitempty"`
	Interfaces      []HostInterface  `json:"interfaces,omitempty"`
	Tags            []HostTag        `json:"tags,omitempty"`
	Templates       []TemplateID     `json:"templates,omitempty"`
	ParentTemplates []ParentTemplate `json:"parentTemplates,omitempty"`
}

// UnmarshalJSON handles Zabbix API returning numeric values as strings.
func (h *Host) UnmarshalJSON(data []byte) error {
	var hj hostJSON
	if err := json.Unmarshal(data, &hj); err != nil {
		return err
	}

	h.HostID = hj.HostID
	h.Host = hj.Host
	h.Name = hj.Name
	h.Groups = hj.Groups
	h.Interfaces = hj.Interfaces
	h.Tags = hj.Tags
	h.Templates = hj.Templates
	h.ParentTemplates = hj.ParentTemplates

	if hj.Status != "" {
		status, err := strconv.Atoi(hj.Status)
		if err != nil {
			return fmt.Errorf("invalid status value: %s", hj.Status)
		}
		h.Status = status
	}

	return nil
}

// HostGroupID represents a host group reference by ID.
type HostGroupID struct {
	GroupID string `json:"groupid"`
	Name    string `json:"name,omitempty"`
}

// HostInterface represents a host interface configuration.
type HostInterface struct {
	InterfaceID string `json:"interfaceid,omitempty"`
	Type        int    `json:"-"`
	Main        int    `json:"-"`
	UseIP       int    `json:"-"`
	IP          string `json:"ip"`
	DNS         string `json:"dns"`
	Port        string `json:"port"`
}

// hostInterfaceJSON is used for JSON unmarshaling with string numeric fields.
type hostInterfaceJSON struct {
	InterfaceID string `json:"interfaceid,omitempty"`
	Type        string `json:"type"`
	Main        string `json:"main"`
	UseIP       string `json:"useip"`
	IP          string `json:"ip"`
	DNS         string `json:"dns"`
	Port        string `json:"port"`
}

// UnmarshalJSON handles Zabbix API returning numeric values as strings.
func (hi *HostInterface) UnmarshalJSON(data []byte) error {
	var hij hostInterfaceJSON
	if err := json.Unmarshal(data, &hij); err != nil {
		return err
	}

	hi.InterfaceID = hij.InterfaceID
	hi.IP = hij.IP
	hi.DNS = hij.DNS
	hi.Port = hij.Port

	if hij.Type != "" {
		t, err := strconv.Atoi(hij.Type)
		if err != nil {
			return fmt.Errorf("invalid interface type value: %s", hij.Type)
		}
		hi.Type = t
	}

	if hij.Main != "" {
		m, err := strconv.Atoi(hij.Main)
		if err != nil {
			return fmt.Errorf("invalid interface main value: %s", hij.Main)
		}
		hi.Main = m
	}

	if hij.UseIP != "" {
		u, err := strconv.Atoi(hij.UseIP)
		if err != nil {
			return fmt.Errorf("invalid interface useip value: %s", hij.UseIP)
		}
		hi.UseIP = u
	}

	return nil
}

// MarshalJSON handles sending numeric values as integers to Zabbix API.
func (hi HostInterface) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"type":  hi.Type,
		"main":  hi.Main,
		"useip": hi.UseIP,
		"ip":    hi.IP,
		"dns":   hi.DNS,
		"port":  hi.Port,
	}
	if hi.InterfaceID != "" {
		m["interfaceid"] = hi.InterfaceID
	}
	return json.Marshal(m)
}

// HostTag represents a host tag.
type HostTag struct {
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

// TemplateID represents a template reference by ID.
type TemplateID struct {
	TemplateID string `json:"templateid"`
}

// ParentTemplate represents a linked template returned from host.get.
type ParentTemplate struct {
	TemplateID string `json:"templateid"`
	Name       string `json:"name,omitempty"`
}

// CreateHostResponse contains the response from host.create.
type CreateHostResponse struct {
	HostIDs []string `json:"hostids"`
}

// GetHostParams contains parameters for retrieving hosts.
type GetHostParams struct {
	HostIDs               []string               `json:"hostids,omitempty"`
	Filter                map[string]interface{} `json:"filter,omitempty"`
	Output                interface{}            `json:"output,omitempty"`
	SelectGroups          interface{}            `json:"selectGroups,omitempty"`
	SelectInterfaces      interface{}            `json:"selectInterfaces,omitempty"`
	SelectTags            interface{}            `json:"selectTags,omitempty"`
	SelectParentTemplates interface{}            `json:"selectParentTemplates,omitempty"`
}

// UpdateHostResponse contains the response from host.update.
type UpdateHostResponse struct {
	HostIDs []string `json:"hostids"`
}

// DeleteHostResponse contains the response from host.delete.
type DeleteHostResponse struct {
	HostIDs []string `json:"hostids"`
}

// CreateHost creates a new host and returns the created host ID.
func (c *Client) CreateHost(ctx context.Context, host *Host) (string, error) {
	params := map[string]interface{}{
		"host":   host.Host,
		"status": host.Status,
	}

	if host.Name != "" {
		params["name"] = host.Name
	}

	if len(host.Groups) > 0 {
		groups := make([]map[string]string, len(host.Groups))
		for i, g := range host.Groups {
			groups[i] = map[string]string{"groupid": g.GroupID}
		}
		params["groups"] = groups
	}

	if len(host.Interfaces) > 0 {
		interfaces := make([]map[string]interface{}, len(host.Interfaces))
		for i, iface := range host.Interfaces {
			interfaces[i] = map[string]interface{}{
				"type":  iface.Type,
				"main":  iface.Main,
				"useip": iface.UseIP,
				"ip":    iface.IP,
				"dns":   iface.DNS,
				"port":  iface.Port,
			}
		}
		params["interfaces"] = interfaces
	}

	if len(host.Templates) > 0 {
		templates := make([]map[string]string, len(host.Templates))
		for i, t := range host.Templates {
			templates[i] = map[string]string{"templateid": t.TemplateID}
		}
		params["templates"] = templates
	}

	if len(host.Tags) > 0 {
		tags := make([]map[string]string, len(host.Tags))
		for i, t := range host.Tags {
			tags[i] = map[string]string{"tag": t.Tag, "value": t.Value}
		}
		params["tags"] = tags
	}

	result, err := c.RequestWithContext(ctx, "host.create", params)
	if err != nil {
		return "", err
	}

	var resp CreateHostResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to unmarshal host.create response: %w", err)
	}

	if len(resp.HostIDs) == 0 {
		return "", fmt.Errorf("host.create returned no host IDs")
	}

	return resp.HostIDs[0], nil
}

// GetHost retrieves a host by ID with all related data.
func (c *Client) GetHost(ctx context.Context, hostID string) (*Host, error) {
	params := GetHostParams{
		HostIDs:               []string{hostID},
		Output:                "extend",
		SelectGroups:          "extend",
		SelectInterfaces:      "extend",
		SelectTags:            "extend",
		SelectParentTemplates: "extend",
	}

	result, err := c.RequestWithContext(ctx, "host.get", params)
	if err != nil {
		return nil, err
	}

	var hosts []Host
	if err := json.Unmarshal(result, &hosts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal host.get response: %w", err)
	}

	if len(hosts) == 0 {
		return nil, nil
	}

	return &hosts[0], nil
}

// GetHostByName retrieves a host by technical name.
func (c *Client) GetHostByName(ctx context.Context, hostname string) (*Host, error) {
	params := GetHostParams{
		Filter: map[string]interface{}{
			"host": hostname,
		},
		Output:                "extend",
		SelectGroups:          "extend",
		SelectInterfaces:      "extend",
		SelectTags:            "extend",
		SelectParentTemplates: "extend",
	}

	result, err := c.RequestWithContext(ctx, "host.get", params)
	if err != nil {
		return nil, err
	}

	var hosts []Host
	if err := json.Unmarshal(result, &hosts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal host.get response: %w", err)
	}

	if len(hosts) == 0 {
		return nil, nil
	}

	return &hosts[0], nil
}

// UpdateHost updates a host.
func (c *Client) UpdateHost(ctx context.Context, host *Host) error {
	params := map[string]interface{}{
		"hostid": host.HostID,
	}

	if host.Host != "" {
		params["host"] = host.Host
	}

	if host.Name != "" {
		params["name"] = host.Name
	}

	// Status is always included since 0 is a valid value
	params["status"] = host.Status

	if len(host.Groups) > 0 {
		groups := make([]map[string]string, len(host.Groups))
		for i, g := range host.Groups {
			groups[i] = map[string]string{"groupid": g.GroupID}
		}
		params["groups"] = groups
	}

	if len(host.Interfaces) > 0 {
		interfaces := make([]map[string]interface{}, len(host.Interfaces))
		for i, iface := range host.Interfaces {
			ifaceMap := map[string]interface{}{
				"type":  iface.Type,
				"main":  iface.Main,
				"useip": iface.UseIP,
				"ip":    iface.IP,
				"dns":   iface.DNS,
				"port":  iface.Port,
			}
			if iface.InterfaceID != "" {
				ifaceMap["interfaceid"] = iface.InterfaceID
			}
			interfaces[i] = ifaceMap
		}
		params["interfaces"] = interfaces
	}

	if len(host.Templates) > 0 {
		templates := make([]map[string]string, len(host.Templates))
		for i, t := range host.Templates {
			templates[i] = map[string]string{"templateid": t.TemplateID}
		}
		params["templates"] = templates
	}

	if host.Tags != nil {
		tags := make([]map[string]string, len(host.Tags))
		for i, t := range host.Tags {
			tags[i] = map[string]string{"tag": t.Tag, "value": t.Value}
		}
		params["tags"] = tags
	}

	result, err := c.RequestWithContext(ctx, "host.update", params)
	if err != nil {
		return err
	}

	var resp UpdateHostResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return fmt.Errorf("failed to unmarshal host.update response: %w", err)
	}

	if len(resp.HostIDs) == 0 {
		return fmt.Errorf("host.update returned no host IDs")
	}

	return nil
}

// DeleteHost deletes a host by ID.
func (c *Client) DeleteHost(ctx context.Context, hostID string) error {
	params := []string{hostID}

	result, err := c.RequestWithContext(ctx, "host.delete", params)
	if err != nil {
		return err
	}

	var resp DeleteHostResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return fmt.Errorf("failed to unmarshal host.delete response: %w", err)
	}

	if len(resp.HostIDs) == 0 {
		return fmt.Errorf("host.delete returned no host IDs")
	}

	return nil
}

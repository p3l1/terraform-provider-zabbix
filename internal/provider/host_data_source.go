// ABOUTME: Terraform data source for looking up existing Zabbix hosts.
// ABOUTME: Retrieves host information by technical name.

package provider

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/p3l1/zabbix-terraform/internal/zabbix"
)

var _ datasource.DataSource = &HostDataSource{}

// HostDataSource defines the data source implementation.
type HostDataSource struct {
	client *zabbix.Client
}

// HostDataSourceModel describes the data source data model.
type HostDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	Host       types.String `tfsdk:"host"`
	Name       types.String `tfsdk:"name"`
	Groups     types.List   `tfsdk:"groups"`
	Templates  types.List   `tfsdk:"templates"`
	Status     types.Int64  `tfsdk:"status"`
	Interfaces types.List   `tfsdk:"interfaces"`
	Tags       types.List   `tfsdk:"tags"`
}

// NewHostDataSource creates a new data source instance.
func NewHostDataSource() datasource.DataSource {
	return &HostDataSource{}
}

func (d *HostDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

func (d *HostDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Zabbix host by technical name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the host (hostid in Zabbix).",
				Computed:    true,
			},
			"host": schema.StringAttribute{
				Description: "Technical name of the host to look up.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Visible name of the host.",
				Computed:    true,
			},
			"groups": schema.ListAttribute{
				Description: "List of host group IDs the host belongs to.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"templates": schema.ListAttribute{
				Description: "List of template IDs linked to the host.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"status": schema.Int64Attribute{
				Description: "Status of the host. 0 = enabled, 1 = disabled.",
				Computed:    true,
			},
			"interfaces": schema.ListNestedAttribute{
				Description: "Host interfaces.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"interface_id": schema.StringAttribute{
							Description: "ID of the interface.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Interface type: agent, snmp, ipmi, or jmx.",
							Computed:    true,
						},
						"ip": schema.StringAttribute{
							Description: "IP address used by the interface.",
							Computed:    true,
						},
						"dns": schema.StringAttribute{
							Description: "DNS name used by the interface.",
							Computed:    true,
						},
						"port": schema.StringAttribute{
							Description: "Port number used by the interface.",
							Computed:    true,
						},
						"main": schema.BoolAttribute{
							Description: "Whether this is the default interface for the type.",
							Computed:    true,
						},
						"use_ip": schema.BoolAttribute{
							Description: "Whether to use IP address instead of DNS name.",
							Computed:    true,
						},
					},
				},
			},
			"tags": schema.ListNestedAttribute{
				Description: "Host tags.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tag": schema.StringAttribute{
							Description: "Tag name.",
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "Tag value.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *HostDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*zabbix.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *zabbix.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *HostDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data HostDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	host, err := d.client.GetHostByName(ctx, data.Host.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Host",
			fmt.Sprintf("Could not read host with name %q: %s", data.Host.ValueString(), err),
		)
		return
	}

	if host == nil {
		resp.Diagnostics.AddError(
			"Host Not Found",
			fmt.Sprintf("No host found with technical name %q.", data.Host.ValueString()),
		)
		return
	}

	diags := d.apiToModel(ctx, host, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// apiToModel converts the Zabbix API struct to Terraform model.
func (d *HostDataSource) apiToModel(ctx context.Context, host *zabbix.Host, data *HostDataSourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(host.HostID)
	data.Host = types.StringValue(host.Host)
	data.Name = types.StringValue(host.Name)
	data.Status = types.Int64Value(int64(host.Status))

	// Convert groups
	groupIDs := make([]attr.Value, len(host.Groups))
	for i, g := range host.Groups {
		groupIDs[i] = types.StringValue(g.GroupID)
	}
	groupsList, diagsGroups := types.ListValue(types.StringType, groupIDs)
	diags.Append(diagsGroups...)
	data.Groups = groupsList

	// Convert templates from parentTemplates
	if len(host.ParentTemplates) > 0 {
		templateIDs := make([]attr.Value, len(host.ParentTemplates))
		for i, t := range host.ParentTemplates {
			templateIDs[i] = types.StringValue(t.TemplateID)
		}
		templatesList, diagsTemplates := types.ListValue(types.StringType, templateIDs)
		diags.Append(diagsTemplates...)
		data.Templates = templatesList
	} else {
		data.Templates = types.ListNull(types.StringType)
	}

	// Convert interfaces - sort by interface_id for stable ordering
	sort.Slice(host.Interfaces, func(i, j int) bool {
		return host.Interfaces[i].InterfaceID < host.Interfaces[j].InterfaceID
	})
	interfaceType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"interface_id": types.StringType,
			"type":         types.StringType,
			"ip":           types.StringType,
			"dns":          types.StringType,
			"port":         types.StringType,
			"main":         types.BoolType,
			"use_ip":       types.BoolType,
		},
	}
	interfaceValues := make([]attr.Value, len(host.Interfaces))
	for i, iface := range host.Interfaces {
		obj, diagsIface := types.ObjectValue(interfaceType.AttrTypes, map[string]attr.Value{
			"interface_id": types.StringValue(iface.InterfaceID),
			"type":         types.StringValue(interfaceTypeToString(iface.Type)),
			"ip":           types.StringValue(iface.IP),
			"dns":          types.StringValue(iface.DNS),
			"port":         types.StringValue(iface.Port),
			"main":         types.BoolValue(iface.Main == 1),
			"use_ip":       types.BoolValue(iface.UseIP == 1),
		})
		diags.Append(diagsIface...)
		interfaceValues[i] = obj
	}
	interfacesList, diagsInterfaces := types.ListValue(interfaceType, interfaceValues)
	diags.Append(diagsInterfaces...)
	data.Interfaces = interfacesList

	// Convert tags
	tagType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"tag":   types.StringType,
			"value": types.StringType,
		},
	}
	if len(host.Tags) > 0 {
		tagValues := make([]attr.Value, len(host.Tags))
		for i, tag := range host.Tags {
			obj, diagsTag := types.ObjectValue(tagType.AttrTypes, map[string]attr.Value{
				"tag":   types.StringValue(tag.Tag),
				"value": types.StringValue(tag.Value),
			})
			diags.Append(diagsTag...)
			tagValues[i] = obj
		}
		tagsList, diagsTags := types.ListValue(tagType, tagValues)
		diags.Append(diagsTags...)
		data.Tags = tagsList
	} else {
		data.Tags = types.ListNull(tagType)
	}

	return diags
}

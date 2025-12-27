// ABOUTME: Terraform resource for managing Zabbix hosts.
// ABOUTME: Implements CRUD operations and import functionality with interfaces, templates, and tags.

package provider

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/p3l1/terraform-provider-zabbix/internal/zabbix"
)

var (
	_ resource.Resource                = &HostResource{}
	_ resource.ResourceWithImportState = &HostResource{}
)

// HostResource defines the resource implementation.
type HostResource struct {
	client *zabbix.Client
}

// HostResourceModel describes the resource data model.
type HostResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Host       types.String `tfsdk:"host"`
	Name       types.String `tfsdk:"name"`
	Groups     types.List   `tfsdk:"groups"`
	Templates  types.List   `tfsdk:"templates"`
	Status     types.Int64  `tfsdk:"status"`
	Interfaces types.List   `tfsdk:"interfaces"`
	Tags       types.List   `tfsdk:"tags"`
}

// HostInterfaceModel describes a host interface.
type HostInterfaceModel struct {
	InterfaceID types.String `tfsdk:"interface_id"`
	Type        types.String `tfsdk:"type"`
	IP          types.String `tfsdk:"ip"`
	DNS         types.String `tfsdk:"dns"`
	Port        types.String `tfsdk:"port"`
	Main        types.Bool   `tfsdk:"main"`
	UseIP       types.Bool   `tfsdk:"use_ip"`
}

// HostTagModel describes a host tag.
type HostTagModel struct {
	Tag   types.String `tfsdk:"tag"`
	Value types.String `tfsdk:"value"`
}

// NewHostResource creates a new resource instance.
func NewHostResource() resource.Resource {
	return &HostResource{}
}

func (r *HostResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

func (r *HostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Zabbix host.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the host (hostid in Zabbix).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"host": schema.StringAttribute{
				Description: "Technical name of the host.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Visible name of the host. Defaults to the host value if not set.",
				Optional:    true,
				Computed:    true,
			},
			"groups": schema.ListAttribute{
				Description: "List of host group IDs the host belongs to.",
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"templates": schema.ListAttribute{
				Description: "List of template IDs to link to the host.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"status": schema.Int64Attribute{
				Description: "Status of the host. 0 = enabled (default), 1 = disabled.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"interfaces": schema.ListNestedAttribute{
				Description: "Host interfaces for monitoring.",
				Required:    true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"interface_id": schema.StringAttribute{
							Description: "ID of the interface (computed by Zabbix).",
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"type": schema.StringAttribute{
							Description: "Interface type: agent, snmp, ipmi, or jmx.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("agent", "snmp", "ipmi", "jmx"),
							},
						},
						"ip": schema.StringAttribute{
							Description: "IP address used by the interface.",
							Required:    true,
						},
						"dns": schema.StringAttribute{
							Description: "DNS name used by the interface.",
							Optional:    true,
							Computed:    true,
						},
						"port": schema.StringAttribute{
							Description: "Port number used by the interface.",
							Required:    true,
						},
						"main": schema.BoolAttribute{
							Description: "Whether this is the default interface for the type.",
							Required:    true,
						},
						"use_ip": schema.BoolAttribute{
							Description: "Whether to use IP address instead of DNS name.",
							Required:    true,
						},
					},
				},
			},
			"tags": schema.ListNestedAttribute{
				Description: "Host tags.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tag": schema.StringAttribute{
							Description: "Tag name.",
							Required:    true,
						},
						"value": schema.StringAttribute{
							Description: "Tag value.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *HostResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*zabbix.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *zabbix.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *HostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data HostResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	host, diags := r.modelToAPI(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hostID, err := r.client.CreateHost(ctx, host)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Host",
			fmt.Sprintf("Could not create host: %s", err),
		)
		return
	}

	apiHost, err := r.client.GetHost(ctx, hostID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Host",
			fmt.Sprintf("Could not read host after creation: %s", err),
		)
		return
	}

	if apiHost == nil {
		resp.Diagnostics.AddError(
			"Error Reading Host",
			fmt.Sprintf("Host %s was created but could not be found", hostID),
		)
		return
	}

	diags = r.apiToModel(ctx, apiHost, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data HostResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	host, err := r.client.GetHost(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Host",
			fmt.Sprintf("Could not read host ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}

	if host == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	diags := r.apiToModel(ctx, host, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data HostResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state HostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	host, diags := r.modelToAPI(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	host.HostID = state.ID.ValueString()

	err := r.client.UpdateHost(ctx, host)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Host",
			fmt.Sprintf("Could not update host ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}

	apiHost, err := r.client.GetHost(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Host",
			fmt.Sprintf("Could not read host after update: %s", err),
		)
		return
	}

	if apiHost == nil {
		resp.Diagnostics.AddError(
			"Error Reading Host",
			fmt.Sprintf("Host %s was updated but could not be found", state.ID.ValueString()),
		)
		return
	}

	diags = r.apiToModel(ctx, apiHost, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data HostResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteHost(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Host",
			fmt.Sprintf("Could not delete host ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}
}

func (r *HostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// modelToAPI converts the Terraform model to Zabbix API struct.
func (r *HostResource) modelToAPI(ctx context.Context, data *HostResourceModel) (*zabbix.Host, diag.Diagnostics) {
	var diags diag.Diagnostics

	host := &zabbix.Host{
		Host:   data.Host.ValueString(),
		Name:   data.Name.ValueString(),
		Status: int(data.Status.ValueInt64()),
	}

	// Convert groups
	var groupIDs []string
	diags.Append(data.Groups.ElementsAs(ctx, &groupIDs, false)...)
	if diags.HasError() {
		return nil, diags
	}
	for _, id := range groupIDs {
		host.Groups = append(host.Groups, zabbix.HostGroupID{GroupID: id})
	}

	// Convert templates
	if !data.Templates.IsNull() {
		var templateIDs []string
		diags.Append(data.Templates.ElementsAs(ctx, &templateIDs, false)...)
		if diags.HasError() {
			return nil, diags
		}
		for _, id := range templateIDs {
			host.Templates = append(host.Templates, zabbix.TemplateID{TemplateID: id})
		}
	}

	// Convert interfaces
	var interfaces []HostInterfaceModel
	diags.Append(data.Interfaces.ElementsAs(ctx, &interfaces, false)...)
	if diags.HasError() {
		return nil, diags
	}
	for _, iface := range interfaces {
		apiIface := zabbix.HostInterface{
			Type:  interfaceTypeToInt(iface.Type.ValueString()),
			IP:    iface.IP.ValueString(),
			DNS:   iface.DNS.ValueString(),
			Port:  iface.Port.ValueString(),
			Main:  boolToInt(iface.Main.ValueBool()),
			UseIP: boolToInt(iface.UseIP.ValueBool()),
		}
		if !iface.InterfaceID.IsNull() && !iface.InterfaceID.IsUnknown() {
			apiIface.InterfaceID = iface.InterfaceID.ValueString()
		}
		host.Interfaces = append(host.Interfaces, apiIface)
	}

	// Convert tags
	if !data.Tags.IsNull() {
		var tags []HostTagModel
		diags.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if diags.HasError() {
			return nil, diags
		}
		for _, tag := range tags {
			host.Tags = append(host.Tags, zabbix.HostTag{
				Tag:   tag.Tag.ValueString(),
				Value: tag.Value.ValueString(),
			})
		}
	}

	return host, diags
}

// apiToModel converts the Zabbix API struct to Terraform model.
func (r *HostResource) apiToModel(ctx context.Context, host *zabbix.Host, data *HostResourceModel) diag.Diagnostics {
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
	groupsList, d := types.ListValue(types.StringType, groupIDs)
	diags.Append(d...)
	data.Groups = groupsList

	// Convert templates from parentTemplates
	if len(host.ParentTemplates) > 0 {
		templateIDs := make([]attr.Value, len(host.ParentTemplates))
		for i, t := range host.ParentTemplates {
			templateIDs[i] = types.StringValue(t.TemplateID)
		}
		templatesList, d := types.ListValue(types.StringType, templateIDs)
		diags.Append(d...)
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
		obj, d := types.ObjectValue(interfaceType.AttrTypes, map[string]attr.Value{
			"interface_id": types.StringValue(iface.InterfaceID),
			"type":         types.StringValue(interfaceTypeToString(iface.Type)),
			"ip":           types.StringValue(iface.IP),
			"dns":          types.StringValue(iface.DNS),
			"port":         types.StringValue(iface.Port),
			"main":         types.BoolValue(iface.Main == 1),
			"use_ip":       types.BoolValue(iface.UseIP == 1),
		})
		diags.Append(d...)
		interfaceValues[i] = obj
	}
	interfacesList, d := types.ListValue(interfaceType, interfaceValues)
	diags.Append(d...)
	data.Interfaces = interfacesList

	// Convert tags
	if len(host.Tags) > 0 {
		tagType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"tag":   types.StringType,
				"value": types.StringType,
			},
		}
		tagValues := make([]attr.Value, len(host.Tags))
		for i, tag := range host.Tags {
			obj, d := types.ObjectValue(tagType.AttrTypes, map[string]attr.Value{
				"tag":   types.StringValue(tag.Tag),
				"value": types.StringValue(tag.Value),
			})
			diags.Append(d...)
			tagValues[i] = obj
		}
		tagsList, d := types.ListValue(tagType, tagValues)
		diags.Append(d...)
		data.Tags = tagsList
	} else {
		tagType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"tag":   types.StringType,
				"value": types.StringType,
			},
		}
		data.Tags = types.ListNull(tagType)
	}

	return diags
}

// interfaceTypeToInt converts interface type string to Zabbix API integer.
func interfaceTypeToInt(t string) int {
	switch t {
	case "agent":
		return 1
	case "snmp":
		return 2
	case "ipmi":
		return 3
	case "jmx":
		return 4
	default:
		return 1
	}
}

// interfaceTypeToString converts Zabbix API integer to interface type string.
func interfaceTypeToString(t int) string {
	switch t {
	case 1:
		return "agent"
	case 2:
		return "snmp"
	case 3:
		return "ipmi"
	case 4:
		return "jmx"
	default:
		return "agent"
	}
}

// boolToInt converts bool to Zabbix API integer (0 or 1).
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

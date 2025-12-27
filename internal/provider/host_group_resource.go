// ABOUTME: Terraform resource for managing Zabbix host groups.
// ABOUTME: Implements CRUD operations and import functionality.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/p3l1/terraform-provider-zabbix/internal/zabbix"
)

var (
	_ resource.Resource                = &HostGroupResource{}
	_ resource.ResourceWithImportState = &HostGroupResource{}
)

// HostGroupResource defines the resource implementation.
type HostGroupResource struct {
	client *zabbix.Client
}

// HostGroupResourceModel describes the resource data model.
type HostGroupResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	UUID types.String `tfsdk:"uuid"`
}

// NewHostGroupResource creates a new resource instance.
func NewHostGroupResource() resource.Resource {
	return &HostGroupResource{}
}

func (r *HostGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host_group"
}

func (r *HostGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Zabbix host group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the host group (groupid in Zabbix).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the host group.",
				Required:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The universally unique identifier of the host group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *HostGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *HostGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data HostGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID, err := r.client.CreateHostGroup(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Host Group",
			fmt.Sprintf("Could not create host group: %s", err),
		)
		return
	}

	group, err := r.client.GetHostGroup(ctx, groupID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Host Group",
			fmt.Sprintf("Could not read host group after creation: %s", err),
		)
		return
	}

	data.ID = types.StringValue(group.GroupID)
	data.Name = types.StringValue(group.Name)
	data.UUID = types.StringValue(group.UUID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HostGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data HostGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.GetHostGroup(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Host Group",
			fmt.Sprintf("Could not read host group ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}

	if group == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data.ID = types.StringValue(group.GroupID)
	data.Name = types.StringValue(group.Name)
	data.UUID = types.StringValue(group.UUID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HostGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data HostGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state HostGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateHostGroup(ctx, state.ID.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Host Group",
			fmt.Sprintf("Could not update host group ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}

	group, err := r.client.GetHostGroup(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Host Group",
			fmt.Sprintf("Could not read host group after update: %s", err),
		)
		return
	}

	data.ID = types.StringValue(group.GroupID)
	data.Name = types.StringValue(group.Name)
	data.UUID = types.StringValue(group.UUID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HostGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data HostGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteHostGroup(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Host Group",
			fmt.Sprintf("Could not delete host group ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}
}

func (r *HostGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

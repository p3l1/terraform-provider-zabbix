// ABOUTME: Terraform resource for managing Zabbix template groups.
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
	_ resource.Resource                = &TemplateGroupResource{}
	_ resource.ResourceWithImportState = &TemplateGroupResource{}
)

// TemplateGroupResource defines the resource implementation.
type TemplateGroupResource struct {
	client *zabbix.Client
}

// TemplateGroupResourceModel describes the resource data model.
type TemplateGroupResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	UUID types.String `tfsdk:"uuid"`
}

// NewTemplateGroupResource creates a new resource instance.
func NewTemplateGroupResource() resource.Resource {
	return &TemplateGroupResource{}
}

func (r *TemplateGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template_group"
}

func (r *TemplateGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Zabbix template group. Template groups are used to organize templates.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the template group (groupid in Zabbix).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the template group.",
				Required:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The universally unique identifier of the template group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *TemplateGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TemplateGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TemplateGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID, err := r.client.CreateTemplateGroup(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Template Group",
			fmt.Sprintf("Could not create template group: %s", err),
		)
		return
	}

	group, err := r.client.GetTemplateGroup(ctx, groupID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Template Group",
			fmt.Sprintf("Could not read template group after creation: %s", err),
		)
		return
	}

	data.ID = types.StringValue(group.GroupID)
	data.Name = types.StringValue(group.Name)
	data.UUID = types.StringValue(group.UUID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TemplateGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TemplateGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.GetTemplateGroup(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Template Group",
			fmt.Sprintf("Could not read template group ID %s: %s", data.ID.ValueString(), err),
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

func (r *TemplateGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TemplateGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TemplateGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateTemplateGroup(ctx, state.ID.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Template Group",
			fmt.Sprintf("Could not update template group ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}

	group, err := r.client.GetTemplateGroup(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Template Group",
			fmt.Sprintf("Could not read template group after update: %s", err),
		)
		return
	}

	data.ID = types.StringValue(group.GroupID)
	data.Name = types.StringValue(group.Name)
	data.UUID = types.StringValue(group.UUID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TemplateGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TemplateGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTemplateGroup(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Template Group",
			fmt.Sprintf("Could not delete template group ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}
}

func (r *TemplateGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

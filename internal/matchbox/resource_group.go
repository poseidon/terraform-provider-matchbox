package matchbox

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	matchbox "github.com/poseidon/matchbox/matchbox/client"
	"github.com/poseidon/matchbox/matchbox/server/serverpb"
	"github.com/poseidon/matchbox/matchbox/storage/storagepb"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &GroupResource{}
var _ resource.ResourceWithImportState = &GroupResource{}

func NewGroupResource() resource.Resource {
	return &GroupResource{}
}

// GroupResource defines the resource implementation.
type GroupResource struct {
	client *matchbox.Client
}

// GroupResourceModel describes the resource data model.
type GroupResourceModel struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Profile  types.String `tfsdk:"profile"`
	Selector types.Map    `tfsdk:"selector"`
	Metadata types.Map    `tfsdk:"metadata"`
}

func (r *GroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *GroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Groups define selectors which match zero or more machines. Machine(s) matching a group will boot and provision according to the group's `Profile`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Example identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"profile": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"selector": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"metadata": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *GroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*matchbox.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *matchbox.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *GroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data GroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()

	selectors := map[string]string{}
	resp.Diagnostics.Append(data.Selector.ElementsAs(ctx, &selectors, false)...)
	metadata := map[string]string{}
	resp.Diagnostics.Append(data.Metadata.ElementsAs(ctx, &metadata, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	metadataInt := make(map[string]interface{}, len(metadata))
	for key, value := range metadata {
		metadataInt[key] = value
	}

	richGroup := &storagepb.RichGroup{
		Id:       name,
		Profile:  data.Profile.ValueString(),
		Selector: selectors,
		Metadata: metadataInt,
	}
	group, err := richGroup.ToGroup()
	if err != nil {
		resp.Diagnostics.AddError("Serializing Group Failed", err.Error())
		return
	}

	_, err = r.client.Groups.GroupPut(ctx, &serverpb.GroupPutRequest{
		Group: group,
	})
	if err != nil {
		resp.Diagnostics.AddError("Group Create Failed", err.Error())
	}

	data.Id = types.StringValue(group.GetId())

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created matchbox group")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	groupGetResponse, err := r.client.Groups.GroupGet(ctx, &serverpb.GroupGetRequest{
		Id: name,
	})
	if err != nil {
		// Resouurce doesn't exist anymore
		data.Id = types.StringNull()
		return
	}

	group := groupGetResponse.GetGroup()

	var diags diag.Diagnostics
	selector, diags := types.MapValueFrom(ctx, types.StringType, group.Selector)
	resp.Diagnostics.Append(diags...)

	var metadata map[string]string
	if len(group.Metadata) > 0 {
		if err := json.Unmarshal(group.Metadata, &metadata); err != nil {
			resp.Diagnostics.AddError("Metadata Unmarshal Failed", err.Error())
			return
		}
	}

	data.Profile = types.StringValue(group.Profile)
	data.Selector = selector
	data.Metadata, diags = types.MapValueFrom(ctx, types.StringType, metadata)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Unexpected Group Update",
		"Updating a group is not supported, and the provider should not have tried to. Please report this bug to the provider developers.",
	)
}

func (r *GroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data GroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	_, err := r.client.Groups.GroupDelete(ctx, &serverpb.GroupDeleteRequest{
		Id: name,
	})
	if err != nil {
		resp.Diagnostics.AddError("Group Delete Failed", err.Error())
	}

	data.Id = types.StringNull()
}

func (r *GroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

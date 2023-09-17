package matchbox

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	matchbox "github.com/poseidon/matchbox/matchbox/client"
	"github.com/poseidon/matchbox/matchbox/server/serverpb"
	"github.com/poseidon/matchbox/matchbox/storage/storagepb"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ProfileResource{}
var _ resource.ResourceWithImportState = &ProfileResource{}

func NewProfileResource() resource.Resource {
	return &ProfileResource{}
}

// ProfileResource defines the resource implementation.
type ProfileResource struct {
	client *matchbox.Client
}

// ProfileResourceModel describes the resource data model.
type ProfileResourceModel struct {
	Id                   types.String
	Name                 types.String `tfsdk:"name"`
	Kernel               types.String `tfsdk:"kernel"`
	Initrd               types.List   `tfsdk:"initrd"`
	Args                 types.List   `tfsdk:"args"`
	ContainerLinuxConfig types.String `tfsdk:"container_linux_config"`
	RawIgnition          types.String `tfsdk:"raw_ignition"`
	GenericConfig        types.String `tfsdk:"generic_config"`
}

func (r *ProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile"
}

func (r *ProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Profiles reference an Ignition config, Butane Config, Cloud-Config, and/or generic config by name and define network boot settings.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
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
			"kernel": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"initrd": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"args": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"container_linux_config": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"raw_ignition": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"generic_config": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *ProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a Profile and its associated configs. Partial
// creates do not modify state and can be retried safely.
func (r *ProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProfileResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := validateResourceProfile(&data); err != nil {
		resp.Diagnostics.AddError("Invalid Profile", err.Error())
		return
	}

	// Profile
	name := data.Name.ValueString()
	kernel := data.Kernel.ValueString()
	// NetBoot
	var initrds, args []string
	resp.Diagnostics.Append(data.Initrd.ElementsAs(ctx, initrds, false)...)
	resp.Diagnostics.Append(data.Args.ElementsAs(ctx, args, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Container Linux config / Ignition config
	clcName, clcContent := containerLinuxConfig(&data)
	genericName, genericContent := genericConfig(&data)

	profile := &storagepb.Profile{
		Id: name,
		Boot: &storagepb.NetBoot{
			Kernel: kernel,
			Initrd: initrds,
			Args:   args,
		},
		IgnitionId: clcName,
		GenericId:  genericName,
	}

	// Profile
	_, err := r.client.Profiles.ProfilePut(ctx, &serverpb.ProfilePutRequest{
		Profile: profile,
	})
	if err != nil {
		resp.Diagnostics.AddError("Profile Create Failed", err.Error())
		return
	}

	// Container Linux Config
	if clcContent != "" {
		_, err := r.client.Ignition.IgnitionPut(ctx, &serverpb.IgnitionPutRequest{
			Name:   name,
			Config: []byte(clcContent),
		})
		if err != nil {
			resp.Diagnostics.AddError("Ignition Config Create Failed", err.Error())
			return
		}
	}

	// Generic Config
	if genericContent != "" {
		_, err := r.client.Generic.GenericPut(ctx, &serverpb.GenericPutRequest{
			Name:   name,
			Config: []byte(genericContent),
		})
		if err != nil {
			resp.Diagnostics.AddError("Generic Config Create Failed", err.Error())
			return
		}
	}

	data.Id = types.StringValue(profile.GetId())

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created matchbox profile")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func validateResourceProfile(data *ProfileResourceModel) error {

	hasRAW := !data.RawIgnition.IsNull()
	hasCLC := !data.ContainerLinuxConfig.IsNull()
	if hasCLC && hasRAW {
		return errors.New("container_linux_config and raw_ignition are mutually exclusive")
	}
	return nil
}

func (r *ProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProfileResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Profile
	name := data.Name.ValueString()
	profileGetResponse, err := r.client.Profiles.ProfileGet(ctx, &serverpb.ProfileGetRequest{
		Id: name,
	})
	if err != nil {
		// resource doesn't exist or is corrupted and needs creating
		data.Id = types.StringNull()
	}

	profile := profileGetResponse.Profile
	data.Kernel = types.StringValue(profile.Boot.Kernel)

	var diags diag.Diagnostics
	data.Initrd, diags = types.ListValueFrom(ctx, types.StringType, profile.Boot.Initrd)
	resp.Diagnostics.Append(diags...)
	data.Args, diags = types.ListValueFrom(ctx, types.StringType, profile.Boot.Args)
	resp.Diagnostics.Append(diags...)

	if profile.IgnitionId != "" {
		ignition, err := r.client.Ignition.IgnitionGet(ctx, &serverpb.IgnitionGetRequest{
			Name: profile.IgnitionId,
		})
		if err != nil {
			// resource doesn't exist or is corrupted and needs creating
			data.Id = types.StringNull()
			return
		}
		// .ign and .ignition files indicate raw ignition,
		// see https://github.com/poseidon/matchbox/blob/d6bb21d5853e7af7c3c54b74537176caf5460482/matchbox/http/ignition.go#L18
		if strings.HasSuffix(profile.IgnitionId, ".ign") || strings.HasSuffix(profile.IgnitionId, ".ignition") {
			data.RawIgnition = types.StringValue(string(ignition.Config))
		} else {
			data.ContainerLinuxConfig = types.StringValue(string(ignition.Config))
		}
	}

	if profile.GenericId != "" {
		generic, err := r.client.Generic.GenericGet(ctx, &serverpb.GenericGetRequest{
			Name: profile.GenericId,
		})
		if err != nil {
			// resource doesn't exist or is corrupted and needs creating
			data.Id = types.StringNull()
			return
		}
		data.GenericConfig = types.StringValue(string(generic.Config))
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Unexpected Profile Update",
		"Updating a profile is not supported, and the provider should not have tried to. Please report this bug to the provider developers.",
	)
}

func (r *ProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProfileResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Profile
	name := data.Name.ValueString()
	_, err := r.client.Profiles.ProfileDelete(ctx, &serverpb.ProfileDeleteRequest{
		Id: name,
	})
	if err != nil {
		resp.Diagnostics.AddError("Profile Delete Failed", err.Error())
	}

	// Container Linux Config
	if name, content := containerLinuxConfig(&data); content != "" {
		_, err = r.client.Ignition.IgnitionDelete(ctx, &serverpb.IgnitionDeleteRequest{
			Name: name,
		})
		if err != nil {
			resp.Diagnostics.AddError("Profile Ignition Config Delete Failed", err.Error())
		}
	}

	// Generic Config
	if name, content := genericConfig(&data); content != "" {
		_, err = r.client.Generic.GenericDelete(ctx, &serverpb.GenericDeleteRequest{
			Name: name,
		})
		if err != nil {
			resp.Diagnostics.AddError("Profile Generic Config Delete Failed", err.Error())
		}
	}

	data.Id = types.StringNull()
}

func (r *ProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// resourceProfileDelete deletes a Profile and its associated configs. Partial
// deletes leave state unchanged and can be retried (deleting resources which
// no longer exist is a no-op).
func resourceProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*matchbox.Client)

	// Profile
	name := d.Get("name").(string)
	_, err := client.Profiles.ProfileDelete(ctx, &serverpb.ProfileDeleteRequest{
		Id: name,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	// Container Linux Config
	if name, content := containerLinuxConfig(d); content != "" {
		_, err = client.Ignition.IgnitionDelete(ctx, &serverpb.IgnitionDeleteRequest{
			Name: name,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Generic Config
	if name, content := genericConfig(d); content != "" {
		_, err = client.Generic.GenericDelete(ctx, &serverpb.GenericDeleteRequest{
			Name: name,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// resource can be destroyed in state
	d.SetId("")
	return diags
}

func containerLinuxConfig(data *ProfileResourceModel) (filename, config string) {
	// use profile name to generate Container Linux and Ignition filenames
	name := data.Name.ValueString()

	if !data.ContainerLinuxConfig.IsNull() {
		content := data.ContainerLinuxConfig.ValueString()
		return fmt.Sprintf("%s.yaml.tmpl", name), content
	}

	if !data.RawIgnition.IsNull() {
		content := data.RawIgnition.ValueString()
		return fmt.Sprintf("%s.ign", name), content
	}

	return
}

func genericConfig(data *ProfileResourceModel) (filename, config string) {
	// use profile name to generate generic config filename
	name := data.Name.ValueString()

	if !data.GenericConfig.IsNull() {
		content := data.GenericConfig.ValueString()
		return name, content
	}

	return
}

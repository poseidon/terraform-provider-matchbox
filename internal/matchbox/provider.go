package matchbox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure MatchboxProvider satisfies various provider interfaces.
var _ provider.Provider = &MatchboxProvider{}

// MatchboxProvider defines the provider implementation.
type MatchboxProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// MatchboxProviderModel describes the provider data model.
type MatchboxProviderModel struct {
	Endpoint   types.String `tfsdk:"endpoint"`
	ClientCert types.String `tfsdk:"client_cert"`
	ClientKey  types.String `tfsdk:"client_key"`
	CA         types.String `tfsdk:"ca"`
}

func (p *MatchboxProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "matchbox"
	resp.Version = p.version
}

func (p *MatchboxProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Required: true,
			},
			"client_cert": schema.StringAttribute{
				Required: true,
			},
			"client_key": schema.StringAttribute{
				Required: true,
			},
			"ca": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (p *MatchboxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data MatchboxProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := NewMatchboxClient(&data)
	if err != nil {
		resp.Diagnostics.AddError("oops", err.Error())
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *MatchboxProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resourceGroup,
		resourceProfile,
	}
}

func (p *MatchboxProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// New returns a Provider for Matchbox.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MatchboxProvider{
			version: version,
		}
	}
}

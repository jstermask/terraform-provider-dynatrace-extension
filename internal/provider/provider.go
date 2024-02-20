// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"terraform-provider-dynatrace-extension/internal/dynatrace"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &DynatraceExtensionProvider{}

// ScaffoldingProvider defines the provider implementation.
type DynatraceExtensionProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type DynatraceExtensionProviderModel struct {
	EnvironmentUrl types.String `tfsdk:"envUrl"`
	ApiToken       types.String `tfsdk:"apiToken"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DynatraceExtensionProvider{
			version: version,
		}
	}
}

func (p *DynatraceExtensionProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dynaext"
	resp.Version = p.version
}

func (p *DynatraceExtensionProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"envUrl": schema.StringAttribute{
				Optional: false,
			},
			"apiToken": schema.StringAttribute{
				Optional: false,
			},
		},
	}
}

func (p *DynatraceExtensionProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// first checking provider configuration values
	var config DynatraceExtensionProviderModel

	// reading configuration...
	diags := req.Config.Get(ctx, &config)

	// appending all diagnostics to the response diagnostics
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		// we already have an error here, just by reading configuration, so we return here
		return
	}

	// checking configuration values...
	if config.ApiToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("apiToken"), "Unknown API Token", "Provider cannot create the client because API Token is unknown")
	}

	if config.EnvironmentUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("envUrl"), "Unknown Environment URL", "Provider cannot create the client because Environment URL is unknown")
	}

	if resp.Diagnostics.HasError() {
		// we already have an error here, just by reading configuration, so we return here
		return
	}


	// reading environment variables and overriding with configurations if required
	envUrl := os.Getenv("DYNATRACE_ENV_URL")
	apiToken := os.Getenv("DYNATRACE_API_TOKEN")

	if(!config.ApiToken.IsNull()) {
		apiToken = config.ApiToken.ValueString()
	}

	if(!config.EnvironmentUrl.IsNull()) {
		envUrl = config.EnvironmentUrl.ValueString()
	}


	// checking values validity
	if envUrl == "" {
		resp.Diagnostics.AddAttributeError(path.Root("envUrl"), "Missing Environment URL", "Provider cannot create the client because Environment URL is missing")
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(path.Root("apiToken"), "Missing API Token", "Provider cannot create the client because API Token is missing")
	}

	if resp.Diagnostics.HasError() {
		// we already have an error here, just by reading configuration, so we return here
		return
	}

	client, err := dynatrace.NewClient(&envUrl, &apiToken)
	if(err != nil) {
		resp.Diagnostics.AddError("Unable to created Dynatrace API Client", "Unexpected error :" + err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

}

func (p *DynatraceExtensionProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (p *DynatraceExtensionProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

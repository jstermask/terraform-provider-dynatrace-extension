// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
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


func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DynatraceExtensionProvider{
			version: version,
		}
	}
}

func (p *DynatraceExtensionProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dynatraceextension"
	resp.Version = p.version
}

func (p *DynatraceExtensionProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *DynatraceExtensionProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *DynatraceExtensionProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (p *DynatraceExtensionProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}




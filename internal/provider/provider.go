// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ka2n/terraform-provider-n8ncloud/internal/client"
)

// Ensure N8nCloudProvider satisfies various provider interfaces.
var _ provider.Provider = &N8nCloudProvider{}
var _ provider.ProviderWithFunctions = &N8nCloudProvider{}
var _ provider.ProviderWithEphemeralResources = &N8nCloudProvider{}

// N8nCloudProvider defines the provider implementation.
type N8nCloudProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// N8nCloudProviderModel describes the provider data model.
type N8nCloudProviderModel struct {
	APIKey      types.String `tfsdk:"api_key"`
	InstanceURL types.String `tfsdk:"instance_url"`
	Timeout     types.Int64  `tfsdk:"timeout"`
}

func (p *N8nCloudProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "n8ncloud"
	resp.Version = p.version
}

func (p *N8nCloudProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The n8n Cloud provider enables Terraform to manage n8n Cloud resources.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The API key for n8n cloud authentication. Can also be set via N8N_API_KEY environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"instance_url": schema.StringAttribute{
				MarkdownDescription: "The URL of your n8n cloud instance (e.g., https://yourinstance.app.n8n.cloud). Can also be set via N8N_INSTANCE_URL environment variable.",
				Optional:            true,
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "The timeout for API requests in seconds. Defaults to 30.",
				Optional:            true,
			},
		},
	}
}

func (p *N8nCloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data N8nCloudProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check configuration values
	if data.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown n8n Cloud API Key",
			"The provider cannot create the n8n Cloud API client as there is an unknown configuration value for the n8n Cloud API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the N8N_API_KEY environment variable.",
		)
	}

	if data.InstanceURL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("instance_url"),
			"Unknown n8n Cloud Instance URL",
			"The provider cannot create the n8n Cloud API client as there is an unknown configuration value for the n8n Cloud instance URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the N8N_INSTANCE_URL environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	apiKey := os.Getenv("N8N_API_KEY")
	instanceURL := os.Getenv("N8N_INSTANCE_URL")
	timeout := int64(30)

	if !data.APIKey.IsNull() {
		apiKey = data.APIKey.ValueString()
	}

	if !data.InstanceURL.IsNull() {
		instanceURL = data.InstanceURL.ValueString()
	}

	if !data.Timeout.IsNull() {
		timeout = data.Timeout.ValueInt64()
	}

	// Validate configuration
	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing n8n Cloud API Key",
			"The provider cannot create the n8n Cloud API client as there is a missing or empty value for the n8n Cloud API key. "+
				"Set the api_key value in the configuration or use the N8N_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if instanceURL == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("instance_url"),
			"Missing n8n Cloud Instance URL",
			"The provider cannot create the n8n Cloud API client as there is a missing or empty value for the n8n Cloud instance URL. "+
				"Set the instance_url value in the configuration or use the N8N_INSTANCE_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create the API client
	clientConfig := &client.Config{
		BaseURL: instanceURL,
		APIKey:  apiKey,
		Timeout: time.Duration(timeout) * time.Second,
	}

	apiClient, err := client.NewClient(clientConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create n8n Cloud API Client",
			fmt.Sprintf("An unexpected error occurred when creating the n8n Cloud API client: %s", err),
		)
		return
	}

	// Make the n8n Cloud client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

func (p *N8nCloudProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserResource,
	}
}

func (p *N8nCloudProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		// No ephemeral resources for n8n cloud provider currently
	}
}

func (p *N8nCloudProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewUserDataSource,
	}
}

func (p *N8nCloudProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		// No functions for n8n cloud provider currently
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &N8nCloudProvider{
			version: version,
		}
	}
}

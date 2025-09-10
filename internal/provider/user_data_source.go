// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ka2n/terraform-provider-n8ncloud/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &UserDataSource{}

func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

// UserDataSource defines the data source implementation.
type UserDataSource struct {
	client *client.Client
}

// UserDataSourceModel describes the data source data model.
type UserDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	Email           types.String `tfsdk:"email"`
	Role            types.String `tfsdk:"role"`
	FirstName       types.String `tfsdk:"first_name"`
	LastName        types.String `tfsdk:"last_name"`
	IsPending       types.Bool   `tfsdk:"is_pending"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
	InviteAcceptURL types.String `tfsdk:"invite_accept_url"`
}

func (d *UserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *UserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "User data source for querying existing n8n cloud users. You must specify either `id` or `email` to identify the user.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the user. Either id or email must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email address of the user. Either id or email must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The role of the user",
				Computed:            true,
			},
			"first_name": schema.StringAttribute{
				MarkdownDescription: "The first name of the user",
				Computed:            true,
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: "The last name of the user",
				Computed:            true,
			},
			"is_pending": schema.BoolAttribute{
				MarkdownDescription: "Whether the user has not yet set up their account",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the user was created",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the user was last updated",
				Computed:            true,
			},
			"invite_accept_url": schema.StringAttribute{
				MarkdownDescription: "The URL for the user to accept their invitation",
				Computed:            true,
			},
		},
	}
}

func (d *UserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that either ID or email is specified
	if data.ID.IsNull() && data.Email.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Attribute",
			"Either 'id' or 'email' must be specified",
		)
		return
	}

	var user *client.User
	var err error

	// Query by ID or email
	if !data.ID.IsNull() {
		user, err = d.client.GetUser(ctx, data.ID.ValueString())
	} else {
		user, err = d.client.GetUserByEmail(ctx, data.Email.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read user, got error: %s", err))
		return
	}

	// Map response body to model
	data.ID = types.StringValue(user.ID)
	data.Email = types.StringValue(user.Email)
	data.IsPending = types.BoolValue(user.IsPending)
	data.CreatedAt = types.StringValue(user.CreatedAt.Format(time.RFC3339))
	data.UpdatedAt = types.StringValue(user.UpdatedAt.Format(time.RFC3339))

	// Set role from API response
	if user.Role != "" {
		data.Role = types.StringValue(user.Role)
	} else {
		data.Role = types.StringNull()
	}

	if user.FirstName != nil {
		data.FirstName = types.StringValue(*user.FirstName)
	} else {
		data.FirstName = types.StringNull()
	}

	if user.LastName != nil {
		data.LastName = types.StringValue(*user.LastName)
	} else {
		data.LastName = types.StringNull()
	}

	if user.InviteAcceptUrl != "" {
		data.InviteAcceptURL = types.StringValue(user.InviteAcceptUrl)
	} else {
		data.InviteAcceptURL = types.StringNull()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

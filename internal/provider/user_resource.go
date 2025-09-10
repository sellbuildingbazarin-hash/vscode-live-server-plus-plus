// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ka2n/terraform-provider-n8ncloud/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &UserResource{}
var _ resource.ResourceWithImportState = &UserResource{}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

// UserResource defines the resource implementation.
type UserResource struct {
	client *client.Client
}

// UserResourceModel describes the resource data model.
type UserResourceModel struct {
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

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "User resource for managing n8n cloud users. Users can be imported using their email address: `terraform import n8ncloud_user.example user@example.com`",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the user",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email address of the user",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The role of the user (global:admin or global:member)",
				Required:            true,
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
				MarkdownDescription: "Whether the user has not yet set up their account. This value is managed externally and will change when the user accepts their invitation.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the user was created",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the user was last updated. This value is updated externally when the user's information changes.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"invite_accept_url": schema.StringAttribute{
				MarkdownDescription: "The URL for the user to accept their invitation",
				Computed:            true,
			},
		},
	}
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create the user
	createReq := &client.CreateUserRequest{
		Email: data.Email.ValueString(),
		Role:  data.Role.ValueString(),
	}

	tflog.Debug(ctx, "Creating n8n cloud user", map[string]interface{}{
		"email": createReq.Email,
		"role":  createReq.Role,
	})

	user, err := r.client.CreateUser(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create user, got error: %s", err))
		return
	}

	// Map response body to schema and populate computed attributes
	data.ID = types.StringValue(user.ID)
	data.IsPending = types.BoolValue(user.IsPending)
	data.CreatedAt = types.StringValue(user.CreatedAt.Format(time.RFC3339))
	data.UpdatedAt = types.StringValue(user.UpdatedAt.Format(time.RFC3339))

	// Set role from API response
	if user.Role != "" {
		data.Role = types.StringValue(user.Role)
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

	tflog.Trace(ctx, "Created n8n cloud user resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get fresh user data from API
	user, err := r.client.GetUser(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read user, got error: %s", err))
		return
	}

	// Update the model with the latest data
	data.Email = types.StringValue(user.Email)
	data.IsPending = types.BoolValue(user.IsPending)
	data.CreatedAt = types.StringValue(user.CreatedAt.Format(time.RFC3339))
	data.UpdatedAt = types.StringValue(user.UpdatedAt.Format(time.RFC3339))

	// Set role from API response
	if user.Role != "" {
		data.Role = types.StringValue(user.Role)
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data UserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update user role (only field that can be updated)
	err := r.client.UpdateUserRole(ctx, data.ID.ValueString(), data.Role.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update user role, got error: %s", err))
		return
	}

	// Get updated user data
	user, err := r.client.GetUser(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read updated user, got error: %s", err))
		return
	}

	// Update the model with the latest data
	data.UpdatedAt = types.StringValue(user.UpdatedAt.Format(time.RFC3339))

	tflog.Trace(ctx, "Updated n8n cloud user resource")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data UserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUser(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete user, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Deleted n8n cloud user resource")
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use email as import ID
	email := req.ID

	// Get user directly by email (API supports email as identifier)
	user, err := r.client.GetUser(ctx, email)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get user by email %s, got error: %s", email, err))
		return
	}

	// Set the resource ID and email
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), user.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("email"), user.Email)...)
}

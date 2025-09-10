// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"time"
)

// User represents an n8n cloud user.
type User struct {
	ID              string    `json:"id"`
	Email           string    `json:"email"`
	FirstName       *string   `json:"firstName,omitempty"`
	LastName        *string   `json:"lastName,omitempty"`
	IsPending       bool      `json:"isPending"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Role            string    `json:"role,omitempty"` // Role as string: "global:admin" or "global:member"
	InviteAcceptUrl string    `json:"inviteAcceptUrl,omitempty"`
}

// CreateUserRequest represents the request to create a new user.
type CreateUserRequest struct {
	Email     string `json:"email"`
	Role      string `json:"role,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

// UpdateUserRoleRequest represents the request to update a user's role.
type UpdateUserRoleRequest struct {
	NewRoleName string `json:"newRoleName"`
}

// UsersResponse represents the response from the list users endpoint.
type UsersResponse struct {
	Data       []User  `json:"data"`
	NextCursor *string `json:"nextCursor"`
}

// ErrorResponse represents an error response from the API.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Hint    string `json:"hint,omitempty"`
}

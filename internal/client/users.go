// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ListUsers retrieves all users from the n8n instance.
func (c *Client) ListUsers(ctx context.Context) ([]User, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/users?includeRole=true", nil)
	if err != nil {
		return nil, err
	}

	var resp UsersResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal users response: %w", err)
	}

	return resp.Data, nil
}

// GetUser retrieves a user by ID with role information.
func (c *Client) GetUser(ctx context.Context, id string) (*User, error) {
	path := fmt.Sprintf("/users/%s?includeRole=true", id)
	body, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// API returns user object directly, not wrapped in response
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user response: %w", err)
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email.
func (c *Client) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	// The API uses email as the identifier in the URL
	return c.GetUser(ctx, email)
}

// CreateUser creates a new user.
func (c *Client) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/users", req)
	if err != nil {
		return nil, err
	}

	// API returns user object directly, not wrapped in response
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create user response: %w", err)
	}

	return &user, nil
}

// UpdateUserRole updates a user's role.
func (c *Client) UpdateUserRole(ctx context.Context, id string, newRole string) error {
	path := fmt.Sprintf("/users/%s/role", id)
	req := &UpdateUserRoleRequest{
		NewRoleName: newRole,
	}

	_, err := c.doRequest(ctx, http.MethodPatch, path, req)
	return err
}

// DeleteUser deletes a user.
func (c *Client) DeleteUser(ctx context.Context, id string) error {
	path := fmt.Sprintf("/users/%s", id)
	_, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	return err
}

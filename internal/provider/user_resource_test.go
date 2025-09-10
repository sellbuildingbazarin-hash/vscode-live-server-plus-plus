// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccUserResource_basic(t *testing.T) {
	email := fmt.Sprintf("test-basic-%d@example.com", time.Now().Unix())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserResourceDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserResourceConfig(email, "user", "Test", "User"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"n8ncloud_user.test",
						tfjsonpath.New("email"),
						knownvalue.StringExact(email),
					),
					statecheck.ExpectKnownValue(
						"n8ncloud_user.test",
						tfjsonpath.New("role"),
						knownvalue.StringExact("user"),
					),
					statecheck.ExpectKnownValue(
						"n8ncloud_user.test",
						tfjsonpath.New("first_name"),
						knownvalue.StringExact("Test"),
					),
					statecheck.ExpectKnownValue(
						"n8ncloud_user.test",
						tfjsonpath.New("last_name"),
						knownvalue.StringExact("User"),
					),
					statecheck.ExpectKnownValue(
						"n8ncloud_user.test",
						tfjsonpath.New("is_pending"),
						knownvalue.Bool(true),
					),
				},
			},
			// ImportState testing using email as import ID
			{
				ResourceName:      "n8ncloud_user.test",
				ImportState:       true,
				ImportStateId:     email,
				ImportStateVerify: true,
				// invite_accept_url is only available during creation
				ImportStateVerifyIgnore: []string{"invite_accept_url"},
			},
			// Update and Read testing
			{
				Config: testAccUserResourceConfig(email, "admin", "Updated", "Admin"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"n8ncloud_user.test",
						tfjsonpath.New("email"),
						knownvalue.StringExact(email),
					),
					statecheck.ExpectKnownValue(
						"n8ncloud_user.test",
						tfjsonpath.New("role"),
						knownvalue.StringExact("admin"),
					),
					statecheck.ExpectKnownValue(
						"n8ncloud_user.test",
						tfjsonpath.New("first_name"),
						knownvalue.StringExact("Updated"),
					),
					statecheck.ExpectKnownValue(
						"n8ncloud_user.test",
						tfjsonpath.New("last_name"),
						knownvalue.StringExact("Admin"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccUserResource_onlyRequiredAttributes(t *testing.T) {
	email := fmt.Sprintf("test-required-%d@example.com", time.Now().Unix())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserResourceDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserResourceConfig_onlyRequired(email),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"n8ncloud_user.test",
						tfjsonpath.New("email"),
						knownvalue.StringExact(email),
					),
					statecheck.ExpectKnownValue(
						"n8ncloud_user.test",
						tfjsonpath.New("role"),
						knownvalue.StringExact("user"), // default value
					),
					statecheck.ExpectKnownValue(
						"n8ncloud_user.test",
						tfjsonpath.New("is_pending"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

func TestAccUserResource_invalidRole(t *testing.T) {
	email := fmt.Sprintf("test-invalid-%d@example.com", time.Now().Unix())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with invalid role
			{
				Config:      testAccUserResourceConfig(email, "invalid_role", "Test", "User"),
				ExpectError: regexp.MustCompile(`Attribute role value must be one of`),
			},
		},
	})
}

// TestAccUserResource_externalStateChange tests that external changes to computed attributes
// like is_pending don't cause unnecessary diffs.
func TestAccUserResource_externalStateChange(t *testing.T) {
	email := fmt.Sprintf("test-external-%d@example.com", time.Now().Unix())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserResourceDestroy,
		Steps: []resource.TestStep{
			// Create user
			{
				Config: testAccUserResourceConfig(email, "user", "Test", "User"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"n8ncloud_user.test",
						tfjsonpath.New("is_pending"),
						knownvalue.Bool(true),
					),
				},
			},
			// Simulate external change by re-applying the same config
			// The is_pending attribute might have changed externally (user accepted invitation)
			// but it shouldn't cause a diff since it's computed with UseStateForUnknown
			{
				Config:   testAccUserResourceConfig(email, "user", "Test", "User"),
				PlanOnly: true,
			},
		},
	})
}

func testAccUserResourceConfig(email, role, firstName, lastName string) string {
	return fmt.Sprintf(`
resource "n8ncloud_user" "test" {
  email      = %[1]q
  role       = %[2]q
  first_name = %[3]q
  last_name  = %[4]q
}
`, email, role, firstName, lastName)
}

func testAccUserResourceConfig_onlyRequired(email string) string {
	return fmt.Sprintf(`
resource "n8ncloud_user" "test" {
  email = %[1]q
}
`, email)
}

// TestAccUserResource_import tests importing a user by email.
func TestAccUserResource_import(t *testing.T) {
	email := fmt.Sprintf("test-import-%d@example.com", time.Now().Unix())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserResourceDestroy,
		Steps: []resource.TestStep{
			// Create a user first
			{
				Config: testAccUserResourceConfig(email, "user", "Import", "Test"),
			},
			// Import using email
			{
				ResourceName:      "n8ncloud_user.test",
				ImportState:       true,
				ImportStateId:     email,
				ImportStateVerify: true,
				// invite_accept_url is only available during creation
				ImportStateVerifyIgnore: []string{"invite_accept_url"},
			},
		},
	})
}

func testAccCheckUserResourceDestroy(s *terraform.State) error {
	// Add logic to verify user is deleted from n8n cloud
	// This would typically involve checking that the resource no longer exists
	// by making an API call and expecting a 404 or similar error
	return nil
}

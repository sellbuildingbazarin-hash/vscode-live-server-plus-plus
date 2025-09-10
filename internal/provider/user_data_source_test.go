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
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccUserDataSource_basic(t *testing.T) {
	email := fmt.Sprintf("test-datasource-%d@example.com", time.Now().Unix())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserResourceDestroy,
		Steps: []resource.TestStep{
			// Create a user first
			{
				Config: testAccUserResourceConfig(email, "user", "Test", "DataSource"),
			},
			// Read the user using data source by ID
			{
				Config: testAccUserDataSourceConfig_byId(email, "user", "Test", "DataSource"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.n8ncloud_user.test",
						tfjsonpath.New("email"),
						knownvalue.StringExact(email),
					),
					statecheck.ExpectKnownValue(
						"data.n8ncloud_user.test",
						tfjsonpath.New("role"),
						knownvalue.StringExact("user"),
					),
					statecheck.ExpectKnownValue(
						"data.n8ncloud_user.test",
						tfjsonpath.New("first_name"),
						knownvalue.StringExact("Test"),
					),
					statecheck.ExpectKnownValue(
						"data.n8ncloud_user.test",
						tfjsonpath.New("last_name"),
						knownvalue.StringExact("DataSource"),
					),
				},
			},
			// Read the user using data source by email
			{
				Config: testAccUserDataSourceConfig_byEmail(email, "user", "Test", "DataSource"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.n8ncloud_user.test_by_email",
						tfjsonpath.New("email"),
						knownvalue.StringExact(email),
					),
					statecheck.ExpectKnownValue(
						"data.n8ncloud_user.test_by_email",
						tfjsonpath.New("role"),
						knownvalue.StringExact("user"),
					),
					statecheck.ExpectKnownValue(
						"data.n8ncloud_user.test_by_email",
						tfjsonpath.New("first_name"),
						knownvalue.StringExact("Test"),
					),
					statecheck.ExpectKnownValue(
						"data.n8ncloud_user.test_by_email",
						tfjsonpath.New("last_name"),
						knownvalue.StringExact("DataSource"),
					),
				},
			},
		},
	})
}

func TestAccUserDataSource_notFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Try to read non-existent user by ID
			{
				Config:      testAccUserDataSourceConfig_notFoundById(),
				ExpectError: regexp.MustCompile(`User with ID "non-existent-id" not found`),
			},
			// Try to read non-existent user by email
			{
				Config:      testAccUserDataSourceConfig_notFoundByEmail(),
				ExpectError: regexp.MustCompile(`User with email "non-existent@example.com" not found`),
			},
		},
	})
}

func testAccUserDataSourceConfig_byId(email, role, firstName, lastName string) string {
	return fmt.Sprintf(`
resource "n8ncloud_user" "test" {
  email      = %[1]q
  role       = %[2]q
  first_name = %[3]q
  last_name  = %[4]q
}

data "n8ncloud_user" "test" {
  id = n8ncloud_user.test.id
}
`, email, role, firstName, lastName)
}

func testAccUserDataSourceConfig_byEmail(email, role, firstName, lastName string) string {
	return fmt.Sprintf(`
resource "n8ncloud_user" "test" {
  email      = %[1]q
  role       = %[2]q
  first_name = %[3]q
  last_name  = %[4]q
}

data "n8ncloud_user" "test_by_email" {
  email = n8ncloud_user.test.email
}
`, email, role, firstName, lastName)
}

func testAccUserDataSourceConfig_notFoundById() string {
	return `
data "n8ncloud_user" "test" {
  id = "non-existent-id"
}
`
}

func testAccUserDataSourceConfig_notFoundByEmail() string {
	return `
data "n8ncloud_user" "test" {
  email = "non-existent@example.com"
}
`
}

# Terraform Provider for n8n Cloud

This provider allows you to manage n8n cloud users using Terraform.

## Features

- **User Management**: Create, read, update, and delete n8n cloud users
- **Role Management**: Support for global:admin and global:member roles
- **Data Sources**: Query existing users by ID or email
- **Import Support**: Import existing users into Terraform state

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23 (for development)
- n8n Cloud instance with API access

## Using the Provider

### Install

This provider is available from the [Terraform Registry](https://registry.terraform.io/providers/ka2n/n8ncloud).

```hcl
terraform {
  required_providers {
    n8ncloud = {
      source = "ka2n/n8ncloud"
      version = "~> 0.1"
    }
  }
}
```

### Configuration

```hcl
provider "n8ncloud" {
  api_key      = var.n8n_api_key      # or set N8N_API_KEY environment variable
  instance_url = var.n8n_instance_url # or set N8N_INSTANCE_URL environment variable
  timeout      = 30                   # optional, defaults to 30 seconds
}
```

### Environment Variables

You can configure the provider using environment variables:

```bash
export N8N_API_KEY="your-api-key"
export N8N_INSTANCE_URL="https://yourinstance.app.n8n.cloud"
```

## Usage Examples

### Create a User

```hcl
resource "n8ncloud_user" "developer" {
  email = "developer@example.com"
  role  = "global:member"
}
```

### Create an Admin User

```hcl
resource "n8ncloud_user" "admin" {
  email = "admin@example.com"
  role  = "global:admin"
}
```

### Query an Existing User

```hcl
data "n8ncloud_user" "existing" {
  email = "existing@example.com"
}

output "user_info" {
  value = {
    id         = data.n8ncloud_user.existing.id
    role       = data.n8ncloud_user.existing.role
    is_pending = data.n8ncloud_user.existing.is_pending
  }
}
```

### Import an Existing User

```bash
terraform import n8ncloud_user.existing "user-uuid-here"
```

## Resource Reference

### `n8ncloud_user`

#### Schema

- `email` (String, Required) - The email address of the user. Changing this forces a new resource.
- `role` (String, Required) - The role of the user (`global:admin` or `global:member`).
- `id` (String, Read-only) - The unique identifier of the user.
- `first_name` (String, Read-only) - The first name of the user.
- `last_name` (String, Read-only) - The last name of the user.
- `is_pending` (Boolean, Read-only) - Whether the user has not yet set up their account.
- `created_at` (String, Read-only) - The timestamp when the user was created.
- `updated_at` (String, Read-only) - The timestamp when the user was last updated.
- `invite_accept_url` (String, Read-only) - The URL for the user to accept their invitation.

## Data Source Reference

### `n8ncloud_user`

#### Schema

- `id` (String, Optional) - The unique identifier of the user. Either `id` or `email` must be specified.
- `email` (String, Optional) - The email address of the user. Either `id` or `email` must be specified.
- `role` (String, Read-only) - The role of the user.
- `first_name` (String, Read-only) - The first name of the user.
- `last_name` (String, Read-only) - The last name of the user.
- `is_pending` (Boolean, Read-only) - Whether the user has not yet set up their account.
- `created_at` (String, Read-only) - The timestamp when the user was created.
- `updated_at` (String, Read-only) - The timestamp when the user was last updated.
- `invite_accept_url` (String, Read-only) - The URL for the user to accept their invitation.

## Development

### Prerequisites

- [Go](https://golang.org/doc/install) 1.23+
- [Terraform](https://developer.hashicorp.com/terraform/downloads) 1.0+

### Building the Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

### Testing

```bash
# Run unit tests
go test -v ./...

# Run acceptance tests (requires API key)
TF_ACC=1 go test -v ./internal/provider/
```

### Local Development

```bash
# Install the provider locally
go install

# Generate documentation
go generate ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for your changes
5. Ensure tests pass (`go test -v ./...`)
6. Commit your changes (`git commit -am 'Add some amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## License

This project is licensed under the Mozilla Public License v2.0 - see the [LICENSE](LICENSE) file for details.

## Support

- [GitHub Issues](https://github.com/ka2n/terraform-provider-n8ncloud/issues)

## Acknowledgments

- [HashiCorp Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework)
- [n8n Public API](https://docs.n8n.io/api/)

# Create a new n8n cloud user
resource "n8ncloud_user" "developer" {
  email = "developer@example.com"
  role  = "global:member"
}

# Create an admin user
resource "n8ncloud_user" "admin" {
  email = "admin@example.com"
  role  = "global:admin"
}

# Output the user information
output "developer_user" {
  value = {
    id         = n8ncloud_user.developer.id
    email      = n8ncloud_user.developer.email
    role       = n8ncloud_user.developer.role
    is_pending = n8ncloud_user.developer.is_pending
  }
}
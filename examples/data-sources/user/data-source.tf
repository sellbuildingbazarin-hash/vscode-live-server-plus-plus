# Query an existing user by email
data "n8ncloud_user" "existing_by_email" {
  email = "admin@example.com"
}

# Query an existing user by ID
data "n8ncloud_user" "existing_by_id" {
  id = "b282f6c9-7bd9-4276-be66-14431896d36d"
}

# Output the user information
output "existing_user_info" {
  value = {
    id         = data.n8ncloud_user.existing_by_email.id
    email      = data.n8ncloud_user.existing_by_email.email
    role       = data.n8ncloud_user.existing_by_email.role
    first_name = data.n8ncloud_user.existing_by_email.first_name
    last_name  = data.n8ncloud_user.existing_by_email.last_name
    is_pending = data.n8ncloud_user.existing_by_email.is_pending
    created_at = data.n8ncloud_user.existing_by_email.created_at
  }
}
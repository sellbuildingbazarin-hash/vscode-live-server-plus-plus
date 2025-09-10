terraform {
  required_providers {
    n8ncloud = {
      source = "ka2n/n8ncloud"
    }
  }
}

# Configure the n8n Cloud Provider
provider "n8ncloud" {
  api_key      = var.n8n_api_key      # or set N8N_API_KEY environment variable
  instance_url = var.n8n_instance_url # or set N8N_INSTANCE_URL environment variable
  timeout      = 30                   # optional, defaults to 30 seconds
}

# Local values for common tags
locals {
  common_tags = {
    Environment = "dev"
    Project     = "example"
    # ManagedBy = "Terraform" # Uncomment when ready
  }
}
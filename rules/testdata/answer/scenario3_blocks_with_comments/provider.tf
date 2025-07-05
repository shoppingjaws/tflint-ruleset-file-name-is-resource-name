# Provider configuration
provider "aws" {
  region = var.aws_region
  
  # Default tags applied to all resources
  default_tags {
    tags = {
      ManagedBy = "Terraform"
      # Environment = "dev" # Uncomment to add environment tag globally
    }
  }
}
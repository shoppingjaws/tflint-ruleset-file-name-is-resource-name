# Import existing resources
# Use this to import resources created outside of Terraform
import {
  to = aws_instance.legacy
  id = "i-1234567890abcdef0"
  # Run: terraform plan -generate-config-out=generated.tf
}
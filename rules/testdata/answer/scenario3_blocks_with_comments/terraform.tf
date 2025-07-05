# Terraform configuration block
terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  
  # Backend configuration can be added here
  # backend "s3" {
  #   bucket = "my-terraform-state"
  #   key    = "example/terraform.tfstate"
  #   region = "us-west-2"
  # }
}
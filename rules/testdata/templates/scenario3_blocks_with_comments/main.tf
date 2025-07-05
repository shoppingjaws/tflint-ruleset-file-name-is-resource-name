# This file contains various blocks with comments that should be preserved when moved

# Variables section
# These variables are used throughout the configuration
variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t2.micro"
  # TODO: Consider using t3.micro for better performance
}

# AWS region configuration
variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
}

# Output definitions
# This output exports the instance ID for use in other modules
output "instance_id" {
  description = "The ID of the EC2 instance"
  value       = aws_instance.example.id
  # Note: This ID can be used for monitoring setup
}

# Local values for common tags
locals {
  common_tags = {
    Environment = "dev"
    Project     = "example"
    # ManagedBy = "Terraform" # Uncomment when ready
  }
}

# Main EC2 instance resource
# This creates the primary application server
resource "aws_instance" "example" {
  ami           = "ami-12345678" # Amazon Linux 2
  instance_type = var.instance_type
  tags          = local.common_tags
  
  # TODO: Add user_data script for initialization
  # user_data = file("init.sh")
}

# Security group for the EC2 instance
resource "aws_security_group" "example" {
  name_prefix = "example-"
  
  # Allow SSH access
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # WARNING: Open to the world, restrict in production
  }
  
  # Allow all outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Data source for latest Ubuntu AMI
# This ensures we always use the latest available AMI
data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"] # Canonical
  
  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }
}

# VPC data source
data "aws_vpc" "default" {
  default = true
  # Using default VPC for simplicity
}

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

# VPC module for network setup
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  
  name = "my-vpc"
  cidr = "10.0.0.0/16"
  
  # Add more configuration as needed
  # azs             = ["us-west-2a", "us-west-2b"]
  # private_subnets = ["10.0.1.0/24", "10.0.2.0/24"]
  # public_subnets  = ["10.0.101.0/24", "10.0.102.0/24"]
}

# Import existing resources
# Use this to import resources created outside of Terraform
import {
  to = aws_instance.legacy
  id = "i-1234567890abcdef0"
  # Run: terraform plan -generate-config-out=generated.tf
}

# Track resource movements
moved {
  from = aws_instance.old_name
  to   = aws_instance.example
  # This helps preserve state during refactoring
}

# Mark resources for removal
removed {
  from = aws_security_group.deprecated
  
  lifecycle {
    destroy = false # Keep the resource in AWS even after removal from state
  }
}

# Health checks
check "instance_health" {
  # Verify instance is properly configured
  assert {
    condition     = aws_instance.example.instance_type != ""
    error_message = "Instance type must be specified"
  }
  
  # Additional checks can be added here
  # assert {
  #   condition     = contains(["t2.micro", "t3.micro"], aws_instance.example.instance_type)
  #   error_message = "Instance type must be t2.micro or t3.micro"
  # }
}
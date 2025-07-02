terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t2.micro"
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
}

output "instance_id" {
  description = "The ID of the EC2 instance"
  value       = aws_instance.example.id
}

output "vpc_id" {
  description = "The VPC ID"
  value       = module.vpc.vpc_id
}

locals {
  common_tags = {
    Environment = "dev"
    Project     = "example"
  }
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  version = "~> 5.0"
  
  cidr = "10.0.0.0/16"
  azs  = ["us-west-2a", "us-west-2b"]
  
  tags = local.common_tags
}

resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = var.instance_type
  subnet_id     = module.vpc.public_subnets[0]
  tags          = local.common_tags
}

resource "aws_security_group" "web" {
  name_prefix = "web-"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.common_tags
}

# Data sources
data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"] # Canonical

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

data "aws_availability_zones" "available" {
  state = "available"
}

data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

# Ephemeral resources (example - might not be real AWS resources yet)
ephemeral "aws_ephemeral_instance" "temp_compute" {
  instance_type = "t2.micro"
  duration      = "1h"
  
  lifecycle {
    cleanup_on_destroy = true
  }
}

ephemeral "aws_ephemeral_storage" "temp_storage" {
  size     = "10GB"
  type     = "gp3"
  duration = "2h"
}

ephemeral "aws_ephemeral_network" "temp_network" {
  cidr_block = "10.0.0.0/24"
  duration   = "30m"
}

# Import blocks (Terraform 1.5+)
import {
  to = aws_instance.existing
  id = "i-1234567890abcdef0"
}

import {
  to = aws_security_group.imported
  id = "sg-0123456789abcdef0"
}

# Moved blocks (Terraform 1.1+)
moved {
  from = aws_instance.old_instance
  to   = aws_instance.example
}

moved {
  from = module.old_vpc
  to   = module.vpc
}

# Removed blocks (Terraform 1.7+)
removed {
  from = aws_instance.deprecated
  lifecycle {
    destroy = false
  }
}

removed {
  from = aws_security_group.legacy
  lifecycle {
    destroy = true
  }
}

# Check blocks (Terraform 1.5+)
check "instance_health" {
  assert {
    condition     = aws_instance.example.instance_state == "running"
    error_message = "EC2 instance must be in running state"
  }
}

check "vpc_configuration" {
  assert {
    condition     = length(module.vpc.public_subnets) >= 2
    error_message = "VPC must have at least 2 public subnets"
  }
  
  assert {
    condition     = module.vpc.enable_nat_gateway == true
    error_message = "NAT gateway must be enabled for the VPC"
  }
}
terraform {
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

locals {
  common_tags = {
    Environment = "dev"
    Project     = "example"
  }
}

output "instance_id" {
  description = "The ID of the EC2 instance"
  value       = aws_instance.example.id
}

output "vpc_id" {
  description = "The VPC ID"
  value       = module.vpc.vpc_id
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "my-vpc"
  cidr = "10.0.0.0/16"

  azs             = ["us-west-2a", "us-west-2b", "us-west-2c"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]

  enable_nat_gateway = true
  enable_vpn_gateway = true

  tags = local.common_tags
}

resource "aws_instance" "example" {
  ami           = "ami-0c02fb55956c7d316"
  instance_type = var.instance_type

  vpc_security_group_ids = [aws_security_group.example.id]
  subnet_id              = module.vpc.public_subnets[0]

  tags = local.common_tags
}

resource "aws_security_group" "example" {
  name_prefix = "example-"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.common_tags
}

# Data sources with custom prefix (datasource_)
data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"] # Canonical

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }
}

data "aws_vpc" "default" {
  default = true
}

# Ephemeral resources with custom prefix (temp_)
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

# Import blocks (Terraform 1.5+)
import {
  to = aws_instance.imported_ec2
  id = "i-0987654321fedcba0"
}

# Moved blocks (Terraform 1.1+)
moved {
  from = aws_instance.legacy_instance
  to   = aws_instance.example
}

# Removed blocks (Terraform 1.7+)
removed {
  from = aws_security_group.deprecated_sg
  lifecycle {
    destroy = false
  }
}

# Check blocks (Terraform 1.5+)
check "security_validation" {
  assert {
    condition     = length(aws_security_group.example.ingress) > 0
    error_message = "Security group must have at least one ingress rule"
  }
}

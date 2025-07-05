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
  name   = "my-vpc"
  cidr   = "10.0.0.0/16"
}

resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = var.instance_type
  subnet_id     = module.vpc.public_subnets[0]
  tags          = local.common_tags
}

resource "aws_security_group" "example" {
  name_prefix = "example-"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"]
}

data "aws_vpc" "default" {
  default = true
}

data "aws_availability_zones" "available" {
  state = "available"
}

data "aws_subnets" "private" {
  filter {
    name   = "vpc-id"
    values = [module.vpc.vpc_id]
  }
}

ephemeral "aws_ephemeral_instance" "temp_compute" {
  instance_type = "t2.micro"
  duration      = "1h"
}

ephemeral "aws_ephemeral_storage" "temp_storage" {
  size = "10GB"
  type = "gp3"
}

ephemeral "aws_ephemeral_network" "temp_network" {
  type = "isolated"
}

import {
  to = aws_instance.imported
  id = "i-1234567890abcdef0"
}

import {
  to = aws_security_group.imported
  id = "sg-0987654321fedcba0"
}

moved {
  from = aws_instance.old
  to   = aws_instance.example
}

removed {
  from = aws_instance.deprecated

  lifecycle {
    destroy = false
  }
}

removed {
  from = aws_security_group.old_sg
}

check "instance_validation" {
  assert {
    condition     = aws_instance.example.instance_type == var.instance_type
    error_message = "Instance type must match variable"
  }
}

check "security_validation" {
  assert {
    condition     = length(aws_security_group.example.ingress) > 0
    error_message = "Security group must have at least one ingress rule"  
  }
}
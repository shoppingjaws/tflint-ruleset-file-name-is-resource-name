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

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
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

# Data sources that might use ephemeral resources
data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"] # Canonical

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }
}

# Regular resources that might reference ephemeral resources
resource "aws_instance" "example" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = ephemeral.aws_ephemeral_instance.temp_compute.instance_type

  tags = {
    Name = "Example Instance"
    Mode = "Using Ephemeral"
  }
}
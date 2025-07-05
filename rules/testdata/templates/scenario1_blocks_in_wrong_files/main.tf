variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t2.micro"
}

output "instance_id" {
  description = "The ID of the EC2 instance"
  value       = aws_instance.example.id
}

locals {
  common_tags = {
    Environment = "dev"
    Project     = "example"
  }
}

resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = var.instance_type
  tags          = local.common_tags
}
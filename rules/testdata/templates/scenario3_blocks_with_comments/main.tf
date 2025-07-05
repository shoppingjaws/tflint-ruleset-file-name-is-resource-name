# This is the main AWS provider configuration
provider "aws" {
  region = "us-west-2"
}

// This variable defines the instance type
variable "instance_type" {
  type    = string
  default = "t2.micro"
}

# Main EC2 instance for the application
# This runs our web server
resource "aws_instance" "web" {
  ami           = "ami-12345678"
  instance_type = var.instance_type
}

// Security group for web traffic
// Allows HTTP and HTTPS
resource "aws_security_group" "web" {
  name = "web-sg"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Data source to get availability zones
data "aws_availability_zones" "available" {
  state = "available"
}

// Fetch the latest Ubuntu AMI
data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"]
}

# Output the instance ID
output "instance_id" {
  value = aws_instance.web.id
}

// Output the security group ID
output "security_group_id" {
  value = aws_security_group.web.id
}
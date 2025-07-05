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
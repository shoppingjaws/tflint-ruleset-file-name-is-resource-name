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
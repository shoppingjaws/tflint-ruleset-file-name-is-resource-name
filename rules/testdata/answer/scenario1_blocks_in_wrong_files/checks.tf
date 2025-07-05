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
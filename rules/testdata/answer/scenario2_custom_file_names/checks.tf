# Check blocks (Terraform 1.5+)
check "security_validation" {
  assert {
    condition     = length(aws_security_group.example.ingress) > 0
    error_message = "Security group must have at least one ingress rule"
  }
}
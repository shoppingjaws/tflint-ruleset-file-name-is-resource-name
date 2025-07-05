# Removed blocks (Terraform 1.7+)
removed {
  from = aws_security_group.deprecated_sg
  lifecycle {
    destroy = false
  }
}
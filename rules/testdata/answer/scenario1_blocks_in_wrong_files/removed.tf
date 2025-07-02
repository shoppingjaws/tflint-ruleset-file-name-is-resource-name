# Removed blocks (Terraform 1.7+)
removed {
  from = aws_instance.deprecated
  lifecycle {
    destroy = false
  }
}

removed {
  from = aws_security_group.legacy
  lifecycle {
    destroy = true
  }
}
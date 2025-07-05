# Mark resources for removal
removed {
  from = aws_security_group.deprecated
  
  lifecycle {
    destroy = false # Keep the resource in AWS even after removal from state
  }
}
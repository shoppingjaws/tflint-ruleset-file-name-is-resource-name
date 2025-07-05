removed {
  from = aws_instance.deprecated

  lifecycle {
    destroy = false
  }
}

removed {
  from = aws_security_group.old_sg
}
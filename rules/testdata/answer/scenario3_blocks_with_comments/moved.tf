# Track resource movements
moved {
  from = aws_instance.old_name
  to   = aws_instance.example
  # This helps preserve state during refactoring
}
# Import blocks (Terraform 1.5+)
import {
  to = aws_instance.existing
  id = "i-1234567890abcdef0"
}

import {
  to = aws_security_group.imported
  id = "sg-0123456789abcdef0"
}
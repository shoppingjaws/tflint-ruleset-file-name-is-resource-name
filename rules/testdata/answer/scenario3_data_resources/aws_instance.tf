resource "aws_instance" "example" {
  ami                    = data.aws_ami.ubuntu.id
  instance_type          = "t2.micro"
  availability_zone      = data.aws_availability_zones.available.names[0]
  subnet_id              = data.aws_subnets.default.ids[0]
  vpc_security_group_ids = [aws_security_group.example.id]

  tags = {
    Name = "Example Instance"
  }
}
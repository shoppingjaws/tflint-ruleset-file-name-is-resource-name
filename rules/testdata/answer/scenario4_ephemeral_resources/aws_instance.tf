resource "aws_instance" "example" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = ephemeral.aws_ephemeral_instance.temp_compute.instance_type

  tags = {
    Name = "Example Instance"
    Mode = "Using Ephemeral"
  }
}
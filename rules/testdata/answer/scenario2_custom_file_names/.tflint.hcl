plugin "file-name-is-resource-name" {
  enabled = true
}

rule "terraform_block_file_name" {
  enabled = true
  variable_file  = "vars.tf"
  output_file    = "out.tf"
  locals_file    = "local.tf"
  terraform_file = "tf.tf"
  provider_file  = "providers.tf"
  module_file    = "modules.tf"
}

rule "terraform_required_version" {
  enabled = false
}

rule "terraform_required_providers" {
  enabled = false
}
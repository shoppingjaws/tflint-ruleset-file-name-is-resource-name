plugin "file-name-is-resource-name" {
  enabled = true
}

rule "terraform_block_file_name" {
  enabled = true
  ephemeral_prefix = "temp_"
}

rule "terraform_required_version" {
  enabled = false
}

rule "terraform_required_providers" {
  enabled = false
}
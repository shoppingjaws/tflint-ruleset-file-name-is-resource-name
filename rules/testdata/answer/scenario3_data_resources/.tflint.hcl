plugin "file-name-is-resource-name" {
  enabled = true
}

rule "terraform_block_file_name" {
  enabled = true
  data_prefix = "datasource_"
}

rule "terraform_required_version" {
  enabled = false
}

rule "terraform_required_providers" {
  enabled = false
}
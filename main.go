package main

import (
	"github.com/shoppingjaws/tflint-ruleset-file-name-is-resource-name/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "file-name-is-resource-name",
			Version: "0.1.0",
			Rules: []tflint.Rule{
				rules.NewTerraformBlockFileNameRule(),
			},
		},
	})
}

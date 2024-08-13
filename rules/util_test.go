package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
)

func Test_GetBlocksFromBody(t *testing.T) {
	t.Run("", func(t *testing.T) {
		file, diags := hclsyntax.ParseConfig([]byte(`
			resource "resource_type" "resource_name" {}
			data "data_type" "data_name" {}
			variable "variable_name" {}
			locals { local_name = "value" }
			provider "provider_name" {}
			module "module_name" {}
			output "output_name" {}
			`), "test.hcl", hcl.Pos{Line: 1, Column: 1})
		if diags.HasErrors() {
			t.Fatal(diags)
		}
		result, _ := GetBlocksFromBody(file.Body)
		assert.Equal(t, 1, len(result.Filter(Resource)))
		assert.Equal(t, 1, len(result.Filter(Data)))
		assert.Equal(t, 1, len(result.Filter(Variable)))
		assert.Equal(t, 1, len(result.Filter(Locals)))
		assert.Equal(t, 1, len(result.Filter(Provider)))
		assert.Equal(t, 1, len(result.Filter(Module)))
		assert.Equal(t, 1, len(result.Filter(Output)))
	})

}

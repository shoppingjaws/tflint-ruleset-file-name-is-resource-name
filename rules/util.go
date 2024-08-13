package rules

import (
	hcl "github.com/hashicorp/hcl/v2"
)

type BlockKindEnum int

const (
	Resource BlockKindEnum = iota
	Data
	Variable
	Locals
	Provider
	Module
	Output
)

var blockKindMap = map[BlockKindEnum]hcl.BlockHeaderSchema{
	Resource: {Type: "resource", LabelNames: []string{"name", "type"}},
	Data:     {Type: "data", LabelNames: []string{"name", "type"}},
	Variable: {Type: "variable", LabelNames: []string{"name"}},
	Locals:   {Type: "locals"},
	Provider: {Type: "provider", LabelNames: []string{"name"}},
	Module:   {Type: "module", LabelNames: []string{"name"}},
	Output:   {Type: "output", LabelNames: []string{"name"}},
}

type Block struct {
	Kind  BlockKindEnum
	Name  *string
	Type  *string
	Range hcl.Range
}
type BlockList []Block

func toBlockKindEnum(s string) BlockKindEnum {
	switch s {
	case "resource":
		return Resource
	case "data":
		return Data
	case "variable":
		return Variable
	case "locals":
		return Locals
	case "provider":
		return Provider
	case "module":
		return Module
	case "output":
		return Output
	}
	panic("invalid block kind")
}

func (b BlockList) Filter(kind BlockKindEnum) []Block {
	var blocks []Block
	for _, block := range b {
		if block.Kind == kind {
			blocks = append(blocks, block)
		}
	}
	return blocks
}

func (b BlockList) Exclude(kind BlockKindEnum) []Block {
	var blocks []Block
	for _, block := range b {
		if block.Kind != kind {
			blocks = append(blocks, block)
		}
	}
	return blocks
}

func GetBlocksFromBody(body hcl.Body) (*BlockList, error) {
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			blockKindMap[Resource],
			blockKindMap[Data],
			blockKindMap[Variable],
			blockKindMap[Locals],
			blockKindMap[Provider],
			blockKindMap[Module],
			blockKindMap[Output],
		},
	}
	content, _, diags := body.PartialContent(schema)
	if diags.HasErrors() {
		return nil, diags
	}
	var blockList BlockList
	for _, block := range content.Blocks {
		blockList = append(blockList, Block{
			Kind: BlockKindEnum(toBlockKindEnum(block.Type)),
			/*
			 * HCLのパース結果とLabelNamesは逆順になっているため、後ろからアクセスする
			 * resource "resource_type" "resource_name" {}
			 *           [0]            [1]
			 * output   "output_name" {}
			 *           [0]
			 */
			Name:  safeAccess(block.Labels, -1),
			Type:  safeAccess(block.Labels, -2),
			Range: block.DefRange,
		})
	}
	return &blockList, nil
}

func safeAccess[T any](arr []T, index int) *T {
	if index >= 0 && index < len(arr) {
		return &arr[index]
	} else if index < 0 && -index <= len(arr) {
		return &arr[len(arr)+index]
	}
	return nil
}

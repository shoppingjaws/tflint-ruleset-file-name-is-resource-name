package rules

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestBlockManager_ExtractBlockText_PreservesComments(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		blockType      string
		blockLabel     string
		expectedOutput string
	}{
		{
			name: "resource block with single line comment using #",
			content: `# This is an important EC2 instance
resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}`,
			blockType:  "resource",
			blockLabel: "aws_instance",
			expectedOutput: `# This is an important EC2 instance
resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}`,
		},
		{
			name: "resource block with single line comment using //",
			content: `// This is an important security group
resource "aws_security_group" "web" {
  name_prefix = "web-"
  vpc_id      = "vpc-12345"
}`,
			blockType:  "resource",
			blockLabel: "aws_security_group",
			expectedOutput: `// This is an important security group
resource "aws_security_group" "web" {
  name_prefix = "web-"
  vpc_id      = "vpc-12345"
}`,
		},
		{
			name: "data block with comment",
			content: `# Fetch the latest Ubuntu AMI
data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"]
}`,
			blockType:  "data",
			blockLabel: "aws_ami",
			expectedOutput: `# Fetch the latest Ubuntu AMI
data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"]
}`,
		},
		{
			name: "variable block with comment",
			content: `// The AWS region to deploy resources
variable "aws_region" {
  type    = string
  default = "us-west-2"
}`,
			blockType:  "variable",
			blockLabel: "aws_region",
			expectedOutput: `// The AWS region to deploy resources
variable "aws_region" {
  type    = string
  default = "us-west-2"
}`,
		},
		{
			name: "output block with comment",
			content: `# The ID of the created instance
output "instance_id" {
  value = aws_instance.example.id
}`,
			blockType:  "output",
			blockLabel: "instance_id",
			expectedOutput: `# The ID of the created instance
output "instance_id" {
  value = aws_instance.example.id
}`,
		},
		{
			name: "resource block with multiple comment lines",
			content: `# This resource creates a VPC
# It includes both public and private subnets
# Use this for production deployments
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}`,
			blockType:  "resource",
			blockLabel: "aws_vpc",
			expectedOutput: `# This resource creates a VPC
# It includes both public and private subnets
# Use this for production deployments
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}`,
		},
		{
			name: "resource block with comment after empty line should not include comment",
			content: `# This comment is separated by empty line

resource "aws_instance" "example" {
  ami = "ami-12345678"
}`,
			blockType:  "resource",
			blockLabel: "aws_instance",
			expectedOutput: `resource "aws_instance" "example" {
  ami = "ami-12345678"
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{
				"test.tf": tt.content,
			})

			bm := NewBlockManager(runner)

			// Find the block
			var targetBlock *hclext.Block
			body, _ := runner.GetModuleContent(&hclext.BodySchema{
				Blocks: []hclext.BlockSchema{
					{Type: tt.blockType, LabelNames: getLabelNamesForBlockType(tt.blockType)},
				},
			}, nil)

			for _, block := range body.Blocks {
				if block.Type == tt.blockType {
					targetBlock = block
					break
				}
			}

			assert.NotNil(t, targetBlock, "Block not found")

			// Extract block text
			extractedText, err := bm.ExtractBlockText(targetBlock)
			assert.NoError(t, err)

			// Normalize whitespace for comparison
			expectedNormalized := strings.TrimSpace(tt.expectedOutput)
			actualNormalized := strings.TrimSpace(extractedText)

			assert.Equal(t, expectedNormalized, actualNormalized, "Extracted text doesn't match expected output")
		})
	}
}

func getLabelNamesForBlockType(blockType string) []string {
	switch blockType {
	case "resource", "data":
		return []string{"type", "name"}
	case "variable", "output", "module":
		return []string{"name"}
	case "provider":
		return []string{"name"}
	case "locals", "terraform":
		return []string{}
	default:
		return []string{}
	}
}
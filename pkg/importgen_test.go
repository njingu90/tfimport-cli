package pkg

import (
	"strings"
	"testing"
)

func TestGenerateImportBlock(t *testing.T) {
	resource := Resource{
		Type:    "aws_vpc",
		Name:    "main",
		Module:  "",
		Address: "aws_vpc.main",
	}

	block := GenerateImportBlock(resource, "vpc-123456")

	if block.To != "aws_vpc.main" {
		t.Errorf("Expected 'aws_vpc.main', got '%s'", block.To)
	}

	if block.ID != "vpc-123456" {
		t.Errorf("Expected 'vpc-123456', got '%s'", block.ID)
	}
}

func TestFormatImportBlock(t *testing.T) {
	block := &ImportBlock{
		To: "aws_vpc.main",
		ID: "vpc-123456",
	}

	formatted := FormatImportBlock(block)

	if !strings.Contains(formatted, "import {") {
		t.Error("Missing 'import {' in formatted output")
	}

	if !strings.Contains(formatted, "to = aws_vpc.main") {
		t.Error("Missing 'to = aws_vpc.main' in formatted output")
	}

	if !strings.Contains(formatted, "id = \"vpc-123456\"") {
		t.Error("Missing 'id = \"vpc-123456\"' in formatted output")
	}

	if !strings.Contains(formatted, "}") {
		t.Error("Missing closing brace in formatted output")
	}
}

func TestSortImportBlocks(t *testing.T) {
	blocks := []ImportBlock{
		{To: "module.network.aws_vpc.main", ID: "vpc-3"},
		{To: "aws_vpc.secondary", ID: "vpc-1"},
		{To: "aws_vpc.primary", ID: "vpc-2"},
	}

	sorted := SortImportBlocks(blocks)

	expected := []string{
		"aws_vpc.primary",
		"aws_vpc.secondary",
		"module.network.aws_vpc.main",
	}

	for i, block := range sorted {
		if block.To != expected[i] {
			t.Errorf("Position %d: expected %s, got %s", i, expected[i], block.To)
		}
	}
}

func TestFormatImportBlocksAsFile(t *testing.T) {
	blocks := []ImportBlock{
		{To: "aws_vpc.main", ID: "vpc-123456"},
		{To: "aws_subnet.main", ID: "subnet-789012"},
	}

	formatted := FormatImportBlocksAsFile(blocks)

	if !strings.Contains(formatted, "# This file contains Terraform import blocks") {
		t.Error("Missing header comment")
	}

	if !strings.Contains(formatted, "import {") {
		t.Error("Missing import block")
	}

	if !strings.Contains(formatted, "terraform apply") {
		t.Error("Missing terraform apply instruction")
	}
}

func TestFilterImportBlocksByModule(t *testing.T) {
	blocks := []ImportBlock{
		{To: "module.network.aws_vpc.main", ID: "vpc-1"},
		{To: "module.network.aws_subnet.main", ID: "subnet-1"},
		{To: "module.compute.aws_instance.app", ID: "i-1"},
	}

	filtered := FilterImportBlocksByModule(blocks, "module.network")

	if len(filtered) != 2 {
		t.Errorf("Expected 2 blocks, got %d", len(filtered))
	}

	for _, block := range filtered {
		if !strings.HasPrefix(block.To, "module.network.") {
			t.Errorf("Block %s not in module.network", block.To)
		}
	}
}

func TestGroupImportBlocksByModule(t *testing.T) {
	blocks := []ImportBlock{
		{To: "module.network.aws_vpc.main", ID: "vpc-1"},
		{To: "module.network.aws_subnet.main", ID: "subnet-1"},
		{To: "module.compute.aws_instance.app", ID: "i-1"},
		{To: "aws_security_group.root", ID: "sg-1"},
	}

	grouped := GroupImportBlocksByModule(blocks)

	if len(grouped) != 3 {
		t.Errorf("Expected 3 groups, got %d", len(grouped))
	}

	if len(grouped["module.network"]) != 2 {
		t.Errorf("Expected 2 blocks in module.network, got %d", len(grouped["module.network"]))
	}

	if len(grouped["module.compute"]) != 1 {
		t.Errorf("Expected 1 block in module.compute, got %d", len(grouped["module.compute"]))
	}

	if len(grouped["root"]) != 1 {
		t.Errorf("Expected 1 block in root, got %d", len(grouped["root"]))
	}
}

func TestGroupImportBlocksByResourceType(t *testing.T) {
	blocks := []ImportBlock{
		{To: "aws_vpc.main", ID: "vpc-1"},
		{To: "aws_vpc.secondary", ID: "vpc-2"},
		{To: "aws_subnet.main", ID: "subnet-1"},
		{To: "module.network.aws_security_group.main", ID: "sg-1"},
	}

	grouped := GroupImportBlocksByResourceType(blocks)

	if len(grouped) != 3 {
		t.Errorf("Expected 3 groups, got %d", len(grouped))
	}

	if len(grouped["aws_vpc"]) != 2 {
		t.Errorf("Expected 2 aws_vpc blocks, got %d", len(grouped["aws_vpc"]))
	}

	if len(grouped["aws_subnet"]) != 1 {
		t.Errorf("Expected 1 aws_subnet block, got %d", len(grouped["aws_subnet"]))
	}

	if len(grouped["aws_security_group"]) != 1 {
		t.Errorf("Expected 1 aws_security_group block, got %d", len(grouped["aws_security_group"]))
	}
}

func TestCountImportBlocksByResourceType(t *testing.T) {
	blocks := []ImportBlock{
		{To: "aws_vpc.main", ID: "vpc-1"},
		{To: "aws_vpc.secondary", ID: "vpc-2"},
		{To: "aws_subnet.main", ID: "subnet-1"},
	}

	counts := CountImportBlocksByResourceType(blocks)

	if counts["aws_vpc"] != 2 {
		t.Errorf("Expected 2 aws_vpc, got %d", counts["aws_vpc"])
	}

	if counts["aws_subnet"] != 1 {
		t.Errorf("Expected 1 aws_subnet, got %d", counts["aws_subnet"])
	}
}

func TestCountImportBlocksByModule(t *testing.T) {
	blocks := []ImportBlock{
		{To: "module.network.aws_vpc.main", ID: "vpc-1"},
		{To: "module.network.aws_subnet.main", ID: "subnet-1"},
		{To: "aws_security_group.root", ID: "sg-1"},
	}

	counts := CountImportBlocksByModule(blocks)

	if counts["module.network"] != 2 {
		t.Errorf("Expected 2 in module.network, got %d", counts["module.network"])
	}

	if counts["root"] != 1 {
		t.Errorf("Expected 1 in root, got %d", counts["root"])
	}
}

func TestGenerateImportBlocksFromResources(t *testing.T) {
	resources := []Resource{
		{
			Type:    "aws_vpc",
			Address: "aws_vpc.main",
			Attributes: map[string]interface{}{
				"id": "vpc-123456",
			},
		},
		{
			Type:    "aws_subnet",
			Address: "aws_subnet.main",
			Attributes: map[string]interface{}{
				"id": "subnet-789012",
			},
		},
	}

	blocks, skipped := GenerateImportBlocksFromResources(resources)

	if len(blocks) != 2 {
		t.Errorf("Expected 2 blocks, got %d", len(blocks))
	}

	if len(skipped) != 0 {
		t.Errorf("Expected 0 skipped, got %d", len(skipped))
	}
}

func TestGenerateImportBlocksFromResources_Unsupported(t *testing.T) {
	resources := []Resource{
		{
			Type:    "google_compute_instance",
			Address: "google_compute_instance.main",
			Attributes: map[string]interface{}{
				"id": "instance-123",
			},
		},
	}

	blocks, skipped := GenerateImportBlocksFromResources(resources)

	if len(blocks) != 0 {
		t.Errorf("Expected 0 blocks, got %d", len(blocks))
	}

	if len(skipped) != 1 {
		t.Errorf("Expected 1 skipped, got %d", len(skipped))
	}

	if len(skipped) > 0 && skipped[0].Address != "google_compute_instance.main" {
		t.Errorf("Expected 'google_compute_instance.main', got '%s'", skipped[0].Address)
	}
}

package pkg

import (
	"testing"
)

func TestFilterByModule(t *testing.T) {
	resources := []Resource{
		{Type: "aws_vpc", Name: "main", Module: "", Address: "aws_vpc.main"},
		{Type: "aws_subnet", Name: "private", Module: "module.network", Address: "module.network.aws_subnet.private"},
		{Type: "aws_instance", Name: "app", Module: "module.compute", Address: "module.compute.aws_instance.app"},
	}

	filtered := FilterByModule(resources, "module.network")
	if len(filtered) != 1 {
		t.Errorf("Expected 1 resource, got %d", len(filtered))
	}

	if filtered[0].Module != "module.network" {
		t.Errorf("Expected module.network, got %s", filtered[0].Module)
	}
}

func TestFilterByResourceType(t *testing.T) {
	resources := []Resource{
		{Type: "aws_vpc", Name: "main", Module: "", Address: "aws_vpc.main"},
		{Type: "aws_subnet", Name: "private", Module: "module.network", Address: "module.network.aws_subnet.private"},
		{Type: "aws_vpc", Name: "secondary", Module: "module.network", Address: "module.network.aws_vpc.secondary"},
	}

	filtered := FilterByResourceType(resources, "aws_vpc")
	if len(filtered) != 2 {
		t.Errorf("Expected 2 resources, got %d", len(filtered))
	}

	for _, res := range filtered {
		if res.Type != "aws_vpc" {
			t.Errorf("Expected aws_vpc, got %s", res.Type)
		}
	}
}

func TestFilterByResourceAddress(t *testing.T) {
	resources := []Resource{
		{Type: "aws_vpc", Name: "main", Module: "", Address: "aws_vpc.main"},
		{Type: "aws_subnet", Name: "private", Module: "module.network", Address: "module.network.aws_subnet.private"},
	}

	res := FilterByResourceAddress(resources, "aws_vpc.main")
	if res == nil {
		t.Fatal("Expected resource, got nil")
	}

	if res.Address != "aws_vpc.main" {
		t.Errorf("Expected aws_vpc.main, got %s", res.Address)
	}

	// Test non-existent address
	res = FilterByResourceAddress(resources, "nonexistent")
	if res != nil {
		t.Fatal("Expected nil for non-existent resource")
	}
}

func TestFilterByPrefix(t *testing.T) {
	resources := []Resource{
		{Type: "aws_vpc", Name: "main", Module: "", Address: "aws_vpc.main"},
		{Type: "aws_subnet", Name: "private", Module: "module.network", Address: "module.network.aws_subnet.private"},
		{Type: "aws_instance", Name: "app", Module: "module.network.nested", Address: "module.network.nested.aws_instance.app"},
	}

	filtered := FilterByPrefix(resources, "module.network")
	if len(filtered) != 2 {
		t.Errorf("Expected 2 resources, got %d", len(filtered))
	}
}

func TestGetTopResourceTypes(t *testing.T) {
	resources := []Resource{
		{Type: "aws_vpc"},
		{Type: "aws_vpc"},
		{Type: "aws_subnet"},
		{Type: "aws_subnet"},
		{Type: "aws_subnet"},
		{Type: "aws_instance"},
	}

	top := GetTopResourceTypes(resources, 2)
	if len(top) != 2 {
		t.Errorf("Expected 2 types, got %d", len(top))
	}

	// First should be aws_subnet (count: 3)
	if top[0].Type != "aws_subnet" || top[0].Count != 3 {
		t.Errorf("Expected aws_subnet with count 3, got %s with count %d", top[0].Type, top[0].Count)
	}

	// Second should be aws_vpc (count: 2)
	if top[1].Type != "aws_vpc" || top[1].Count != 2 {
		t.Errorf("Expected aws_vpc with count 2, got %s with count %d", top[1].Type, top[1].Count)
	}
}

func TestFilterByResourceTypesSupported(t *testing.T) {
	resources := []Resource{
		{Type: "aws_vpc"},
		{Type: "aws_subnet"},
		{Type: "aws_custom_resource"},
		{Type: "aws_instance"},
	}

	supportedTypes := map[string]bool{
		"aws_vpc":      true,
		"aws_subnet":   true,
		"aws_instance": true,
	}

	supported, unsupported := FilterByResourceTypesSupported(resources, supportedTypes)

	if len(supported) != 3 {
		t.Errorf("Expected 3 supported, got %d", len(supported))
	}

	if len(unsupported) != 1 {
		t.Errorf("Expected 1 unsupported, got %d", len(unsupported))
	}
}

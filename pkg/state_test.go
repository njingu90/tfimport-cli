package pkg

import (
	"path/filepath"
	"runtime"
	"testing"
)

// testdataPath returns the absolute path to a testdata file
func testdataPath(filename string) string {
	// Get the directory of this test file
	_, currentFile, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(currentFile)
	rootDir := filepath.Dir(testDir)
	return filepath.Join(rootDir, "testdata", filename)
}

func TestLoadLocalState_Success(t *testing.T) {
	state, err := LoadLocalState(testdataPath("sample.tfstate"))
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}

	if state == nil {
		t.Fatal("State is nil")
	}

	if state.Version != 4 {
		t.Errorf("Expected version 4, got %d", state.Version)
	}

	if len(state.Resources) == 0 {
		t.Fatal("Expected resources, got none")
	}
}

func TestLoadLocalState_FileNotFound(t *testing.T) {
	_, err := LoadLocalState("nonexistent.tfstate")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestValidateState_ValidState(t *testing.T) {
	state := &TerraformState{
		Version:   4,
		Resources: []StateResource{},
	}

	err := ValidateState(state)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestValidateState_InvalidVersion(t *testing.T) {
	state := &TerraformState{
		Version:   3,
		Resources: []StateResource{},
	}

	err := ValidateState(state)
	if err == nil {
		t.Fatal("Expected error for invalid version")
	}
}

func TestGetAllResources(t *testing.T) {
	state, err := LoadLocalState(testdataPath("sample.tfstate"))
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}

	resources := GetAllResources(state)
	if len(resources) == 0 {
		t.Fatal("Expected resources, got none")
	}

	// Check that root resources have correct format (using a managed resource, not data source)
	found := false
	for _, res := range resources {
		if res.Type == "aws_cloudwatch_log_group" && res.Name == "flow_log" && res.Module == "" {
			found = true
			if res.Address != "aws_cloudwatch_log_group.flow_log" {
				t.Errorf("Expected address 'aws_cloudwatch_log_group.flow_log', got '%s'", res.Address)
			}
		}
	}

	if !found {
		t.Fatal("Did not find root aws_cloudwatch_log_group resource")
	}
}

func TestGetResourcesByModule(t *testing.T) {
	state, err := LoadLocalState(testdataPath("sample_modules.tfstate"))
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}

	resources := GetResourcesByModule(state, "module.vpc")
	if len(resources) == 0 {
		t.Fatal("Expected resources in module.vpc")
	}

	// All resources should belong to module.vpc
	for _, res := range resources {
		if res.Module != "module.vpc" {
			t.Errorf("Expected module 'module.vpc', got '%s'", res.Module)
		}
	}
}

func TestGetAllModules(t *testing.T) {
	state, err := LoadLocalState(testdataPath("sample_modules.tfstate"))
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}

	modules := GetAllModules(state)
	if len(modules) == 0 {
		t.Fatal("Expected modules")
	}

	// Check for known modules from real-world state
	hasVPC := false
	hasEC2 := false
	for _, m := range modules {
		if m == "module.vpc" {
			hasVPC = true
		}
		if m == "module.ec2_instance[\"one\"]" {
			hasEC2 = true
		}
	}

	if !hasVPC || !hasEC2 {
		t.Errorf("Missing expected modules: vpc=%v, ec2_instance=%v", hasVPC, hasEC2)
	}
}

func TestResourceCountByType(t *testing.T) {
	state, err := LoadLocalState(testdataPath("sample.tfstate"))
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}

	counts := ResourceCountByType(state)
	if len(counts) == 0 {
		t.Fatal("Expected resource counts")
	}

	// Check for resource types in real-world state
	if counts["aws_subnet"] == 0 {
		t.Fatal("Expected aws_subnet resources")
	}
	if counts["aws_iam_role"] == 0 {
		t.Fatal("Expected aws_iam_role resources")
	}
}

func TestGetProvidersInUse(t *testing.T) {
	state, err := LoadLocalState(testdataPath("sample.tfstate"))
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}

	providers := GetProvidersInUse(state)
	if len(providers) == 0 {
		t.Fatal("Expected providers")
	}

	// Should have aws provider
	found := false
	for _, p := range providers {
		if p == "provider[\"registry.terraform.io/hashicorp/aws\"]" {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("AWS provider not found")
	}
}

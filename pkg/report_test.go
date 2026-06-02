package pkg

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGenerateReport(t *testing.T) {
	skipped := []SkippedResourceInfo{
		{Address: "aws_custom.main", Type: "aws_custom", Reason: "unsupported"},
	}

	report := GenerateReport(100, 99, 98, 1, skipped)

	if report.ScannedCount != 100 {
		t.Errorf("Expected ScannedCount 100, got %d", report.ScannedCount)
	}

	if report.MatchedCount != 99 {
		t.Errorf("Expected MatchedCount 99, got %d", report.MatchedCount)
	}

	if report.GeneratedCount != 98 {
		t.Errorf("Expected GeneratedCount 98, got %d", report.GeneratedCount)
	}

	if report.SkippedCount != 1 {
		t.Errorf("Expected SkippedCount 1, got %d", report.SkippedCount)
	}

	if len(report.SkippedDetails) != 1 {
		t.Errorf("Expected 1 skipped detail, got %d", len(report.SkippedDetails))
	}
}

func TestFormatReportAsJSON(t *testing.T) {
	report := &Report{
		ScannedCount:   100,
		MatchedCount:   99,
		GeneratedCount: 98,
		SkippedCount:   1,
	}

	jsonStr, err := FormatReportAsJSON(report)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify it's valid JSON
	var data map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		t.Errorf("Invalid JSON: %v", err)
	}

	if data["scanned"] != float64(100) {
		t.Errorf("Expected scanned=100, got %v", data["scanned"])
	}
}

func TestFormatReportAsConsole(t *testing.T) {
	report := &Report{
		ScannedCount:   100,
		GeneratedCount: 98,
		SkippedCount:   2,
		ResourceTypes: map[string]int{
			"aws_vpc":    10,
			"aws_subnet": 20,
		},
		Modules: map[string]int{
			"module.network": 15,
			"root":           5,
		},
	}

	output := FormatReportAsConsole(report)

	if !strings.Contains(output, "Scanned") {
		t.Error("Missing 'Scanned' in output")
	}

	if !strings.Contains(output, "100") {
		t.Error("Missing scanned count in output")
	}

	if !strings.Contains(output, "Generated") {
		t.Error("Missing 'Generated' in output")
	}

	if !strings.Contains(output, "98") {
		t.Error("Missing generated count in output")
	}
}

func TestGenerateAnalysisReport(t *testing.T) {
	state := &TerraformState{
		Version: 4,
		Resources: []StateResource{
			{
				Type: "aws_vpc",
				Instances: []ResourceInstance{
					{Attributes: map[string]interface{}{"id": "vpc-1"}},
				},
			},
			{
				Type: "aws_subnet",
				Instances: []ResourceInstance{
					{Attributes: map[string]interface{}{"id": "subnet-1"}},
				},
			},
		},
	}

	report := GenerateAnalysisReport(state)

	if report.ScannedCount != 2 {
		t.Errorf("Expected ScannedCount 2, got %d", report.ScannedCount)
	}

	if len(report.ResourceTypes) != 2 {
		t.Errorf("Expected 2 resource types, got %d", len(report.ResourceTypes))
	}
}

func TestFormatAnalysisReportAsConsole(t *testing.T) {
	report := &Report{
		ScannedCount: 50,
		ResourceTypes: map[string]int{
			"aws_vpc":      5,
			"aws_subnet":   10,
			"aws_instance": 20,
		},
		Modules: map[string]int{
			"module.network": 15,
			"root":           5,
		},
	}

	output := FormatAnalysisReportAsConsole(report)

	if !strings.Contains(output, "State Analysis") {
		t.Error("Missing 'State Analysis' header")
	}

	if !strings.Contains(output, "50") {
		t.Error("Missing total resource count")
	}

	if !strings.Contains(output, "Top Resource Types") {
		t.Error("Missing 'Top Resource Types' section")
	}

	if !strings.Contains(output, "Modules") {
		t.Error("Missing 'Modules' section")
	}
}

func TestSkippedResourceInfo(t *testing.T) {
	skipped := SkippedResourceInfo{
		Address: "aws_custom.main",
		Type:    "aws_custom",
		Reason:  "unsupported resource type",
	}

	if skipped.Address != "aws_custom.main" {
		t.Errorf("Expected address 'aws_custom.main', got '%s'", skipped.Address)
	}

	if skipped.Type != "aws_custom" {
		t.Errorf("Expected type 'aws_custom', got '%s'", skipped.Type)
	}

	if skipped.Reason != "unsupported resource type" {
		t.Errorf("Expected reason 'unsupported resource type', got '%s'", skipped.Reason)
	}
}

func TestReportJSON(t *testing.T) {
	report := &Report{
		ScannedCount:     100,
		MatchedCount:     99,
		GeneratedCount:   98,
		SkippedCount:     1,
		UnsupportedCount: 2,
		Providers:        []string{"aws"},
		SkippedDetails: []SkippedResourceInfo{
			{
				Address: "aws_custom.main",
				Type:    "aws_custom",
				Reason:  "unsupported",
			},
		},
	}

	data, _ := json.MarshalIndent(report, "", "  ")
	jsonStr := string(data)

	// Verify structure
	var decoded map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &decoded)

	if decoded["scanned"] != float64(100) {
		t.Error("Invalid scanned count in JSON")
	}

	if decoded["matched"] != float64(99) {
		t.Error("Invalid matched count in JSON")
	}

	if skippedArr, ok := decoded["skipped_details"].([]interface{}); ok {
		if len(skippedArr) != 1 {
			t.Error("Invalid skipped_details array in JSON")
		}
	} else {
		t.Error("skipped_details is not an array")
	}
}

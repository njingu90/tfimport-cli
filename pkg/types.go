package pkg

import (
	"encoding/json"
)

// SensitiveAttributes handles both array and map formats for sensitive_attributes
// Terraform v1.0-1.4: map[string]bool
// Terraform v1.5+: []string (array of sensitive attribute paths)
type SensitiveAttributes struct {
	data map[string]bool
}

// UnmarshalJSON implements json.Unmarshaler to handle both formats
func (sa *SensitiveAttributes) UnmarshalJSON(b []byte) error {
	sa.data = make(map[string]bool)

	// Try unmarshaling as array first (Terraform v1.5+)
	var arr []string
	if err := json.Unmarshal(b, &arr); err == nil {
		// Convert array to map (array items are sensitive attribute paths)
		for _, path := range arr {
			sa.data[path] = true
		}
		return nil
	}

	// Fall back to map format (Terraform v1.0-1.4)
	var m map[string]bool
	if err := json.Unmarshal(b, &m); err != nil {
		// If both fail, treat as empty (no sensitive attributes)
		return nil
	}

	sa.data = m
	return nil
}

// Map returns the underlying map representation
func (sa *SensitiveAttributes) Map() map[string]bool {
	if sa.data == nil {
		return make(map[string]bool)
	}
	return sa.data
}
type TerraformState struct {
	Version          int                    `json:"version"`
	TerraformVersion string                 `json:"terraform_version"`
	Serial           int64                  `json:"serial"`
	Lineage          string                 `json:"lineage"`
	Outputs          map[string]interface{} `json:"outputs,omitempty"`
	Resources        []StateResource        `json:"resources"`
}

// StateResource represents a resource block in the state
type StateResource struct {
	Mode      string             `json:"mode"`      // "managed" or "data"
	Type      string             `json:"type"`
	Name      string             `json:"name"`
	Provider  string             `json:"provider"`
	Instances []ResourceInstance `json:"instances"`
	Module    string             `json:"module,omitempty"`
}

// ResourceInstance represents an instance of a resource
type ResourceInstance struct {
	SchemaVersion       int                    `json:"schema_version"`
	Attributes          map[string]interface{} `json:"attributes"`
	SensitiveAttributes SensitiveAttributes    `json:"sensitive_attributes,omitempty"`
	Private             string                 `json:"private,omitempty"`
	DependsOn           []string               `json:"depends_on,omitempty"`
	IndexKey            interface{}            `json:"index_key,omitempty"`
}

// Resource represents a processed resource for import
type Resource struct {
	Type           string
	Name           string
	Module         string
	Address        string // Full address like "module.network.aws_vpc.main"
	ImportID       string
	InstanceIndex  interface{} // for count/for_each
	Attributes     map[string]interface{}
	SupportedByAWS bool
	SkipReason     string // Reason if not supported or cannot generate import ID
}

// ImportBlock represents a Terraform import block
type ImportBlock struct {
	To string // Terraform address like "module.network.aws_vpc.main"
	ID string // Import ID
}

// Report represents the summary of an analysis/generation run
type Report struct {
	ScannedCount     int                   `json:"scanned"`
	MatchedCount     int                   `json:"matched"`
	GeneratedCount   int                   `json:"generated"`
	SkippedCount     int                   `json:"skipped"`
	SkippedDetails   []SkippedResourceInfo `json:"skipped_details"`
	ResourceTypes    map[string]int        `json:"resource_types,omitempty"`
	Modules          map[string]int        `json:"modules,omitempty"`
	UnsupportedCount int                   `json:"unsupported,omitempty"`
	Providers        []string              `json:"providers,omitempty"`
}

// SkippedResourceInfo tracks why a resource was skipped
type SkippedResourceInfo struct {
	Address string `json:"address"`
	Type    string `json:"type"`
	Reason  string `json:"reason"`
}

// TFCStateVersion represents a Terraform Cloud state version response
type TFCStateVersion struct {
	ID    string `json:"id"`
	State string `json:"state"` // Base64-encoded state JSON
}

// GenerateOptions holds options for import generation
type GenerateOptions struct {
	Module     string // Module filter (e.g., "module.network")
	Resource   string // Resource address filter (e.g., "aws_vpc.main")
	DryRun     bool
	OutputFile string
	Summary    string
	Verbose    bool
}

// AnalyzeOptions holds options for analysis
type AnalyzeOptions struct {
	Verbose bool
}

// ListOptions holds options for listing resources
type ListOptions struct {
	Module  string
	Type    string // For listing resource types
	Verbose bool
}

// ValidateOptions holds options for validation
type ValidateOptions struct {
	Verbose bool
}

// StateSource represents the source of a Terraform state (local or TFC)
type StateSource struct {
	Type         string // "local" or "tfc"
	LocalPath    string
	TFCOrg       string
	TFCWorkspace string
}

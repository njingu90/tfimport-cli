package pkg

import (
	"testing"
)

// TestIsSupportedAWSResource_DynamicSupport tests that all AWS resources are supported
func TestIsSupportedAWSResource_DynamicSupport(t *testing.T) {
	testCases := []struct {
		resourceType string
		expected     bool
		description  string
	}{
		// Standard AWS resources
		{"aws_vpc", true, "vpc is supported"},
		{"aws_subnet", true, "subnet is supported"},
		{"aws_instance", true, "instance is supported"},
		{"aws_s3_bucket", true, "s3_bucket is supported"},
		{"aws_lambda_function", true, "lambda_function is supported"},

		// Future AWS resources (not yet defined in schema)
		{"aws_new_future_service", true, "future AWS services are supported"},
		{"aws_another_service_resource", true, "any aws_* resource is supported"},

		// Non-AWS resources
		{"google_compute_instance", false, "google resources not supported"},
		{"azurerm_virtual_machine", false, "azure resources not supported"},
		{"aws", false, "incomplete aws_ prefix not supported"},
	}

	for _, tc := range testCases {
		t.Run(tc.resourceType, func(t *testing.T) {
			result := IsSupportedAWSResource(tc.resourceType)
			if result != tc.expected {
				t.Errorf("%s: got %v, expected %v", tc.description, result, tc.expected)
			}
		})
	}
}

// TestGetImportID_CommonResources tests ID extraction for common resources
func TestGetImportID_CommonResources(t *testing.T) {
	tests := []struct {
		name      string
		resource  Resource
		expectID  string
		expectErr bool
	}{
		{
			name: "aws_vpc with id",
			resource: Resource{
				Type: "aws_vpc",
				Attributes: map[string]interface{}{
					"id": "vpc-123456",
				},
			},
			expectID:  "vpc-123456",
			expectErr: false,
		},
		{
			name: "aws_instance with id",
			resource: Resource{
				Type: "aws_instance",
				Attributes: map[string]interface{}{
					"id": "i-1234567890abcdef",
				},
			},
			expectID:  "i-1234567890abcdef",
			expectErr: false,
		},
		{
			name: "aws_iam_role with name",
			resource: Resource{
				Type: "aws_iam_role",
				Attributes: map[string]interface{}{
					"name": "my-role",
				},
			},
			expectID:  "my-role",
			expectErr: false,
		},
		{
			name: "aws_s3_bucket with id",
			resource: Resource{
				Type: "aws_s3_bucket",
				Attributes: map[string]interface{}{
					"id": "my-bucket-name",
				},
			},
			expectID:  "my-bucket-name",
			expectErr: false,
		},
		{
			name: "aws_lambda_function with function_name",
			resource: Resource{
				Type: "aws_lambda_function",
				Attributes: map[string]interface{}{
					"function_name": "my-function",
				},
			},
			expectID:  "my-function",
			expectErr: false,
		},
		{
			name: "aws_iam_policy with arn",
			resource: Resource{
				Type: "aws_iam_policy",
				Attributes: map[string]interface{}{
					"arn": "arn:aws:iam::123456789012:policy/MyPolicy",
				},
			},
			expectID:  "arn:aws:iam::123456789012:policy/MyPolicy",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := GetImportID(tt.resource)
			if (err != nil) != tt.expectErr {
				t.Errorf("error = %v, expectErr = %v", err, tt.expectErr)
			}
			if id != tt.expectID {
				t.Errorf("id = %s, expectID = %s", id, tt.expectID)
			}
		})
	}
}

// TestGetImportID_CompoundIDs tests compound ID generation
func TestGetImportID_CompoundIDs(t *testing.T) {
	tests := []struct {
		name      string
		resource  Resource
		expectID  string
		expectErr bool
	}{
		{
			name: "aws_iam_role_policy_attachment",
			resource: Resource{
				Type: "aws_iam_role_policy_attachment",
				Attributes: map[string]interface{}{
					"role":       "my-role",
					"policy_arn": "arn:aws:iam::aws:policy/ReadOnlyAccess",
				},
			},
			expectID:  "my-role/arn:aws:iam::aws:policy/ReadOnlyAccess",
			expectErr: false,
		},
		{
			name: "aws_route_table_association with subnet",
			resource: Resource{
				Type: "aws_route_table_association",
				Attributes: map[string]interface{}{
					"route_table_id": "rtb-12345678",
					"subnet_id":      "subnet-87654321",
				},
			},
			expectID:  "subnet-87654321/rtb-12345678",
			expectErr: false,
		},
		{
			name: "aws_route with destination_cidr_block",
			resource: Resource{
				Type: "aws_route",
				Attributes: map[string]interface{}{
					"route_table_id":         "rtb-12345678",
					"destination_cidr_block": "10.0.0.0/16",
				},
			},
			expectID:  "rtb-12345678_10.0.0.0/16",
			expectErr: false,
		},
		{
			name: "aws_lambda_alias with compound",
			resource: Resource{
				Type: "aws_lambda_alias",
				Attributes: map[string]interface{}{
					"function_name": "my-function",
					"name":          "prod",
				},
			},
			expectID:  "my-function:prod",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := GetImportID(tt.resource)
			if (err != nil) != tt.expectErr {
				t.Errorf("error = %v, expectErr = %v", err, tt.expectErr)
			}
			if id != tt.expectID {
				t.Errorf("id = %s, expectID = %s", id, tt.expectID)
			}
		})
	}
}

// TestGetImportID_IntelligentFallback tests fallback for unmapped resources
func TestGetImportID_IntelligentFallback(t *testing.T) {
	tests := []struct {
		name      string
		resource  Resource
		expectID  string
		expectErr bool
	}{
		{
			name: "unmapped resource with id",
			resource: Resource{
				Type: "aws_future_service_unknown",
				Attributes: map[string]interface{}{
					"id": "fs-12345",
				},
			},
			expectID:  "fs-12345",
			expectErr: false,
		},
		{
			name: "unmapped resource with arn",
			resource: Resource{
				Type: "aws_another_new_service",
				Attributes: map[string]interface{}{
					"arn": "arn:aws:service:region:account:resource",
				},
			},
			expectID:  "arn:aws:service:region:account:resource",
			expectErr: false,
		},
		{
			name: "unmapped resource with name",
			resource: Resource{
				Type: "aws_yet_another_service",
				Attributes: map[string]interface{}{
					"name": "my-resource",
				},
			},
			expectID:  "my-resource",
			expectErr: false,
		},
		{
			name: "unmapped resource with no recognizable attributes",
			resource: Resource{
				Type: "aws_mystery_service",
				Attributes: map[string]interface{}{
					"custom_field": "value",
				},
			},
			expectID:  "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := GetImportID(tt.resource)
			if (err != nil) != tt.expectErr {
				t.Errorf("error = %v, expectErr = %v", err, tt.expectErr)
			}
			if id != tt.expectID {
				t.Errorf("id = %s, expectID = %s", id, tt.expectID)
			}
		})
	}
}

// TestGetImportID_NonAWSResource tests rejection of non-AWS resources
func TestGetImportID_NonAWSResource(t *testing.T) {
	tests := []struct {
		name         string
		resourceType string
	}{
		{"google compute", "google_compute_instance"},
		{"azure resource", "azurerm_virtual_machine"},
		{"terraform module", "module.example"},
		{"local resource", "local.var"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := Resource{
				Type: tt.resourceType,
				Attributes: map[string]interface{}{
					"id": "test-id",
				},
			}
			_, err := GetImportID(resource)
			if err == nil {
				t.Errorf("expected error for non-AWS resource %s", tt.resourceType)
			}
		})
	}
}

// TestValidateResourceForImport tests resource validation
func TestValidateResourceForImport(t *testing.T) {
	tests := []struct {
		name      string
		resource  Resource
		expectErr bool
	}{
		{
			name: "valid aws_vpc",
			resource: Resource{
				Type: "aws_vpc",
				Attributes: map[string]interface{}{
					"id": "vpc-123",
				},
			},
			expectErr: false,
		},
		{
			name: "invalid - missing id",
			resource: Resource{
				Type: "aws_vpc",
				Attributes: map[string]interface{}{
					"other": "field",
				},
			},
			expectErr: true,
		},
		{
			name: "invalid - non-AWS resource",
			resource: Resource{
				Type: "google_compute_instance",
				Attributes: map[string]interface{}{
					"id": "test",
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateResourceForImport(tt.resource)
			if (err != nil) != tt.expectErr {
				t.Errorf("error = %v, expectErr = %v", err, tt.expectErr)
			}
		})
	}
}

// TestGetUnsupportedReason tests unsupported resource messages
func TestGetUnsupportedReason(t *testing.T) {
	tests := []struct {
		resourceType string
		hasReason    bool
	}{
		{"aws_vpc", false},
		{"aws_future_service", false},
		{"google_compute", true},
		{"azure_vm", true},
	}

	for _, tt := range tests {
		t.Run(tt.resourceType, func(t *testing.T) {
			reason := GetUnsupportedReason(tt.resourceType)
			hasReason := reason != ""
			if hasReason != tt.hasReason {
				t.Errorf("hasReason = %v, expectHasReason = %v", hasReason, tt.hasReason)
			}
		})
	}
}

// TestRouteTableAssociationVariants tests different route table association patterns
func TestRouteTableAssociationVariants(t *testing.T) {
	tests := []struct {
		name      string
		resource  Resource
		expectID  string
		expectErr bool
	}{
		{
			name: "subnet association",
			resource: Resource{
				Type: "aws_route_table_association",
				Attributes: map[string]interface{}{
					"route_table_id": "rtb-12345678",
					"subnet_id":      "subnet-87654321",
				},
			},
			expectID:  "subnet-87654321/rtb-12345678",
			expectErr: false,
		},
		{
			name: "main route table (gateway)",
			resource: Resource{
				Type: "aws_route_table_association",
				Attributes: map[string]interface{}{
					"route_table_id": "rtb-main",
					"gateway_id":     "igw-12345678",
				},
			},
			expectID:  "igw-12345678/rtb-main",
			expectErr: false,
		},
		{
			name: "missing route_table_id",
			resource: Resource{
				Type: "aws_route_table_association",
				Attributes: map[string]interface{}{
					"subnet_id": "subnet-87654321",
				},
			},
			expectID:  "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := generateRouteTableAssociationID(tt.resource)
			if (err != nil) != tt.expectErr {
				t.Errorf("error = %v, expectErr = %v", err, tt.expectErr)
			}
			if id != tt.expectID {
				t.Errorf("id = %s, expectID = %s", id, tt.expectID)
			}
		})
	}
}

// TestRouteVariants tests different route patterns
func TestRouteVariants(t *testing.T) {
	tests := []struct {
		name      string
		resource  Resource
		expectID  string
		expectErr bool
	}{
		{
			name: "route with cidr",
			resource: Resource{
				Type: "aws_route",
				Attributes: map[string]interface{}{
					"route_table_id":         "rtb-12345678",
					"destination_cidr_block": "10.0.0.0/16",
				},
			},
			expectID:  "rtb-12345678_10.0.0.0/16",
			expectErr: false,
		},
		{
			name: "route with ipv6",
			resource: Resource{
				Type: "aws_route",
				Attributes: map[string]interface{}{
					"route_table_id":              "rtb-12345678",
					"destination_ipv6_cidr_block": "2001:db8::/32",
				},
			},
			expectID:  "rtb-12345678_2001:db8::/32",
			expectErr: false,
		},
		{
			name: "route with prefix list",
			resource: Resource{
				Type: "aws_route",
				Attributes: map[string]interface{}{
					"route_table_id":             "rtb-12345678",
					"destination_prefix_list_id": "pl-12345678",
				},
			},
			expectID:  "rtb-12345678_pl-12345678",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := generateRouteID(tt.resource)
			if (err != nil) != tt.expectErr {
				t.Errorf("error = %v, expectErr = %v", err, tt.expectErr)
			}
			if id != tt.expectID {
				t.Errorf("id = %s, expectID = %s", id, tt.expectID)
			}
		})
	}
}

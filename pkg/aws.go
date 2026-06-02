package pkg

import (
	"fmt"
	"strings"
)

// ResourceIDSchema defines how to extract import ID for a resource type
type ResourceIDSchema struct {
	// Simple ID attributes to try (in order)
	IDAttributes []string
	// Compound ID attributes with separator
	CompoundAttributes []string
	Separator          string
	// Custom generator function
	Generator func(Resource) (string, error)
}

// AWSDynamicResourceSchema provides universal support for all AWS resources
// Maps resource types to their import ID extraction patterns
var awsDynamicResourceSchema = map[string]ResourceIDSchema{
	// EC2 - Compute
	"aws_instance": {
		IDAttributes: []string{"id"},
	},
	"aws_ami": {
		IDAttributes: []string{"id"},
	},
	"aws_ami_copy": {
		IDAttributes: []string{"id"},
	},
	"aws_ami_from_container": {
		IDAttributes: []string{"id"},
	},
	"aws_ami_launch_permission": {
		CompoundAttributes: []string{"image_id", "account_id"},
		Separator:          ":",
	},
	"aws_ami_deprecation_notice": {
		IDAttributes: []string{"ami_id"},
	},

	// EC2 - VPC
	"aws_vpc": {
		IDAttributes: []string{"id"},
	},
	"aws_subnet": {
		IDAttributes: []string{"id"},
	},
	"aws_network_interface": {
		IDAttributes: []string{"id"},
	},
	"aws_network_interface_attachment": {
		CompoundAttributes: []string{"instance_id", "network_interface_id", "device_index"},
		Separator:          ":",
	},
	"aws_network_interface_sg_attachment": {
		CompoundAttributes: []string{"network_interface_id", "security_group_id"},
		Separator:          "_",
	},
	"aws_internet_gateway": {
		IDAttributes: []string{"id"},
	},
	"aws_internet_gateway_attachment": {
		CompoundAttributes: []string{"internet_gateway_id", "vpc_id"},
		Separator:          ":",
	},
	"aws_nat_gateway": {
		IDAttributes: []string{"id"},
	},
	"aws_eip": {
		IDAttributes: []string{"id"},
	},
	"aws_eip_association": {
		IDAttributes: []string{"id"},
	},
	"aws_route_table": {
		IDAttributes: []string{"id"},
	},
	"aws_route_table_association": {
		Generator: generateRouteTableAssociationID,
	},
	"aws_route": {
		Generator: generateRouteID,
	},
	"aws_vpc_dhcp_options_association": {
		IDAttributes: []string{"vpc_id"},
	},
	"aws_vpc_peering_connection": {
		IDAttributes: []string{"id"},
	},
	"aws_vpc_peering_connection_accepter": {
		IDAttributes: []string{"id"},
	},
	"aws_vpc_endpoint": {
		IDAttributes: []string{"id"},
	},
	"aws_vpc_endpoint_service": {
		IDAttributes: []string{"id"},
	},

	// Security Groups
	"aws_security_group": {
		IDAttributes: []string{"id"},
	},
	"aws_security_group_rule": {
		Generator: generateSecurityGroupRuleID,
	},

	// IAM
	"aws_iam_role": {
		IDAttributes: []string{"name"},
	},
	"aws_iam_role_policy": {
		CompoundAttributes: []string{"role", "name"},
		Separator:          ":",
	},
	"aws_iam_role_policy_attachment": {
		Generator: generateIAMRolePolicyAttachmentID,
	},
	"aws_iam_policy": {
		IDAttributes: []string{"arn"},
	},
	"aws_iam_policy_version": {
		IDAttributes: []string{"arn"},
	},
	"aws_iam_user": {
		IDAttributes: []string{"name"},
	},
	"aws_iam_user_policy": {
		CompoundAttributes: []string{"user", "name"},
		Separator:          ":",
	},
	"aws_iam_user_policy_attachment": {
		Generator: generateIAMUserPolicyAttachmentID,
	},
	"aws_iam_group": {
		IDAttributes: []string{"name"},
	},
	"aws_iam_group_policy": {
		CompoundAttributes: []string{"group", "name"},
		Separator:          ":",
	},
	"aws_iam_group_policy_attachment": {
		Generator: generateIAMGroupPolicyAttachmentID,
	},
	"aws_iam_instance_profile": {
		IDAttributes: []string{"name"},
	},
	"aws_iam_access_key": {
		IDAttributes: []string{"id"},
	},
	"aws_iam_saml_provider": {
		IDAttributes: []string{"arn"},
	},
	"aws_iam_openid_connect_provider": {
		IDAttributes: []string{"arn"},
	},

	// S3
	"aws_s3_bucket": {
		IDAttributes: []string{"id"},
	},
	"aws_s3_bucket_versioning": {
		IDAttributes: []string{"bucket"},
	},
	"aws_s3_bucket_acl": {
		IDAttributes: []string{"bucket"},
	},
	"aws_s3_bucket_cors_configuration": {
		IDAttributes: []string{"bucket"},
	},
	"aws_s3_bucket_lifecycle_configuration": {
		IDAttributes: []string{"bucket"},
	},
	"aws_s3_bucket_logging": {
		IDAttributes: []string{"bucket"},
	},
	"aws_s3_bucket_policy": {
		IDAttributes: []string{"bucket"},
	},
	"aws_s3_bucket_public_access_block": {
		IDAttributes: []string{"bucket"},
	},
	"aws_s3_bucket_replication_configuration": {
		IDAttributes: []string{"role"},
	},
	"aws_s3_object": {
		CompoundAttributes: []string{"bucket", "key"},
		Separator:          ",",
	},
	"aws_s3_object_copy": {
		CompoundAttributes: []string{"bucket", "key"},
		Separator:          ",",
	},

	// RDS
	"aws_db_instance": {
		IDAttributes: []string{"identifier"},
	},
	"aws_db_parameter_group": {
		IDAttributes: []string{"name"},
	},
	"aws_db_option_group": {
		IDAttributes: []string{"name"},
	},
	"aws_db_subnet_group": {
		IDAttributes: []string{"name"},
	},
	"aws_db_cluster": {
		IDAttributes: []string{"cluster_identifier"},
	},
	"aws_db_cluster_parameter_group": {
		IDAttributes: []string{"name"},
	},
	"aws_rds_cluster_instance": {
		IDAttributes: []string{"identifier"},
	},

	// Lambda
	"aws_lambda_function": {
		IDAttributes: []string{"function_name"},
	},
	"aws_lambda_alias": {
		CompoundAttributes: []string{"function_name", "name"},
		Separator:          ":",
	},
	"aws_lambda_layer_version": {
		CompoundAttributes: []string{"layer_name", "version"},
		Separator:          ":",
	},
	"aws_lambda_function_url": {
		IDAttributes: []string{"function_name"},
	},
	"aws_lambda_permission": {
		CompoundAttributes: []string{"function_name", "statement_id"},
		Separator:          "/",
	},

	// ECS
	"aws_ecs_cluster": {
		IDAttributes: []string{"name"},
	},
	"aws_ecs_service": {
		CompoundAttributes: []string{"cluster", "name"},
		Separator:          "/",
	},
	"aws_ecs_task_definition": {
		IDAttributes: []string{"arn"},
	},
	"aws_ecs_container_definition": {
		IDAttributes: []string{"container_name"},
	},

	// EKS
	"aws_eks_cluster": {
		IDAttributes: []string{"name"},
	},
	"aws_eks_node_group": {
		CompoundAttributes: []string{"cluster_name", "node_group_name"},
		Separator:          ":",
	},
	"aws_eks_addon": {
		CompoundAttributes: []string{"cluster_name", "addon_name"},
		Separator:          ":",
	},

	// Load Balancing
	"aws_lb": {
		IDAttributes: []string{"arn"},
	},
	"aws_lb_target_group": {
		IDAttributes: []string{"arn"},
	},
	"aws_lb_listener": {
		IDAttributes: []string{"arn"},
	},
	"aws_lb_listener_rule": {
		IDAttributes: []string{"arn"},
	},
	"aws_lb_target_group_attachment": {
		CompoundAttributes: []string{"target_group_arn", "target_id"},
		Separator:          "-",
	},
	"aws_alb": {
		IDAttributes: []string{"arn"},
	},
	"aws_alb_target_group": {
		IDAttributes: []string{"arn"},
	},
	"aws_alb_listener": {
		IDAttributes: []string{"arn"},
	},
	"aws_alb_listener_rule": {
		IDAttributes: []string{"arn"},
	},
	"aws_alb_target_group_attachment": {
		CompoundAttributes: []string{"target_group_arn", "target_id"},
		Separator:          "-",
	},
	"aws_elb": {
		IDAttributes: []string{"name"},
	},
	"aws_elb_attachment": {
		CompoundAttributes: []string{"load_balancer_name", "instance"},
		Separator:          "-",
	},

	// CloudWatch
	"aws_cloudwatch_log_group": {
		IDAttributes: []string{"name"},
	},
	"aws_cloudwatch_log_stream": {
		CompoundAttributes: []string{"log_group_name", "name"},
		Separator:          ":",
	},
	"aws_cloudwatch_log_resource_policy": {
		IDAttributes: []string{"policy_name"},
	},
	"aws_cloudwatch_log_data_protection_policy": {
		IDAttributes: []string{"log_group_name"},
	},
	"aws_cloudwatch_event_rule": {
		IDAttributes: []string{"name"},
	},
	"aws_cloudwatch_event_target": {
		CompoundAttributes: []string{"rule", "target_id"},
		Separator:          "-",
	},
	"aws_cloudwatch_metric_alarm": {
		IDAttributes: []string{"alarm_name"},
	},

	// KMS
	"aws_kms_key": {
		IDAttributes: []string{"id"},
	},
	"aws_kms_alias": {
		IDAttributes: []string{"name"},
	},
	"aws_kms_grant": {
		IDAttributes: []string{"grant_id"},
	},
	"aws_kms_key_policy": {
		IDAttributes: []string{"key_id"},
	},

	// Secrets Manager
	"aws_secretsmanager_secret": {
		IDAttributes: []string{"id"},
	},
	"aws_secretsmanager_secret_version": {
		IDAttributes: []string{"secret_id"},
	},

	// Route53
	"aws_route53_zone": {
		IDAttributes: []string{"zone_id"},
	},
	"aws_route53_record": {
		Generator: generateRoute53RecordID,
	},
	"aws_route53_health_check": {
		IDAttributes: []string{"id"},
	},
	"aws_route53_traffic_policy": {
		IDAttributes: []string{"id"},
	},

	// VPN
	"aws_vpn_gateway": {
		IDAttributes: []string{"id"},
	},
	"aws_vpn_connection": {
		IDAttributes: []string{"id"},
	},
	"aws_vpn_connection_route": {
		CompoundAttributes: []string{"vpn_connection_id", "destination_cidr_block"},
		Separator:          "_",
	},
	"aws_customer_gateway": {
		IDAttributes: []string{"id"},
	},

	// CloudFront
	"aws_cloudfront_distribution": {
		IDAttributes: []string{"id"},
	},
	"aws_cloudfront_origin_access_identity": {
		IDAttributes: []string{"id"},
	},
	"aws_cloudfront_cache_policy": {
		IDAttributes: []string{"id"},
	},
	"aws_cloudfront_origin_request_policy": {
		IDAttributes: []string{"id"},
	},

	// CloudFormation
	"aws_cloudformation_stack": {
		IDAttributes: []string{"name"},
	},
	"aws_cloudformation_stack_set": {
		IDAttributes: []string{"name"},
	},

	// SNS
	"aws_sns_topic": {
		IDAttributes: []string{"arn"},
	},
	"aws_sns_topic_policy": {
		IDAttributes: []string{"arn"},
	},
	"aws_sns_topic_subscription": {
		IDAttributes: []string{"arn"},
	},

	// SQS
	"aws_sqs_queue": {
		IDAttributes: []string{"url"},
	},
	"aws_sqs_queue_policy": {
		IDAttributes: []string{"queue_url"},
	},

	// DynamoDB
	"aws_dynamodb_table": {
		IDAttributes: []string{"name"},
	},
	"aws_dynamodb_table_item": {
		CompoundAttributes: []string{"table_name", "hash_key"},
		Separator:          ",",
	},
	"aws_dynamodb_global_table": {
		IDAttributes: []string{"name"},
	},
	"aws_dynamodb_ttl": {
		IDAttributes: []string{"table_name"},
	},

	// Kinesis
	"aws_kinesis_stream": {
		IDAttributes: []string{"name"},
	},
	"aws_kinesis_firehose_delivery_stream": {
		IDAttributes: []string{"name"},
	},

	// ElastiCache
	"aws_elasticache_cluster": {
		IDAttributes: []string{"cluster_id"},
	},
	"aws_elasticache_parameter_group": {
		IDAttributes: []string{"name"},
	},
	"aws_elasticache_replication_group": {
		IDAttributes: []string{"replication_group_id"},
	},

	// Elasticsearch
	"aws_elasticsearch_domain": {
		IDAttributes: []string{"domain_name"},
	},
	"aws_opensearch_domain": {
		IDAttributes: []string{"domain_name"},
	},

	// ACM
	"aws_acm_certificate": {
		IDAttributes: []string{"arn"},
	},

	// API Gateway
	"aws_api_gateway_rest_api": {
		IDAttributes: []string{"id"},
	},
	"aws_api_gateway_resource": {
		IDAttributes: []string{"id"},
	},
	"aws_api_gateway_method": {
		CompoundAttributes: []string{"rest_api_id", "resource_id", "http_method"},
		Separator:          "-",
	},
	"aws_api_gateway_integration": {
		CompoundAttributes: []string{"rest_api_id", "resource_id", "http_method"},
		Separator:          "-",
	},
	"aws_api_gateway_deployment": {
		CompoundAttributes: []string{"rest_api_id", "id"},
		Separator:          "_",
	},
	"aws_api_gateway_stage": {
		CompoundAttributes: []string{"rest_api_id", "stage_name"},
		Separator:          "_",
	},
	"aws_apigatewayv2_api": {
		IDAttributes: []string{"id"},
	},
	"aws_apigatewayv2_stage": {
		CompoundAttributes: []string{"api_id", "name"},
		Separator:          "/",
	},

	// Backup
	"aws_backup_vault": {
		IDAttributes: []string{"name"},
	},
	"aws_backup_plan": {
		IDAttributes: []string{"arn"},
	},

	// SSM
	"aws_ssm_parameter": {
		IDAttributes: []string{"name"},
	},
	"aws_ssm_document": {
		IDAttributes: []string{"name"},
	},

	// AppConfig
	"aws_appconfig_app": {
		IDAttributes: []string{"id"},
	},
	"aws_appconfig_configuration_profile": {
		IDAttributes: []string{"configuration_profile_id"},
	},
	"aws_appconfig_environment": {
		IDAttributes: []string{"environment_id"},
	},

	// CodeBuild
	"aws_codebuild_project": {
		IDAttributes: []string{"name"},
	},

	// CodePipeline
	"aws_codepipeline": {
		IDAttributes: []string{"name"},
	},

	// Default fallback for any unmapped AWS resource
	"aws_default": {
		IDAttributes: []string{"id", "arn", "name"},
	},
}

// IsSupportedAWSResource checks if a resource is an AWS resource (all aws_* resources are supported)
func IsSupportedAWSResource(resourceType string) bool {
	return strings.HasPrefix(resourceType, "aws_")
}

// GetUnsupportedReason returns the reason why a resource type is not supported
func GetUnsupportedReason(resourceType string) string {
	if IsSupportedAWSResource(resourceType) {
		return ""
	}
	return fmt.Sprintf("resource type %s is not an AWS resource (must start with 'aws_')", resourceType)
}

// GetImportID generates the import ID for any AWS resource
func GetImportID(res Resource) (string, error) {
	if !IsSupportedAWSResource(res.Type) {
		return "", fmt.Errorf("resource type %s is not an AWS resource", res.Type)
	}

	// Check if we have a schema for this specific resource type
	if schema, ok := awsDynamicResourceSchema[res.Type]; ok {
		// Use custom generator if defined
		if schema.Generator != nil {
			return schema.Generator(res)
		}

		// Use compound attributes if defined
		if len(schema.CompoundAttributes) > 0 {
			return generateCompoundID(res, schema.CompoundAttributes, schema.Separator)
		}

		// Use simple ID attributes
		if len(schema.IDAttributes) > 0 {
			return generateSimpleID(res, schema.IDAttributes)
		}
	}

	// Fallback: Try default schema (uses common attributes)
	return generateWithIntelligentFallback(res)
}

// Helper function to get string attributes
func getStringAttr(attrs map[string]interface{}, key string) string {
	if val, ok := attrs[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// Helper function to get interface attributes (for nested objects)
func getInterfaceAttr(attrs map[string]interface{}, key string) interface{} {
	return attrs[key]
}

// Helper function to get first element from array attribute
func getFirstArrayElement(attrs map[string]interface{}, key string) string {
	if val, ok := attrs[key]; ok {
		if arr, ok := val.([]interface{}); ok && len(arr) > 0 {
			if str, ok := arr[0].(string); ok {
				return str
			}
		}
	}
	return ""
}

// Helper function to get all elements from array attribute joined by separator
func getAllArrayElements(attrs map[string]interface{}, key string, separator string) string {
	if val, ok := attrs[key]; ok {
		if arr, ok := val.([]interface{}); ok && len(arr) > 0 {
			var elements []string
			for _, item := range arr {
				if str, ok := item.(string); ok {
					elements = append(elements, str)
				}
			}
			if len(elements) > 0 {
				return strings.Join(elements, separator)
			}
		}
	}
	return ""
}

// Helper function to convert number to string
func getNumberAsString(attrs map[string]interface{}, key string) string {
	if val, ok := attrs[key]; ok {
		switch v := val.(type) {
		case float64:
			return fmt.Sprintf("%.0f", v)
		case int:
			return fmt.Sprintf("%d", v)
		case string:
			return v
		}
	}
	return ""
}

// generateSimpleID extracts ID from simple attributes
func generateSimpleID(res Resource, attributes []string) (string, error) {
	for _, attr := range attributes {
		id := getStringAttr(res.Attributes, attr)
		if id != "" {
			return id, nil
		}
	}
	return "", fmt.Errorf("could not find ID from attributes %v for %s", attributes, res.Type)
}

// generateCompoundID creates ID from multiple attributes
func generateCompoundID(res Resource, attributes []string, separator string) (string, error) {
	var parts []string
	for _, attr := range attributes {
		val := getStringAttr(res.Attributes, attr)
		if val == "" {
			return "", fmt.Errorf("missing required attribute '%s' for compound ID in %s", attr, res.Type)
		}
		parts = append(parts, val)
	}
	return strings.Join(parts, separator), nil
}

// generateWithIntelligentFallback tries common attribute patterns for unmapped resources
func generateWithIntelligentFallback(res Resource) (string, error) {
	// Try standard ID attributes in order
	commonIDAttrs := []string{"id", "arn", "name", "identifier", "url"}
	for _, attr := range commonIDAttrs {
		if val := getStringAttr(res.Attributes, attr); val != "" {
			return val, nil
		}
	}

	// If still not found, return error with suggestions
	keys := make([]string, 0, len(res.Attributes))
	for k := range res.Attributes {
		keys = append(keys, k)
	}
	return "", fmt.Errorf("could not determine import ID for %s; available attributes: %v", res.Type, keys)
}

// generateRouteTableAssociationID generates import ID for aws_route_table_association
// Format: subnet-id/rtb-id or igw-id/rtb-id (per Terraform documentation)
func generateRouteTableAssociationID(res Resource) (string, error) {
	rtID := getStringAttr(res.Attributes, "route_table_id")
	subnetID := getStringAttr(res.Attributes, "subnet_id")
	gatewayID := getStringAttr(res.Attributes, "gateway_id")

	if rtID == "" {
		return "", fmt.Errorf("missing route_table_id")
	}

	if subnetID != "" {
		return fmt.Sprintf("%s/%s", subnetID, rtID), nil
	}

	if gatewayID != "" {
		return fmt.Sprintf("%s/%s", gatewayID, rtID), nil
	}

	return "", fmt.Errorf("missing both subnet_id and gateway_id")
}

// generateSecurityGroupRuleID generates import ID for aws_security_group_rule
// Format: sg-xxxxx_ingress_tcp_80_80_10.0.0.0/8 (per Terraform documentation)
func generateSecurityGroupRuleID(res Resource) (string, error) {
	sgID := getStringAttr(res.Attributes, "security_group_id")
	ruleType := getStringAttr(res.Attributes, "type") // "ingress" or "egress"

	if sgID == "" {
		return "", fmt.Errorf("missing security_group_id")
	}

	if ruleType == "" {
		return "", fmt.Errorf("missing rule type")
	}

	// Get protocol
	protocol := getStringAttr(res.Attributes, "protocol")
	if protocol == "" {
		protocol = "-1" // Default for all protocols
	}

	// Get ports (can be numbers or strings)
	fromPort := getNumberAsString(res.Attributes, "from_port")
	toPort := getNumberAsString(res.Attributes, "to_port")

	if fromPort == "" {
		fromPort = "0"
	}
	if toPort == "" {
		toPort = "0"
	}

	// Get CIDR or referenced security group
	// cidr_blocks and ipv6_cidr_blocks are arrays, get ALL elements joined by underscore
	cidr := getAllArrayElements(res.Attributes, "cidr_blocks", "_")
	if cidr == "" {
		cidr = getAllArrayElements(res.Attributes, "ipv6_cidr_blocks", "_")
	}
	if cidr == "" {
		// Try source_security_group_id or referenced_security_group_id
		cidr = getStringAttr(res.Attributes, "source_security_group_id")
		if cidr == "" {
			cidr = getStringAttr(res.Attributes, "referenced_security_group_id")
		}
	}

	if cidr == "" {
		return "", fmt.Errorf("missing CIDR or referenced security group")
	}

	// Format: sg-xxxxx_ingress_tcp_22_22_10.0.0.0/8_172.16.0.0/12_192.168.0.0/16
	// All CIDR blocks are joined with underscores
	return fmt.Sprintf("%s_%s_%s_%s_%s_%s", sgID, ruleType, protocol, fromPort, toPort, cidr), nil
}

// generateIAMRolePolicyAttachmentID generates import ID for aws_iam_role_policy_attachment
// Format: role_name/policy_arn
func generateIAMRolePolicyAttachmentID(res Resource) (string, error) {
	roleName := getStringAttr(res.Attributes, "role")
	policyARN := getStringAttr(res.Attributes, "policy_arn")

	if roleName == "" {
		return "", fmt.Errorf("missing role name")
	}

	if policyARN == "" {
		return "", fmt.Errorf("missing policy_arn")
	}

	return fmt.Sprintf("%s/%s", roleName, policyARN), nil
}

// generateIAMUserPolicyAttachmentID generates import ID for aws_iam_user_policy_attachment
// Format: user_name/policy_arn
func generateIAMUserPolicyAttachmentID(res Resource) (string, error) {
	userName := getStringAttr(res.Attributes, "user")
	policyARN := getStringAttr(res.Attributes, "policy_arn")

	if userName == "" {
		return "", fmt.Errorf("missing user name")
	}

	if policyARN == "" {
		return "", fmt.Errorf("missing policy_arn")
	}

	return fmt.Sprintf("%s/%s", userName, policyARN), nil
}

// generateIAMGroupPolicyAttachmentID generates import ID for aws_iam_group_policy_attachment
// Format: group_name/policy_arn
func generateIAMGroupPolicyAttachmentID(res Resource) (string, error) {
	groupName := getStringAttr(res.Attributes, "group")
	policyARN := getStringAttr(res.Attributes, "policy_arn")

	if groupName == "" {
		return "", fmt.Errorf("missing group name")
	}

	if policyARN == "" {
		return "", fmt.Errorf("missing policy_arn")
	}

	return fmt.Sprintf("%s/%s", groupName, policyARN), nil
}

// generateRouteID generates import ID for aws_route
// Format: rtb-id_destination (underscore separator per Terraform documentation)
func generateRouteID(res Resource) (string, error) {
	rtID := getStringAttr(res.Attributes, "route_table_id")

	if rtID == "" {
		return "", fmt.Errorf("missing route_table_id")
	}

	// Check for destination CIDR
	destCIDR := getStringAttr(res.Attributes, "destination_cidr_block")
	if destCIDR != "" {
		return fmt.Sprintf("%s_%s", rtID, destCIDR), nil
	}

	// Check for destination IPv6 CIDR
	destIPv6CIDR := getStringAttr(res.Attributes, "destination_ipv6_cidr_block")
	if destIPv6CIDR != "" {
		return fmt.Sprintf("%s_%s", rtID, destIPv6CIDR), nil
	}

	// Check for prefix list
	destPrefixList := getStringAttr(res.Attributes, "destination_prefix_list_id")
	if destPrefixList != "" {
		return fmt.Sprintf("%s_%s", rtID, destPrefixList), nil
	}

	return "", fmt.Errorf("missing destination CIDR block, IPv6 CIDR, or prefix list")
}

// generateRoute53RecordID generates import ID for aws_route53_record
// Format: zone_id_record_name_record_type
func generateRoute53RecordID(res Resource) (string, error) {
	zoneID := getStringAttr(res.Attributes, "zone_id")
	fqdn := getStringAttr(res.Attributes, "fqdn")
	recordType := getStringAttr(res.Attributes, "type")

	if zoneID == "" {
		return "", fmt.Errorf("missing zone_id")
	}

	if fqdn == "" {
		return "", fmt.Errorf("missing fqdn")
	}

	if recordType == "" {
		return "", fmt.Errorf("missing record type")
	}

	// Format: Z1234567890ABC_example.com_A
	// Normalize FQDN by ensuring trailing dot
	if !strings.HasSuffix(fqdn, ".") {
		fqdn = fqdn + "."
	}

	return fmt.Sprintf("%s_%s_%s", zoneID, fqdn, recordType), nil
}

// ValidateResourceForImport checks if a resource has all required attributes for import
func ValidateResourceForImport(res Resource) error {
	if !IsSupportedAWSResource(res.Type) {
		return fmt.Errorf("resource type %s is not an AWS resource", res.Type)
	}

	_, err := GetImportID(res)
	return err
}

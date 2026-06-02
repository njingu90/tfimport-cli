package pkg

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	// TerraformCloudAPI is the base URL for Terraform Cloud API
	TerraformCloudAPI = "https://app.terraform.io/api/v2"
)

// FetchStateFromTFC retrieves the latest state version from Terraform Cloud
func FetchStateFromTFC(org, workspace, token string) (*TerraformState, error) {
	if org == "" || workspace == "" {
		return nil, fmt.Errorf("organization and workspace are required")
	}

	if token == "" {
		// Try to get token from environment
		token = os.Getenv("TF_API_TOKEN")
		if token == "" {
			return nil, fmt.Errorf("Terraform Cloud API token not provided and TF_API_TOKEN not set")
		}
	}

	// Get current state version
	stateVersion, err := getStateVersion(org, workspace, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get state version: %w", err)
	}

	if stateVersion == nil || stateVersion.State == "" {
		return nil, fmt.Errorf("no state found for organization %s, workspace %s", org, workspace)
	}

	// Try to parse as-is first (might be plain JSON from download URL)
	var state TerraformState
	err = json.Unmarshal([]byte(stateVersion.State), &state)
	if err == nil && state.Version > 0 {
		// Successfully parsed as JSON
		return &state, nil
	}

	// If that failed, try base64 decoding first
	stateData, decodeErr := base64.StdEncoding.DecodeString(stateVersion.State)
	if decodeErr == nil {
		if err := json.Unmarshal(stateData, &state); err != nil {
			return nil, fmt.Errorf("failed to parse state JSON: %w", err)
		}
		return &state, nil
	}

	// If both approaches failed, return the original parse error
	return nil, fmt.Errorf("failed to parse state JSON: %w", err)
}

// getWorkspaceID retrieves the workspace ID from its name
func getWorkspaceID(org, workspace, token string) (string, error) {
	// Endpoint: GET /organizations/:org_name/workspaces/:workspace_name
	url := fmt.Sprintf("%s/organizations/%s/workspaces/%s", TerraformCloudAPI, org, workspace)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/vnd.api+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch workspace info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("authentication failed: invalid token")
	}

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("workspace not found: %s/%s", org, workspace)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var response struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse API response: %w", err)
	}

	if response.Data.ID == "" {
		return "", fmt.Errorf("workspace ID not found in response")
	}

	return response.Data.ID, nil
}

// getStateVersion retrieves the current state version from TFC API
func getStateVersion(org, workspace, token string) (*TFCStateVersion, error) {
	// First, get the workspace ID by querying the workspace info
	workspaceID, err := getWorkspaceID(org, workspace, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace ID: %w", err)
	}

	// Now use the workspace ID to get the current state version metadata
	// Endpoint: GET /workspaces/:workspace_id/current-state-version
	url := fmt.Sprintf("%s/workspaces/%s/current-state-version", TerraformCloudAPI, workspaceID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/vnd.api+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch state version: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("authentication failed: invalid token")
	}

	if resp.StatusCode == http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("workspace not found: %s/%s (URL: %s, response: %s)", org, workspace, url, string(body))
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d (URL: %s): %s", resp.StatusCode, url, string(body))
	}

	// Parse response (TFC returns JSONAPI format)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response struct {
		Data struct {
			ID         string `json:"id"`
			Attributes struct {
				HostedStateDownloadURL string `json:"hosted-state-download-url"`
			} `json:"attributes"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	if response.Data.Attributes.HostedStateDownloadURL == "" {
		return nil, fmt.Errorf("no state download URL found in response")
	}

	// Now download the actual state from the hosted URL (not hosted-json-state, which is filtered)
	stateData, err := downloadStateFromURL(response.Data.Attributes.HostedStateDownloadURL, token)
	if err != nil {
		return nil, fmt.Errorf("failed to download state: %w", err)
	}

	return &TFCStateVersion{
		ID:    response.Data.ID,
		State: stateData,
	}, nil
}

// downloadStateFromURL downloads state from a hosted URL
// The download URL requires Bearer token authentication and may redirect
func downloadStateFromURL(downloadURL, token string) (string, error) {
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set Authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Create client with custom redirect policy to preserve auth header
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Preserve the Authorization header across redirects
			if len(via) > 0 && via[0].Header.Get("Authorization") != "" {
				req.Header.Set("Authorization", via[0].Header.Get("Authorization"))
			}
			return nil
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download state: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to download state (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read state data: %w", err)
	}

	// Return the state as a JSON string
	return string(body), nil
}

// ValidateTFCToken validates that a Terraform Cloud token is valid
func ValidateTFCToken(token string) error {
	if token == "" {
		token = os.Getenv("TF_API_TOKEN")
		if token == "" {
			return fmt.Errorf("no Terraform Cloud token provided")
		}
	}

	if !strings.HasPrefix(token, "skpenv-") && !strings.HasPrefix(token, "atlasv1.") {
		return fmt.Errorf("invalid token format (should start with 'skpenv-' or 'atlasv1.')")
	}

	return nil
}

// GetWorkspaceMetadata retrieves metadata about a workspace
func GetWorkspaceMetadata(org, workspace, token string) (map[string]interface{}, error) {
	if org == "" || workspace == "" {
		return nil, fmt.Errorf("organization and workspace are required")
	}

	if token == "" {
		token = os.Getenv("TF_API_TOKEN")
		if token == "" {
			return nil, fmt.Errorf("no Terraform Cloud token provided")
		}
	}

	url := fmt.Sprintf("%s/organizations/%s/workspaces/%s", TerraformCloudAPI, org, workspace)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/vnd.api+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch workspace metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch workspace metadata (status: %d)", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return response, nil
}

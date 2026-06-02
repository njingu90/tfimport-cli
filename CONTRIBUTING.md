# Contributing to tfimport-cli

Thank you for your interest in contributing to tfimport-cli! This document provides guidelines and instructions for contributing.

## Code of Conduct

Be respectful and inclusive. We're building a community that welcomes all contributors.

## Development Setup

### Prerequisites
- Go 1.21 or later
- Git
- Make

### Clone and Setup

```bash
git clone https://github.com/njingu90/tfimport-cli.git
cd tfimport-cli
go mod download
make dev
```

### Project Structure

```
tfimport-cli/
├── cmd/tfimport-cli/           # CLI entry point
├── pkg/                         # Core packages (state, filter, aws, tfc, importgen, report)
├── internal/cmd/               # Command implementations (analyze, generate, list, validate, version)
├── testdata/                    # Test fixtures (sample state files)
├── docs/                        # Documentation
├── .github/workflows/          # CI/CD pipelines
└── Makefile                     # Build targets
```

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Your Changes

- Follow Go conventions
- Write tests for new functionality
- Keep commits focused and descriptive

### 3. Run Tests Locally

```bash
make test          # Run all tests with coverage
make lint          # Run linters
make build         # Build the binary
```

### 4. Verify Code Quality

```bash
go fmt ./...       # Format code
go vet ./...       # Check for issues
```

### 5. Commit and Push

```bash
git commit -m "Brief description of changes"
git push origin feature/your-feature-name
```

### 6. Submit a Pull Request

- Provide a clear description of changes
- Reference any related issues
- Ensure all tests pass

## Testing Requirements

### Unit Tests

Every package should have corresponding `*_test.go` files:

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/...
```

### Test Coverage

- Minimum 75% code coverage required
- Test both success and error cases
- Test edge cases (empty inputs, missing fields, etc.)

### Integration Tests

Test the full workflow:

```bash
# Test with sample data
./tfimport-cli analyze --state testdata/sample.tfstate
./tfimport-cli generate --state testdata/sample.tfstate
./tfimport-cli validate --state testdata/sample.tfstate
```

## Adding Support for New AWS Resources

To add a new supported AWS resource:

1. **Add to supported list** in `pkg/aws.go`:
   ```go
   var supportedAWSResources = map[string]bool{
       "aws_your_resource": true,
       // ...
   }
   ```

2. **Implement ID generator** in `pkg/aws.go`:
   ```go
   func generateYourResourceID(res Resource) (string, error) {
       // Extract ID from attributes
       id := getStringAttr(res.Attributes, "id")
       if id == "" {
           return "", fmt.Errorf("missing id attribute")
       }
       return id, nil
   }
   ```

3. **Add to GetImportID switch** in `pkg/aws.go`:
   ```go
   case "aws_your_resource":
       return generateYourResourceID(res)
   ```

4. **Write tests** in `pkg/aws_test.go`:
   ```go
   func TestGetImportID_YourResource(t *testing.T) {
       // Test implementation
   }
   ```

5. **Update documentation** in README.md:
   - Add to supported resources list
   - Add to command examples if relevant

## Adding Support for New Providers

To add support for a new cloud provider (e.g., Azure, GCP):

1. **Create new package** `pkg/azure.go` or `pkg/gcp.go`

2. **Implement provider interface**:
   - `IsSupported(resourceType string) bool`
   - `GetImportID(resource Resource) (string, error)`
   - Provider-specific ID generators

3. **Update state filters** to handle provider resources

4. **Update command logic** to route to correct provider

5. **Add tests** for new provider

6. **Update documentation** with new provider support

## Code Style

### Go Conventions

- Follow standard Go style guidelines
- Use `gofmt` for formatting
- Use meaningful variable names
- Keep functions focused and small

### Naming

- `CamelCase` for functions and variables
- `snake_case` for file names
- UPPERCASE for constants

### Comments

- Document exported functions and types
- Explain complex logic
- Keep comments concise

Example:
```go
// GetImportID generates the import ID for an AWS resource.
// It routes to provider-specific ID generators based on resource type.
func GetImportID(res Resource) (string, error) {
    // Implementation
}
```

## Documentation

### Update README.md

If your changes affect user-facing behavior:
- Update command descriptions
- Update usage examples
- Update feature list

### Update INSTALL.md

If installation process changes:
- Update platform-specific instructions
- Update troubleshooting section

### Add Code Comments

Explain non-obvious logic:
```go
// Compound IDs require combining multiple attributes
// Format: subnet_id:route_table_id or gateway_id:route_table_id
func generateRouteTableAssociationID(res Resource) (string, error) {
    // Implementation
}
```

## Performance Considerations

- Avoid unnecessary allocations
- Use efficient data structures
- Handle large state files gracefully
- No heavy external dependencies

## Security

- Never log sensitive data (tokens, IDs)
- Validate all inputs
- Handle errors gracefully
- Don't modify state or infrastructure

## Bug Reports

When reporting issues:
1. Describe the problem clearly
2. Provide steps to reproduce
3. Include tfimport-cli version: `tfimport-cli version`
4. Include state file snippet (sanitized)
5. Include error messages and logs

## Feature Requests

For feature requests:
1. Describe the use case
2. Explain expected behavior
3. Provide examples if possible
4. Discuss implementation approach

## Reviewing PRs

When reviewing others' PRs:
- Check code quality and style
- Verify tests are included
- Verify documentation is updated
- Run the code locally if possible
- Provide constructive feedback

## Release Process

Releases are automated via GitHub Actions:

1. Create a git tag: `git tag v1.0.0`
2. Push the tag: `git push origin v1.0.0`
3. GitHub Actions builds binaries
4. Release is created automatically

Version format: `vMAJOR.MINOR.PATCH` (semantic versioning)

## Questions?

- Open an issue on GitHub
- Check existing issues/discussions
- Read the README and documentation

## Thank You!

Thank you for contributing to tfimport-cli. Your efforts help make infrastructure migration easier for everyone!

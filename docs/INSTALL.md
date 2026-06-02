# Installation Guide

## Platform-Specific Installation

### Linux

#### Using Downloaded Binary

1. **Download the release**:
   ```bash
   VERSION=v1.0.0
   curl -L https://github.com/njingu90/tfimport-cli/releases/download/${VERSION}/tfimport-cli-${VERSION}-linux-amd64 \
     -o tfimport-cli
   ```

2. **Verify checksum** (optional but recommended):
   ```bash
   curl -L https://github.com/njingu90/tfimport-cli/releases/download/${VERSION}/SHA256SUMS \
     -o SHA256SUMS
   sha256sum -c SHA256SUMS | grep tfimport-cli-${VERSION}-linux-amd64
   ```

3. **Make executable**:
   ```bash
   chmod +x tfimport-cli
   ```

4. **Move to PATH**:
   ```bash
   sudo mv tfimport-cli /usr/local/bin/
   ```

5. **Verify installation**:
   ```bash
   tfimport-cli --version
   ```

#### Building from Source

```bash
git clone https://github.com/njingu90/tfimport-cli.git
cd tfimport-cli
make build
sudo cp bin/tfimport-cli /usr/local/bin/
tfimport-cli --version
```

### macOS

#### Using Downloaded Binary (Intel)

```bash
VERSION=v1.0.0
curl -L https://github.com/njingu90/tfimport-cli/releases/download/${VERSION}/tfimport-cli-${VERSION}-darwin-amd64 \
  -o tfimport-cli
chmod +x tfimport-cli
sudo mv tfimport-cli /usr/local/bin/
tfimport-cli --version
```

#### Using Downloaded Binary (Apple Silicon/ARM64)

```bash
VERSION=v1.0.0
curl -L https://github.com/njingu90/tfimport-cli/releases/download/${VERSION}/tfimport-cli-${VERSION}-darwin-arm64 \
  -o tfimport-cli
chmod +x tfimport-cli
sudo mv tfimport-cli /usr/local/bin/
tfimport-cli --version
```

#### Building from Source

```bash
git clone https://github.com/njingu90/tfimport-cli.git
cd tfimport-cli
make build
sudo cp bin/tfimport-cli /usr/local/bin/
tfimport-cli --version
```

### Windows

#### Using Downloaded Binary (PowerShell)

1. **Download the release**:
   ```powershell
   $VERSION = "v1.0.0"
   $Url = "https://github.com/njingu90/tfimport-cli/releases/download/$VERSION/tfimport-cli-$VERSION-windows-amd64.exe"
   Invoke-WebRequest -Uri $Url -OutFile "tfimport-cli.exe"
   ```

2. **Move to PATH** (e.g., `C:\Program Files\tfimport-cli\`):
   ```powershell
   mkdir "C:\Program Files\tfimport-cli" -ErrorAction SilentlyContinue
   Move-Item tfimport-cli.exe "C:\Program Files\tfimport-cli\"
   ```

3. **Add to PATH** environment variable (if not already present):
   - Right-click "This PC" → Properties
   - Click "Advanced system settings"
   - Click "Environment Variables"
   - Under "System variables", click "Path" → "Edit"
   - Click "New" and add `C:\Program Files\tfimport-cli`
   - Click "OK"

4. **Verify installation** (new PowerShell window):
   ```powershell
   tfimport-cli --version
   ```

#### Building from Source (Windows)

```powershell
git clone https://github.com/njingu90/tfimport-cli.git
cd tfimport-cli
go build -o tfimport-cli.exe ./cmd/tfimport-cli
# Move tfimport-cli.exe to your PATH
```

## Terraform Cloud Authentication

### Create an API Token

1. **Log into Terraform Cloud**:
   - Go to https://app.terraform.io

2. **Create personal API token**:
   - Click your profile icon (top-right)
   - Select "User Settings"
   - Click "Tokens"
   - Click "Create an API token"
   - Give it a name (e.g., "tfimport-cli")
   - Click "Create token"
   - **Copy the token** (you won't see it again!)

3. **Or create a team token**:
   - In organization settings → Teams → API Tokens
   - Create and copy the token

### Set Environment Variable

#### Linux/macOS

```bash
export TF_API_TOKEN=<your-token>
```

To make it persistent, add to `~/.bashrc` or `~/.zshrc`:

```bash
echo 'export TF_API_TOKEN=<your-token>' >> ~/.bashrc
source ~/.bashrc
```

#### Windows (PowerShell)

```powershell
$env:TF_API_TOKEN = "<your-token>"
```

To make it persistent:
- Right-click "This PC" → Properties → "Advanced system settings"
- Click "Environment Variables"
- Click "New" under "User variables"
- Variable name: `TF_API_TOKEN`
- Variable value: `<your-token>`
- Click "OK"

#### Windows (Command Prompt)

```cmd
setx TF_API_TOKEN <your-token>
```

### Test Authentication

```bash
tfimport-cli analyze -organization my-org -workspace prod
```

If successful, you'll see a state analysis report.

## Upgrade Instructions

### From v1.0.0 to v1.1.0 (example)

1. **Download new version**:
   ```bash
   curl -L https://github.com/njingu90/tfimport-cli/releases/download/v1.1.0/tfimport-cli-v1.1.0-linux-amd64 \
     -o tfimport-cli
   ```

2. **Replace old binary**:
   ```bash
   chmod +x tfimport-cli
   sudo cp tfimport-cli /usr/local/bin/tfimport-cli
   ```

3. **Verify**:
   ```bash
   tfimport-cli --version
   ```

## Troubleshooting

### "Command not found: tfimport-cli"

1. **Check installation**:
   ```bash
   which tfimport-cli
   ```

2. **If not found, add to PATH**:
   ```bash
   export PATH=$PATH:/usr/local/bin
   ```

3. **Or reinstall**:
   ```bash
   sudo cp /path/to/tfimport-cli /usr/local/bin/
   ```

### "Permission denied"

Make sure the binary is executable:
```bash
chmod +x /usr/local/bin/tfimport-cli
```

### "Could not authenticate with Terraform Cloud"

1. **Verify token format**:
   - Should start with `skpenv-` (user token) or `atlasv1.` (team token)
   - Should be ~25-30 characters

2. **Check environment variable**:
   ```bash
   echo $TF_API_TOKEN
   ```

3. **Verify token isn't expired**:
   - Check in Terraform Cloud UI
   - Regenerate if needed

4. **Check organization and workspace names**:
   ```bash
   tfimport-cli analyze -organization my-org -workspace prod
   ```

### "State file not found"

```bash
# Make sure the path is correct
ls -la terraform.tfstate

# Try with absolute path
tfimport-cli analyze --state /full/path/to/terraform.tfstate
```

### "Unsupported state version"

tfimport-cli requires Terraform state v4. If you're using an older version:

```bash
# Upgrade Terraform
terraform version  # Check your version

# Refresh state to v4 format
terraform init
terraform refresh
```

## Support

- **Documentation**: https://github.com/njingu90/tfimport-cli
- **Issues**: https://github.com/njingu90/tfimport-cli/issues
- **Discussions**: https://github.com/njingu90/tfimport-cli/discussions

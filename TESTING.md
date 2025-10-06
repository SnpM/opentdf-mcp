# OpenTDF MCP Server - Testing Guide

## Quick Test

The MCP server is now working! Here's how to verify:

### 1. Build the server

```bash
cd /workspaces/opentdf-demo-server/opentdf-mcp
go build -o opentdf-mcp-server ./mcp-server
```

### 2. Test initialize handshake

Make sure `OPENTDF_PLATFORM_ENDPOINT` (and credentials) are set in your environment or in a copied `.env` file.

```bash
(echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}'; sleep 1) | timeout 3 ./opentdf-mcp-server 2>&1
```

Expected output:
```json
{"jsonrpc":"2.0","id":1,"result":{"capabilities":{"logging":{},"tools":{"listChanged":true}},"protocolVersion":"2025-06-18","serverInfo":{"name":"opentdf-mcp","version":"1.0.0"}}}
```

### 3. Test Tools List

```bash
(echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}'; echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}'; sleep 1) | timeout 3 ./opentdf-mcp-server 2>&1
```

You should see three tools registered:
- **encrypt**: Encrypt data using OpenTDF with attributes
- **decrypt**: Decrypt TDF or nanoTDF files automatically
- **list_attributes**: List available attributes from the platform

## Integration with Claude Desktop

To use the MCP server with Claude Desktop:

1. Copy the example configuration:
```bash
cp claude_desktop_config.example.json ~/Library/Application\ Support/Claude/claude_desktop_config.json
```

2. Update the configuration with your paths and credentials. Make sure `OPENTDF_PLATFORM_ENDPOINT` includes the protocol (for example `http://localhost:8080`):

```json
{
  "mcpServers": {
    "opentdf": {
      "command": "/path/to/opentdf-mcp-server",
      "env": {
        "OPENTDF_PLATFORM_ENDPOINT": "http://localhost:8080",
        "OPENTDF_CLIENT_ID": "opentdf-sdk",
        "OPENTDF_CLIENT_SECRET": "secret"
      }
    }
  }
}
```

3. Restart Claude Desktop

4. The OpenTDF tools should appear in the client's tools/capabilities list.

## Available Tools

### encrypt
Encrypt data using OpenTDF with specified attributes.

**Input:**
- `data` (required): The data to encrypt
- `attributes` (required): Array of attribute FQNs (e.g., `["https://example.com/attr/Classification/value/S"]`)
- `format` (optional): "tdf" or "nano" (default: "nano")
- `output` (optional): Output file path

**Output:**
- `success`: Boolean indicating success
- `outputFile`: Path to the encrypted file
- `message`: Human-readable result message

### decrypt
Decrypt a TDF or nanoTDF file (auto-detects format).

**Input:**
- `input` (required): Path to encrypted file or base64 encoded data
- `output` (optional): Output file path for decrypted data

**Output:**
- `success`: Boolean indicating success
- `decryptedData`: The plaintext data
- `error`: Error message if failed

### list_attributes
List available data attributes from the OpenTDF platform.

**Input:**
- `namespace` (optional): Filter by namespace (e.g., "https://example.com")
- `verbose` (optional): Show detailed attribute information

**Output:**
- `success`: Boolean indicating success
- `attributes`: Array of attribute objects with namespace, name, values, and FQN
- `error`: Error message if failed

## Troubleshooting

### Server won't start
- Check that the build completed successfully
- Verify Go version is 1.20 or higher
- Check for port conflicts if using network transport

### Connection errors
- Verify the OpenTDF platform is running (`docker-compose up` in the demo-server directory)
- Check the `OPENTDF_PLATFORM_ENDPOINT` environment variable
- Ensure credentials are correct (default: `opentdf-sdk:secret`)

### Permission denied errors
- Make sure you're using the correct credentials (opentdf-sdk:secret for demo)
- Verify the OpenTDF platform is configured properly
- Check KAS is running and accessible

## JSON Schema Tag Format

The MCP Go SDK uses a simple jsonschema tag format:

```go
type MyInput struct {
    Field string `json:"field" jsonschema:"Description of the field"`
}
```

**Important**: The tag value is just the description text itself, NOT `"description=text"` format.

## Next Steps

- See [MCP-SERVER.md](./MCP-SERVER.md) for detailed architecture
- See [README.md](./README.md) for general usage
- See [IMPLEMENTATION.md](./IMPLEMENTATION.md) for implementation notes

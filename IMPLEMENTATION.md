# OpenTDF MCP Implementation Summary

## What Was Implemented

An MCP (Model Context Protocol) server for OpenTDF that exposes encryption, decryption, and attribute management as AI-accessible tools.

## Files Created/Modified

### New Files
1. **mcp-server/main.go** - Main MCP server implementation
2. **mcp-server/config.go** - Configuration helpers
3. **MCP-SERVER.md** - Comprehensive MCP server documentation
4. **claude_desktop_config.example.json** - Example configuration for Claude Desktop
5. **test-mcp.sh** - Build and test script

### Modified Files
1. **go.mod** - Added MCP SDK dependency
2. **go.sum** - Updated dependencies
3. **README.md** - Added MCP server section
4. **cmd/main.go** - Updated default client ID to `opentdf-sdk`

## MCP Tools Provided

### 1. encrypt
- **Purpose:** Encrypt data with OpenTDF using specified attributes
- **Input:** data, attributes[], outputFile (optional), useNano (optional)
- **Output:** Success status, output file path, message/error
- **Use Case:** "Encrypt 'Hello World' with the secret classification attribute"

### 2. decrypt  
- **Purpose:** Decrypt TDF or nanoTDF files
- **Input:** inputFile path
- **Output:** Success status, decrypted data, error
- **Use Case:** "Decrypt encrypted.nano and show me the contents"

### 3. list_attributes
- **Purpose:** List available data attributes from the platform
- **Input:** namespace (optional), verbose (optional)
- **Output:** Success status, list of attributes with namespaces/values
- **Use Case:** "Show me all available attributes in the example.com namespace"

## Technical Details

### Stack
- **Language:** Go 1.24+
- **MCP SDK:** github.com/modelcontextprotocol/go-sdk v1.0.0
- **OpenTDF SDK:** github.com/opentdf/platform/sdk v0.8.0
- **Transport:** stdio (JSON-RPC 2.0 over stdin/stdout)
- **Protocol Version:** MCP 2025-06-18

### Architecture
```
┌─────────────────┐
│  MCP Client     │ (Claude Desktop, VS Code, etc.)
│  (AI Assistant) │
└────────┬────────┘
         │ stdio (JSON-RPC)
         │
┌────────▼────────┐
│ opentdf-mcp-    │
│    server       │
└────────┬────────┘
         │ gRPC/HTTP
         │
┌────────▼────────┐
│  OpenTDF        │
│  Platform       │
└─────────────────┘
```

### Default Configuration
- **Platform:** http://localhost:8080
- **Client ID:** opentdf-sdk (has proper KAS permissions)
- **Client Secret:** secret
- **Format Preference:** nanoTDF (better compatibility)

## Quick Start

### Build
```bash
cd opentdf-mcp
go build -o opentdf-mcp-server ./mcp-server
```

### Test
```bash
./test-mcp.sh
```

### Configure Claude Desktop
Edit `~/Library/Application Support/Claude/claude_desktop_config.json`:
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

### Use
Once configured, interact naturally:
- "Encrypt 'sensitive data' with the secret attribute"
- "List all attributes"
- "Decrypt encrypted.nano"

## Key Features

✅ **Auto-format Detection** - Automatically detects TDF vs nanoTDF  
✅ **Proper Authentication** - Uses demo credentials with KAS permissions  
✅ **Structured Responses** - JSON output with success/error states  
✅ **Natural Language** - Works with conversational AI queries  
✅ **Type Safety** - JSON schema validation for all inputs  

## Testing Notes

The implementation was tested for:
- ✅ Successful build
- ✅ CLI compatibility maintained
- ✅ Proper MCP tool registration
- ✅ Correct default credentials (opentdf-sdk:secret)

To test with a real MCP client:
1. Configure Claude Desktop as shown above
2. Restart Claude Desktop
3. Ask it to list OpenTDF attributes
4. Try encrypting and decrypting data

## Troubleshooting

**Issue:** Permission denied during decrypt  
**Fix:** Default credentials (opentdf-sdk:secret) have proper permissions. Ensure you're using them.

**Issue:** Server not found in Claude  
**Fix:** Use absolute path in config, ensure binary has execute permissions.

**Issue:** Can't connect to platform  
**Fix:** Verify OpenTDF platform is running at http://localhost:8080 or set OPENTDF_PLATFORM_ENDPOINT.

## Future Enhancements

Potential additions:
- Add `get_entitlements` tool for authorization queries
- Support for multiple output formats
- Batch encrypt/decrypt operations
- Resource templates for dynamic content
- Prompts for guided workflows

## Documentation

- **MCP-SERVER.md** - Full MCP server documentation
- **README.md** - Main project documentation with CLI and MCP sections
- **claude_desktop_config.example.json** - Example configuration

## Repository Structure

```
opentdf-mcp/
├── mcp-server/              # MCP server source
│   ├── main.go
│   └── config.go
├── cmd/                     # CLI source
│   ├── main.go
│   ├── encrypt.go
│   ├── decrypt.go
│   ├── attributes.go
│   └── entitlements.go
├── opentdf-mcp-server       # MCP server binary (built)
├── opentdf-cli              # CLI binary (built)
├── README.md                # Main documentation
├── MCP-SERVER.md            # MCP-specific docs
├── claude_desktop_config.example.json
├── test-mcp.sh              # Build/test script
├── go.mod
└── go.sum
```

## Conclusion

The OpenTDF MCP server successfully bridges OpenTDF's data-centric security with AI assistants through the Model Context Protocol, enabling natural language interactions for encryption, decryption, and policy management operations.

# OpenTDF MCP Server

An MCP (Model Context Protocol) server that exposes OpenTDF encryption, decryption, and policy management operations as tools for AI assistants and other MCP clients.

## What is MCP?

[Model Context Protocol (MCP)](https://modelcontextprotocol.io/) is an open protocol that enables AI applications to securely interact with local and remote resources. This server implements MCP to provide AI assistants with the ability to perform OpenTDF operations.

## Features

The OpenTDF MCP server provides three main tools:

### 1. `encrypt`
Encrypt data using OpenTDF with specified data attributes.

**Parameters:**
- `data` (required): The plaintext data to encrypt
- `attributes` (required): Array of attribute FQNs to apply
- `outputFile` (optional): Output file path (default: encrypted.nano)
- `useNano` (optional): Use nanoTDF format (default: true)

**Example:**
```json
{
  "data": "Sensitive information",
  "attributes": ["https://example.com/attr/class/value/secret"],
  "useNano": true
}
```

### 2. `decrypt`
Decrypt a TDF or nanoTDF file and return the plaintext.

**Parameters:**
- `inputFile` (required): Path to the encrypted file

**Example:**
```json
{
  "inputFile": "encrypted.nano"
}
```

### 3. `list_attributes`
List available data attributes from the OpenTDF platform.

**Parameters:**
- `namespace` (optional): Filter by namespace
- `verbose` (optional): Show detailed information including attribute values

**Example:**
```json
{
  "namespace": "https://example.com",
  "verbose": true
}
```

## Installation & Configuration

### Prerequisites

- Go 1.24 or higher
- OpenTDF platform running (default: http://localhost:8080)
- Valid client credentials with KAS permissions

### Build

```bash
cd opentdf-mcp
go build -o opentdf-mcp-server ./mcp-server
```

### Environment variables

Configure the server using these environment variables. Defaults below reflect the demo setup used in examples:

- `OPENTDF_PLATFORM_ENDPOINT` — Platform endpoint (default: `http://localhost:8080`)
- `OPENTDF_CLIENT_ID` — Client ID for authentication (default: `opentdf-sdk`)
- `OPENTDF_CLIENT_SECRET` — Client secret (default: `secret`)

These values are used throughout the docs and example scripts. If you run the platform on a different host or port, update `OPENTDF_PLATFORM_ENDPOINT` accordingly.

## MCP Client Configuration

### Claude Desktop

Add this to your Claude Desktop configuration file:

**macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`  
**Windows:** `%APPDATA%\Claude\claude_desktop_config.json`  
**Linux:** `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "opentdf": {
      "command": "/absolute/path/to/opentdf-mcp-server",
      "env": {
        "OPENTDF_PLATFORM_ENDPOINT": "http://localhost:8080",
        "OPENTDF_CLIENT_ID": "opentdf-sdk",
        "OPENTDF_CLIENT_SECRET": "secret"
      }
    }
  }
}
```

Replace `/absolute/path/to/opentdf-mcp-server` with the actual path to your built binary.

### VS Code with MCP Extension

If using VS Code with an MCP extension, configure it similarly:

```json
{
  "mcp.servers": {
    "opentdf": {
      "command": "/absolute/path/to/opentdf-mcp-server",
      "env": {
        "OPENTDF_PLATFORM_ENDPOINT": "http://localhost:8080",
        "OPENTDF_CLIENT_ID": "opentdf-sdk",
        "OPENTDF_CLIENT_SECRET": "secret"
      }
    }
  }
}
```

### Custom MCP Clients

The server uses stdio transport and follows the MCP specification. To integrate:

1. Start the server process: `./opentdf-mcp-server`
2. Communicate via stdin/stdout using JSON-RPC 2.0
3. The server implements the MCP protocol version 2025-06-18

## Usage Examples (via MCP Client)

Once configured in your MCP client, you can use natural language:

### Encrypt Data
> "Encrypt the text 'Hello World' with the attribute https://example.com/attr/class/value/secret using nanoTDF format"

The AI assistant will call the `encrypt` tool with the appropriate parameters.

### Decrypt a File
> "Decrypt the file encrypted.nano and show me the contents"

### List Attributes
> "Show me all available OpenTDF attributes in verbose mode"

### Combined Operations
> "List the attributes, then encrypt 'Confidential data' with the secret classification attribute"

## Testing

Run the test script to verify the build:

```bash
./test-mcp.sh
```

For end-to-end testing with an actual MCP client:

1. Configure the server in Claude Desktop or another MCP client
2. Restart the client
3. Verify the server appears in the client's tools/capabilities list
4. Try a simple operation like listing attributes

## Troubleshooting

### Server not appearing in MCP client

- Verify the absolute path to `opentdf-mcp-server` is correct
- Check that the binary has execute permissions
- Review the MCP client's logs for connection errors

### Permission denied errors during decrypt

The default client credentials (`opentdf-sdk:secret`) have the necessary KAS permissions. If you change the credentials, ensure the new client has appropriate entitlements in the OpenTDF platform.

### Platform connection errors

- Verify `OPENTDF_PLATFORM_ENDPOINT` is correct and reachable
- Ensure the OpenTDF platform is running
- Check network connectivity if using a remote endpoint

## Development

The MCP server is implemented in Go using the official [go-sdk](https://github.com/modelcontextprotocol/go-sdk).

### Project Structure

```
opentdf-mcp/
├── mcp-server/
│   ├── main.go       # MCP server implementation
│   └── config.go     # Configuration helpers
├── cmd/
│   └── ...           # CLI implementation
└── README.md         # Main documentation
```

### Adding New Tools

To add a new MCP tool:

1. Define input/output structs with JSON schema tags
2. Implement the tool handler function
3. Register the tool with `mcp.AddTool()` in `main.go`

Example:

```go
type MyToolInput struct {
    Param string `json:"param" jsonschema:"required,description=Parameter description"`
}

func MyTool(ctx context.Context, req *mcp.CallToolRequest, input MyToolInput) (*mcp.CallToolResult, MyToolOutput, error) {
    // Implementation
}

// In runMCPServer():
mcp.AddTool(server, &mcp.Tool{
    Name:        "my_tool",
    Description: "Tool description",
}, MyTool)
```

## Security Considerations

- **Credentials:** The server uses client credentials to authenticate with the OpenTDF platform. Keep `OPENTDF_CLIENT_SECRET` secure.
- **File Access:** The server can read/write files in the working directory. Run it in a restricted directory if needed.
- **Network:** The server connects to the configured OpenTDF platform endpoint. Ensure secure connections for production use.

## License

See the main repository LICENSE file.

## Links

- [OpenTDF Platform](https://github.com/opentdf/platform)
- [Model Context Protocol](https://modelcontextprotocol.io/)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)

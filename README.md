# OpenTDF Hello World

A simple Go application demonstrating the use of the OpenTDF Platform SDK v0.11.0.

## Prerequisites

- Go 1.24 or higher
- Docker/Podman (to run OpenTDF platform locally)

## Setup

This project uses:
- `github.com/opentdf/platform/protocol/go v0.11.0` - Protocol definitions
- `github.com/opentdf/platform/sdk v0.8.0` - SDK for interacting with OpenTDF platform
- `github.com/modelcontextprotocol/go-sdk v1.0.0` - Model Context Protocol SDK

Dependencies are already configured in `go.mod`.

---

# OpenTDF CLI (opentdf-mcp)

A small Go-based CLI that demonstrates using the OpenTDF SDK to
create and consume TDF and nanoTDF data, request entitlements, and list
policy attributes against an OpenTDF platform.

This README focuses on using the `opentdf-cli` tool shipped in this
directory (build instructions, examples, troubleshooting notes).

## Prerequisites

- Go 1.24 or higher
- Docker/Podman (optional)

Dependencies are configured in `go.mod` for the example code.

## Build

From this directory build the CLI binary:

```bash
cd opentdf-mcp
go build -o opentdf-cli ./cmd
```

That creates `opentdf-cli` (or use `go run ./cmd ...` to run without building).

To build the MCP server:

```bash
go build -o opentdf-mcp-server ./mcp-server
```

## Environment variables

The project expects the following environment variables. Defaults shown below are the demo values used throughout these docs and example scripts.

- OPENTDF_PLATFORM_ENDPOINT — Platform endpoint (default: `http://localhost:8080`)
- OPENTDF_CLIENT_ID — Client ID for authentication (default: `opentdf-sdk`)
- OPENTDF_CLIENT_SECRET — Client secret for authentication (default: `secret`)

These defaults are the demo credentials used in the repository fixtures and have the KAS permissions needed for the examples. You can set them in a shell or copy the provided `.env.template` into `.env`.

## Platform

This CLI works against an OpenTDF platform endpoint. Set the endpoint
with `OPENTDF_PLATFORM_ENDPOINT` (default: `http://localhost:8080`).

Ensure the platform's KAS (Key Access Server) is reachable if you plan
to create or decrypt TDF files that require KAS operations.

---

# Model Context Protocol (MCP) Server

The OpenTDF MCP server exposes OpenTDF operations as MCP tools that can be used by MCP clients like Claude Desktop, IDEs with MCP support, or custom applications.

## MCP Server Features

The server provides the following tools:

1. **encrypt** - Encrypt data using OpenTDF with specified attributes
   - Supports both TDF and nanoTDF formats
   - Creates encrypted files with policy bindings

2. **decrypt** - Decrypt TDF or nanoTDF files
   - Automatically detects format
   - Returns plaintext data

3. **list_attributes** - List available data attributes from the platform
   - Optional namespace filtering
   - Verbose mode shows attribute values

## Running the MCP Server

Start the MCP server:

```bash
./opentdf-mcp-server
```

The server communicates over stdio using the Model Context Protocol.

## Configuring MCP Clients

### Claude Desktop

Add to your Claude Desktop config file (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

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

### Other MCP Clients

Use the command `opentdf-mcp-server` with stdio transport. The server follows the [Model Context Protocol specification](https://modelcontextprotocol.io/).

## Example MCP Tool Usage

Once configured in an MCP client, you can use natural language to interact with OpenTDF:

- "Encrypt 'Hello World' with the attribute https://example.com/attr/class/value/secret"
- "List all available attributes"
- "Decrypt the file encrypted.nano"

The MCP tools provide structured responses with success/failure status and detailed error messages.

---

# CLI Commands & Examples

Build the binary (see Build section) and then use the following examples.

Encrypt (recommended: nanoTDF)

The easiest way to get a working encrypt/decrypt flow against a platform
without needing external KAS hostnames is to use nanoTDF mode. Nano mode
uses the platform's KAS at `<OPENTDF_PLATFORM_ENDPOINT>/kas`.

```bash
./opentdf-cli encrypt -a https://example.com/attr/attr1/value/value1 -o encrypted.nano "Hello Nano"
```

Encrypt (standard TDF — advanced)

Standard TDF encryption may require contacting the Key Access Server
(KAS) configured for the attribute values used. If an attribute's key
access server points to an external host (for example `kas.example.com`)
and that host is not resolvable or reachable from your environment, the
CLI will fail with an error like:

```
Error: failed to create TDF: unable to retrieve public key from KAS at [https://kas.example.com]: error making request to KAS: unavailable: dial tcp: lookup kas.example.com: no such host
```

If you need to use standard TDF (non-nano):

1. Ensure `OPENTDF_PLATFORM_ENDPOINT` points to a platform whose KAS
    endpoint is reachable and correctly configured for the attribute
    values you plan to use.
2. Use attributes whose key access server is reachable from your host,
    or update platform fixtures / configuration to point to a reachable KAS.
3. (Not recommended for general use) Add a host entry so the KAS
    hostname resolves to an appropriate address.

Example (may fail if KAS is unreachable for the chosen attributes):

```bash
./opentdf-cli encrypt -a https://example.com/attr/attr1/value/value1 -o encrypted.tdf "Hello World"
```

Notes:
- Use `-a` multiple times to supply multiple data attributes.
- `-o` sets the output file (default `encrypted.tdf`).

Decrypt (prints plaintext to stdout)

```bash
# default (uses OPENTDF_CLIENT_ID/OPENTDF_CLIENT_SECRET env vars or defaults)
./opentdf-cli decrypt encrypted.tdf > decrypted.txt
```

If decryption fails with a KAS permission error (see Troubleshooting),
use the demo client credentials which are configured in the fixtures and
have the required KAS permissions for the example attributes:

```bash
OPENTDF_CLIENT_ID=opentdf-sdk OPENTDF_CLIENT_SECRET=secret ./opentdf-cli decrypt encrypted.nano > decrypted.txt
cat decrypted.txt
# expected output: Hello Nano
```

Get entitlements (Authorization V2)

```bash
./opentdf-cli get-entitlements --identifier user@example.com --type email
```

Attributes list

```bash
# list all namespaces/attributes
./opentdf-cli attributes list

# verbose and filter by namespace
./opentdf-cli attributes list -l -N https://example.com
```

Help

```bash
./opentdf-cli help
```

## Troubleshooting

- Unknown attribute FQN (ErrNotFound): The platform will return an error
    if you attempt to use an attribute FQN that doesn't exist. Use
    `./opentdf-cli attributes list` to find valid FQNs.
- KAS unreachable: If encrypt (standard TDF) fails because the CLI
    attempted to contact a KAS (for example `kas.example.com`), ensure
    your `OPENTDF_PLATFORM_ENDPOINT` points to a platform whose KAS
    endpoint is reachable. If your platform exposes KAS at
    `<platform>/kas`, you can use `-nano` so the CLI targets that
    endpoint when creating nanoTDF.
- permission_denied during decrypt (nanoTDF): The KAS must rewrap the
    key for the requesting client. The demo fixtures give that permission
    to the client id `opentdf-sdk`. If you see `permission_denied`, try
    decrypting with `OPENTDF_CLIENT_ID=opentdf-sdk OPENTDF_CLIENT_SECRET=secret`.

## Where the code is

- CLI code: `opentdf-mcp/cmd/*.go` (main, encrypt, decrypt, entitlements, attributes)
- Example README: this file

## Learn more

- OpenTDF Platform: https://github.com/opentdf/platform
- OpenTDF docs: https://opentdf.io/

If you'd like, I can add a short example script in this directory that
runs a build → encrypt (nano) → decrypt flow using the demo client
credentials so you can reproduce the flow with one command.

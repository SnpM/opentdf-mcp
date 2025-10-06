#!/bin/bash
# Test script for OpenTDF MCP server

set -e

echo "Testing OpenTDF MCP Server..."
echo

# Build the server
echo "Building MCP server..."
go build -o opentdf-mcp-server ./mcp-server
echo "✓ Build successful"
echo

# Test help
echo "Testing CLI..."
./opentdf-cli help
echo
echo "✓ CLI help works"
echo

# Note: To actually test the MCP server, you would need an MCP client
# For now, we just verify it builds and can be invoked
echo "MCP server binary ready at: ./opentdf-mcp-server"
echo "To test with an MCP client, configure it to run:"
echo "  ./opentdf-mcp-server"
echo

echo "All tests passed!"

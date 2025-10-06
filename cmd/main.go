package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	var err error
	switch command {
	case "encrypt":
		err = handleEncrypt()
	case "decrypt":
		err = handleDecrypt()
	case "get-entitlements":
		err = handleGetEntitlements()
	case "attributes":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Error: attributes subcommand required\n")
			os.Exit(1)
		}
		subcommand := os.Args[2]
		if subcommand == "list" {
			err = handleAttributesList()
		} else {
			fmt.Fprintf(os.Stderr, "Error: unknown attributes subcommand: %s\n", subcommand)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		printUsage()
		return
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("OpenTDF CLI - Command line interface for OpenTDF operations")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  opentdf-cli <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  encrypt             Encrypt data using TDF")
	fmt.Println("  decrypt             Decrypt a TDF file")
	fmt.Println("  get-entitlements    Get entitlements for an entity")
	fmt.Println("  attributes list     List available attributes")
	fmt.Println("  help                Show this help message")
	fmt.Println()
	fmt.Println("Environment Variables:")
	fmt.Println("  OPENTDF_PLATFORM_ENDPOINT   Platform endpoint (default: http://localhost:8080)")
	fmt.Println("  OPENTDF_CLIENT_ID           Client ID for authentication (default: opentdf-sdk)")
	fmt.Println("  OPENTDF_CLIENT_SECRET       Client secret for authentication (default: secret)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  opentdf-cli encrypt -a https://example.com/attr/class/value/secret \"Hello World\"")
	fmt.Println("  OPENTDF_CLIENT_ID=opentdf-sdk OPENTDF_CLIENT_SECRET=secret ./opentdf-cli decrypt encrypted.tdf")
	fmt.Println("  opentdf-cli get-entitlements --identifier user@example.com --type email")
	fmt.Println("  opentdf-cli attributes list -l")
	fmt.Println()
	fmt.Println("For MCP Server:")
	fmt.Println("  Use the separate 'opentdf-mcp-server' binary for Model Context Protocol support")
}

func getPlatformEndpoint() string {
	if endpoint := os.Getenv("OPENTDF_PLATFORM_ENDPOINT"); endpoint != "" {
		return endpoint
	}
	return "http://localhost:8080"
}

func getClientID() string {
	if clientID := os.Getenv("OPENTDF_CLIENT_ID"); clientID != "" {
		return clientID
	}
	return "opentdf-sdk"
}

func getClientSecret() string {
	if secret := os.Getenv("OPENTDF_CLIENT_SECRET"); secret != "" {
		return secret
	}
	return "secret"
}

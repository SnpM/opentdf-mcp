// OpenTDF MCP Server
// Run as a standalone Model Context Protocol server for OpenTDF operations
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/opentdf/platform/protocol/go/policy/attributes"
	"github.com/opentdf/platform/protocol/go/policy/namespaces"
	"github.com/opentdf/platform/sdk"
)

// Tool input/output schemas
// EncryptToolInput defines the input for the encrypt tool
type EncryptToolInput struct {
	Data       string   `json:"data" jsonschema:"The data to encrypt"`
	Attributes []string `json:"attributes" jsonschema:"Data attributes (FQNs) to apply during encryption"`
	Format     string   `json:"format,omitempty" jsonschema:"Output format: 'tdf' or 'nano' (default: nano)"`
	Output     string   `json:"output,omitempty" jsonschema:"Output file path (optional, returns base64 if not specified)"`
}

type EncryptToolOutput struct {
	Success    bool   `json:"success"`
	OutputFile string `json:"outputFile"`
	Message    string `json:"message,omitempty"`
	Error      string `json:"error,omitempty"`
}

// DecryptToolInput defines the input for the decrypt tool
type DecryptToolInput struct {
	Input  string `json:"input" jsonschema:"Path to encrypted file or base64 encoded data"`
	Output string `json:"output,omitempty" jsonschema:"Output file path (optional, returns plaintext if not specified)"`
}

type DecryptToolOutput struct {
	Success      bool   `json:"success"`
	DecryptedData string `json:"decryptedData,omitempty"`
	Error        string `json:"error,omitempty"`
}

type ListAttributesToolInput struct {
	Namespace string `json:"namespace,omitempty" jsonschema:"Filter by namespace (e.g., https://example.com)"`
	Verbose   bool   `json:"verbose,omitempty" jsonschema:"Show detailed attribute information"`
}

type ListAttributesToolOutput struct {
	Success    bool                   `json:"success"`
	Attributes []AttributeInfo        `json:"attributes,omitempty"`
	Error      string                 `json:"error,omitempty"`
}

type AttributeInfo struct {
	Namespace string   `json:"namespace"`
	Name      string   `json:"name"`
	Values    []string `json:"values,omitempty"`
	FQN       string   `json:"fqn"`
}

// getSDKClientMCP creates an authenticated OpenTDF SDK client for MCP
func getSDKClientMCP() (*sdk.SDK, error) {
	platformEndpoint := getPlatformEndpoint()
	clientID := getClientID()
	clientSecret := getClientSecret()

	var opts []sdk.Option
	if clientID != "" && clientSecret != "" {
		opts = append(opts, sdk.WithClientCredentials(clientID, clientSecret, nil))
	} else {
		opts = append(opts, sdk.WithInsecurePlaintextConn())
	}

	client, err := sdk.New(platformEndpoint, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create SDK client: %w", err)
	}

	return client, nil
}

// MCPEncrypt encrypts data with the given attributes
func MCPEncrypt(ctx context.Context, req *mcp.CallToolRequest, input EncryptToolInput) (*mcp.CallToolResult, EncryptToolOutput, error) {
	client, err := getSDKClientMCP()
	if err != nil {
		return nil, EncryptToolOutput{Success: false, Error: err.Error()}, nil
	}
	defer client.Close()

	// Determine format
	useNano := input.Format == "" || input.Format == "nano"

	outputFile := input.Output
	if outputFile == "" {
		if useNano {
			outputFile = "encrypted.nano"
		} else {
			outputFile = "encrypted.tdf"
		}
	}

	// Encrypt
	reader := strings.NewReader(input.Data)
	file, err := os.Create(outputFile)
	if err != nil {
		return nil, EncryptToolOutput{Success: false, Error: fmt.Sprintf("failed to create output file: %v", err)}, nil
	}
	defer file.Close()

	baseKasURL := getPlatformEndpoint()
	if !strings.HasPrefix(baseKasURL, "http://") && !strings.HasPrefix(baseKasURL, "https://") {
		baseKasURL = "http://" + baseKasURL
	}

	if useNano {
		nanoConfig, err := client.NewNanoTDFConfig()
		if err != nil {
			return nil, EncryptToolOutput{Success: false, Error: fmt.Sprintf("failed to create nanoTDF config: %v", err)}, nil
		}

		if len(input.Attributes) > 0 {
			if err := nanoConfig.SetAttributes(input.Attributes); err != nil {
				return nil, EncryptToolOutput{Success: false, Error: fmt.Sprintf("failed to set attributes: %v", err)}, nil
			}
		}

		nanoConfig.EnableECDSAPolicyBinding()

		if err := nanoConfig.SetKasURL(baseKasURL + "/kas"); err != nil {
			return nil, EncryptToolOutput{Success: false, Error: fmt.Sprintf("failed to set KAS URL: %v", err)}, nil
		}

		if _, err := client.CreateNanoTDF(file, reader, *nanoConfig); err != nil {
			return nil, EncryptToolOutput{Success: false, Error: fmt.Sprintf("failed to encrypt (nano): %v", err)}, nil
		}
	} else {
		if _, err := client.CreateTDF(file, reader, sdk.WithDataAttributes(input.Attributes...)); err != nil {
			return nil, EncryptToolOutput{Success: false, Error: fmt.Sprintf("failed to encrypt: %v", err)}, nil
		}
	}

	msg := fmt.Sprintf("Successfully encrypted data to %s", outputFile)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: msg},
		},
	}, EncryptToolOutput{Success: true, OutputFile: outputFile, Message: msg}, nil
}

// MCPDecrypt decrypts a TDF or nanoTDF file
func MCPDecrypt(ctx context.Context, req *mcp.CallToolRequest, input DecryptToolInput) (*mcp.CallToolResult, DecryptToolOutput, error) {
	client, err := getSDKClientMCP()
	if err != nil {
		return nil, DecryptToolOutput{Success: false, Error: err.Error()}, nil
	}
	defer client.Close()

	file, err := os.Open(input.Input)
	if err != nil {
		return nil, DecryptToolOutput{Success: false, Error: fmt.Sprintf("failed to open input file: %v", err)}, nil
	}
	defer file.Close()

	// Detect format by reading magic bytes
	var magic [3]byte
	var isNano bool
	n, err := io.ReadFull(file, magic[:])
	switch {
	case err != nil:
		return nil, DecryptToolOutput{Success: false, Error: fmt.Sprintf("failed to read magic bytes: %v", err)}, nil
	case n < 3:
		return nil, DecryptToolOutput{Success: false, Error: "file too small; no magic number found"}, nil
	case bytes.HasPrefix(magic[:], []byte("L1L")):
		isNano = true
	default:
		isNano = false
	}

	// Reset file position
	if _, err := file.Seek(0, 0); err != nil {
		return nil, DecryptToolOutput{Success: false, Error: fmt.Sprintf("failed to seek to beginning: %v", err)}, nil
	}

	var output bytes.Buffer

	if isNano {
		if _, err := client.ReadNanoTDF(&output, file); err != nil {
			return nil, DecryptToolOutput{Success: false, Error: fmt.Sprintf("failed to decrypt nanoTDF: %v", err)}, nil
		}
	} else {
		tdfReader, err := client.LoadTDF(file)
		if err != nil {
			return nil, DecryptToolOutput{Success: false, Error: fmt.Sprintf("failed to load TDF: %v", err)}, nil
		}
		if _, err := io.Copy(&output, tdfReader); err != nil && err != io.EOF {
			return nil, DecryptToolOutput{Success: false, Error: fmt.Sprintf("failed to decrypt TDF: %v", err)}, nil
		}
	}

	decryptedData := output.String()
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Successfully decrypted:\n%s", decryptedData)},
		},
	}, DecryptToolOutput{Success: true, DecryptedData: decryptedData}, nil
}

// MCPListAttributes lists available attributes
func MCPListAttributes(ctx context.Context, req *mcp.CallToolRequest, input ListAttributesToolInput) (*mcp.CallToolResult, ListAttributesToolOutput, error) {
	client, err := getSDKClientMCP()
	if err != nil {
		return nil, ListAttributesToolOutput{Success: false, Error: err.Error()}, nil
	}
	defer client.Close()

	// List namespaces
	var nsuris []string
	if input.Namespace == "" {
		listResp, err := client.Namespaces.ListNamespaces(ctx, &namespaces.ListNamespacesRequest{})
		if err != nil {
			return nil, ListAttributesToolOutput{Success: false, Error: fmt.Sprintf("failed to list namespaces: %v", err)}, nil
		}
		for _, n := range listResp.GetNamespaces() {
			nsuris = append(nsuris, n.GetFqn())
		}
	} else {
		nsuris = []string{input.Namespace}
	}

	var attributesList []AttributeInfo
	var textOutput strings.Builder

	for _, ns := range nsuris {
		// Parse namespace to get host
		parts := strings.Split(ns, "//")
		var host string
		if len(parts) > 1 {
			host = strings.Split(parts[1], "/")[0]
		} else {
			host = ns
		}

		lsr, err := client.Attributes.ListAttributes(ctx, &attributes.ListAttributesRequest{
			Namespace: host,
		})
		if err != nil {
			continue
		}

		textOutput.WriteString(fmt.Sprintf("Namespace: %s\n", ns))
		for _, a := range lsr.GetAttributes() {
			values := []string{}
			if input.Verbose {
				for _, v := range a.GetValues() {
					values = append(values, v.GetValue())
				}
			}

			attrInfo := AttributeInfo{
				Namespace: ns,
				Name:      a.GetName(),
				Values:    values,
				FQN:       a.GetFqn(),
			}
			attributesList = append(attributesList, attrInfo)

			if input.Verbose {
				textOutput.WriteString(fmt.Sprintf("  Attribute: %s\n", a.GetFqn()))
				if len(values) > 0 {
					textOutput.WriteString(fmt.Sprintf("  Values: %s\n", strings.Join(values, ", ")))
				}
			} else {
				textOutput.WriteString(fmt.Sprintf("  %s\n", a.GetFqn()))
			}
		}
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: textOutput.String()},
		},
	}, ListAttributesToolOutput{Success: true, Attributes: attributesList}, nil
}

func runMCPServer() error {
	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "opentdf-mcp",
		Version: "1.0.0",
	}, &mcp.ServerOptions{
		InitializedHandler: func(ctx context.Context, req *mcp.InitializedRequest) {
			log.Println("MCP server initialized")
		},
	})

	// Add encrypt tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "encrypt",
		Description: "Encrypt data using OpenTDF with the specified attributes. Creates a TDF or nanoTDF file. Use nanoTDF format for better compatibility.",
	}, MCPEncrypt)

	// Add decrypt tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "decrypt",
		Description: "Decrypt a TDF or nanoTDF file and return the plaintext data. Automatically detects the format.",
	}, MCPDecrypt)

	// Add list attributes tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_attributes",
		Description: "List available data attributes from the OpenTDF platform. Use verbose mode to see attribute values. Filter by namespace if needed.",
	}, MCPListAttributes)

	// Run server over stdio
	log.Println("Starting OpenTDF MCP server on stdio...")
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}

func main() {
	if err := runMCPServer(); err != nil {
		log.Fatalf("MCP server error: %v", err)
	}
}

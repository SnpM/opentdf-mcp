package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/opentdf/platform/sdk"
)

// handleEncrypt processes the encrypt command to create a nanoTDF encrypted file.
// It accepts plaintext data and optional attributes to encrypt the data using OpenTDF.
//
// Usage:
//   encrypt [flags] <plaintext>
//
// Flags:
//   -o string
//       Output file path (default "encrypted.tdf")
//   -a string
//       Data attribute FQN (can be specified multiple times)
//       Example: -a https://example.com/attr/attr1/value/value1
//
// The function:
//  1. Parses command-line flags and plaintext input
//  2. Retrieves platform endpoint and authentication credentials from environment
//  3. Creates an authenticated OpenTDF SDK client
//  4. Configures nanoTDF with attributes and ECDSA policy binding
//  5. Encrypts the plaintext data and writes it to the output file
//
// Returns an error if any step fails, including flag parsing, client creation,
// attribute configuration, or encryption operations.
func handleEncrypt() error {
	fs := flag.NewFlagSet("encrypt", flag.ExitOnError)
	output := fs.String("o", "encrypted.tdf", "Output file path")

	// Parse attributes flag multiple times
	var attributes []string
	fs.Func("a", "Data attribute (can be specified multiple times)", func(s string) error {
		attributes = append(attributes, s)
		return nil
	})

	if err := fs.Parse(os.Args[2:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("plaintext data is required")
	}

	plaintext := fs.Arg(0)

	// Retrieve configuration from environment variables
	platformEndpoint := getPlatformEndpoint()
	clientID := getClientID()
	clientSecret := getClientSecret()

	// Create authenticated client with appropriate security settings
	// If credentials are provided, use OAuth2 client credentials flow
	// Otherwise, use insecure plaintext connection (suitable for local development)
	var opts []sdk.Option
	if clientID != "" && clientSecret != "" {
		opts = append(opts, sdk.WithClientCredentials(clientID, clientSecret, nil))
	} else {
		opts = append(opts, sdk.WithInsecurePlaintextConn())
	}

	client, err := sdk.New(platformEndpoint, opts...)
	if err != nil {
		return fmt.Errorf("failed to create SDK client: %w", err)
	}
	defer client.Close()

	// Open output file
	outFile, err := os.Create(*output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	in := strings.NewReader(plaintext)

	// Ensure the platform endpoint has a valid HTTP/HTTPS scheme
	baseKasURL := platformEndpoint
	if !strings.HasPrefix(baseKasURL, "http://") && !strings.HasPrefix(baseKasURL, "https://") {
		baseKasURL = "http://" + baseKasURL
	}

	// Create nanoTDF configuration (nanoTDF is a compact binary format)
	nanoConfig, err := client.NewNanoTDFConfig()
	if err != nil {
		return fmt.Errorf("failed to create nanoTDF config: %w", err)
	}

	// Apply data attributes if specified
	// Attributes define the access policy for the encrypted data
	if len(attributes) > 0 {
		if err := nanoConfig.SetAttributes(attributes); err != nil {
			return fmt.Errorf("failed to set attributes: %w", err)
		}
	}

	// Enable ECDSA policy binding for enhanced security
	// This cryptographically binds the policy to the encrypted data
	nanoConfig.EnableECDSAPolicyBinding()

	// Configure the Key Access Service (KAS) endpoint
	// KAS is responsible for key management and access control
	if err := nanoConfig.SetKasURL(baseKasURL + "/kas"); err != nil {
		return fmt.Errorf("failed to set KAS URL: %w", err)
	}

	// Perform the encryption operation
	if _, err := client.CreateNanoTDF(outFile, in, *nanoConfig); err != nil {
		return fmt.Errorf("failed to create nanoTDF: %w", err)
	}

	fmt.Printf("Successfully encrypted to nanoTDF: %s\n", *output)

	return nil
}

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/opentdf/platform/sdk"
)

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

	platformEndpoint := getPlatformEndpoint()
	clientID := getClientID()
	clientSecret := getClientSecret()

	// Create authenticated client
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

	baseKasURL := platformEndpoint
	if !strings.HasPrefix(baseKasURL, "http://") && !strings.HasPrefix(baseKasURL, "https://") {
		baseKasURL = "http://" + baseKasURL
	}

	// Create nanoTDF (nano mode is the default behavior)
	nanoConfig, err := client.NewNanoTDFConfig()
	if err != nil {
		return fmt.Errorf("failed to create nanoTDF config: %w", err)
	}

	if len(attributes) > 0 {
		if err := nanoConfig.SetAttributes(attributes); err != nil {
			return fmt.Errorf("failed to set attributes: %w", err)
		}
	}

	nanoConfig.EnableECDSAPolicyBinding()

	if err := nanoConfig.SetKasURL(baseKasURL + "/kas"); err != nil {
		return fmt.Errorf("failed to set KAS URL: %w", err)
	}

	if _, err := client.CreateNanoTDF(outFile, in, *nanoConfig); err != nil {
		return fmt.Errorf("failed to create nanoTDF: %w", err)
	}

	fmt.Printf("Successfully encrypted to nanoTDF: %s\n", *output)

	return nil
}

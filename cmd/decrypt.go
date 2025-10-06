package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/opentdf/platform/sdk"
)

func handleDecrypt() error {
	fs := flag.NewFlagSet("decrypt", flag.ExitOnError)

	if err := fs.Parse(os.Args[2:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("input file is required")
	}

	inputFile := fs.Arg(0)

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

	// Open input file
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	// Detect format by reading magic bytes
	var magic [3]byte
	var isNano bool
	n, err := io.ReadFull(file, magic[:])
	switch {
	case err != nil:
		return fmt.Errorf("failed to read magic bytes: %w", err)
	case n < 3:
		return fmt.Errorf("file too small; no magic number found")
	case bytes.HasPrefix(magic[:], []byte("L1L")):
		isNano = true
	default:
		isNano = false
	}

	// Reset file position
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek to beginning: %w", err)
	}

	if isNano {
		// Decrypt nanoTDF
		if _, err := client.ReadNanoTDF(os.Stdout, file); err != nil {
			return fmt.Errorf("failed to decrypt nanoTDF: %w", err)
		}
	} else {
		// Decrypt standard TDF
		tdfReader, err := client.LoadTDF(file)
		if err != nil {
			return fmt.Errorf("failed to load TDF: %w", err)
		}

		if _, err := io.Copy(os.Stdout, tdfReader); err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("failed to decrypt TDF: %w", err)
		}
	}

	return nil
}

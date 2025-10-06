package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/opentdf/platform/protocol/go/policy/attributes"
	"github.com/opentdf/platform/protocol/go/policy/namespaces"
	"github.com/opentdf/platform/sdk"
)

func handleAttributesList() error {
	fs := flag.NewFlagSet("attributes list", flag.ExitOnError)
	verbose := fs.Bool("l", false, "Include detailed information")
	namespace := fs.String("N", "", "Filter by namespace")

	if err := fs.Parse(os.Args[3:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

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

	ctx := context.Background()

	var nsuris []string
	if *namespace == "" {
		// List all namespaces
		listResp, err := client.Namespaces.ListNamespaces(ctx, &namespaces.ListNamespacesRequest{})
		if err != nil {
			return fmt.Errorf("failed to list namespaces: %w", err)
		}

		for _, n := range listResp.GetNamespaces() {
			nsuris = append(nsuris, n.GetFqn())
		}
	} else {
		nsuris = strings.Split(*namespace, " ")
	}

	for _, ns := range nsuris {
		u, err := url.Parse(ns)
		if err != nil {
			return fmt.Errorf("failed to parse namespace URL: %w", err)
		}

		lsr, err := client.Attributes.ListAttributes(ctx, &attributes.ListAttributesRequest{
			Namespace: u.Host,
		})
		if err != nil {
			return fmt.Errorf("failed to list attributes: %w", err)
		}

		fmt.Printf("Namespace: %s\n", ns)
		for _, a := range lsr.GetAttributes() {
			if *verbose {
				fmt.Printf("  %s\t%s\n", a.GetFqn(), a.GetId())
			} else {
				fmt.Printf("  %s\n", a.GetFqn())
			}

			for _, v := range a.GetValues() {
				if *verbose {
					fmt.Printf("    %s\t%s\n", v.GetFqn(), v.GetId())
				} else {
					fmt.Printf("    %s\n", v.GetFqn())
				}
			}
		}
	}

	return nil
}

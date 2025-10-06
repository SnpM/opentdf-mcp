package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	authorizationv2 "github.com/opentdf/platform/protocol/go/authorization/v2"
	"github.com/opentdf/platform/protocol/go/entity"
	"github.com/opentdf/platform/sdk"
)

func handleGetEntitlements() error {
	fs := flag.NewFlagSet("get-entitlements", flag.ExitOnError)
	identifier := fs.String("identifier", "", "Entity identifier (e.g., email address)")
	identifierType := fs.String("type", "email", "Identifier type (email or username)")

	if err := fs.Parse(os.Args[2:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	if *identifier == "" {
		return fmt.Errorf("identifier is required")
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

	// Build entity based on identifier type
	var entityIdentifier *authorizationv2.EntityIdentifier
	switch *identifierType {
	case "email":
		entityIdentifier = &authorizationv2.EntityIdentifier{
			Identifier: &authorizationv2.EntityIdentifier_EntityChain{
				EntityChain: &entity.EntityChain{
					Entities: []*entity.Entity{
						{
							EphemeralId: "user-" + *identifier,
							EntityType: &entity.Entity_EmailAddress{
								EmailAddress: *identifier,
							},
						},
					},
				},
			},
		}
	case "username":
		entityIdentifier = &authorizationv2.EntityIdentifier{
			Identifier: &authorizationv2.EntityIdentifier_EntityChain{
				EntityChain: &entity.EntityChain{
					Entities: []*entity.Entity{
						{
							EphemeralId: "user-" + *identifier,
							EntityType: &entity.Entity_UserName{
								UserName: *identifier,
							},
						},
					},
				},
			},
		}
	default:
		return fmt.Errorf("unsupported identifier type: %s", *identifierType)
	}

	// Get entitlements
	entitlementReq := &authorizationv2.GetEntitlementsRequest{
		EntityIdentifier: entityIdentifier,
	}

	entitlements, err := client.AuthorizationV2.GetEntitlements(
		context.Background(),
		entitlementReq,
	)
	if err != nil {
		return fmt.Errorf("failed to get entitlements: %w", err)
	}

	// Convert to JSON for output
	output, err := json.MarshalIndent(entitlements, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal entitlements: %w", err)
	}

	fmt.Println(string(output))
	return nil
}

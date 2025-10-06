package main

import (
	"context"
	"log"

	"github.com/opentdf/platform/sdk"
	authorizationv2 "github.com/opentdf/platform/protocol/go/authorization/v2"
	"github.com/opentdf/platform/protocol/go/entity"
)

func main() {
	platformEndpoint := "http://localhost:8080"

	// Create authenticated client
	client, err := sdk.New(
		platformEndpoint,
		sdk.WithClientCredentials("opentdf", "secret", nil),
	)
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	log.Println("OpenTDF SDK client created successfully!")
	log.Printf("Connected to platform at: %s", platformEndpoint)

	// Example: Get entitlements for a user
	getEntitlements(client)
}

func getEntitlements(client *sdk.SDK) {
	log.Println("\nGetting entitlements for user bob@OrgA.com...")

	entitlementReq := &authorizationv2.GetEntitlementsRequest{
		EntityIdentifier: &authorizationv2.EntityIdentifier{
			Identifier: &authorizationv2.EntityIdentifier_EntityChain{
				EntityChain: &entity.EntityChain{
					Entities: []*entity.Entity{
						{
							EphemeralId: "user-bob",
							EntityType: &entity.Entity_EmailAddress{
								EmailAddress: "bob@OrgA.com",
							},
						},
					},
				},
			},
		},
	}

	entitlements, err := client.AuthorizationV2.GetEntitlements(
		context.Background(),
		entitlementReq,
	)
	if err != nil {
		log.Printf("Warning: Failed to get entitlements (this is expected if platform is not running): %v", err)
		return
	}

	log.Printf("Successfully retrieved entitlements: %+v", entitlements)
}

# OpenTDF Hello World

A simple Go application demonstrating the use of the OpenTDF Platform SDK v0.11.0.

## Prerequisites

- Go 1.24 or higher
- Docker/Podman (to run OpenTDF platform locally)

## Setup

This project uses:
- `github.com/opentdf/platform/protocol/go v0.11.0` - Protocol definitions
- `github.com/opentdf/platform/sdk v0.8.0` - SDK for interacting with OpenTDF platform

Dependencies are already configured in `go.mod`.

## Running the Example

### 1. Start the OpenTDF Platform (optional)

If you want to test against a running platform:

```bash
# Follow the OpenTDF platform quickstart guide
docker compose --profile ers-test up
```

### 2. Build and Run

```bash
# Build the application
go build -o hello-world

# Run the application
./hello-world
```

**Note:** The example will attempt to connect to `http://localhost:8080`. If the platform is not running, you'll see a warning message, which is expected for demonstration purposes.

## What This Example Does

The hello-world application demonstrates:

1. **SDK Initialization**: Creates an authenticated OpenTDF SDK client with client credentials
2. **Authorization V2 API**: Uses the v2 authorization API to get entitlements for a user
3. **Entity Identification**: Shows how to construct entity identifiers using email addresses

## Code Structure

- [main.go](main.go) - Main application with SDK initialization and example usage
- `getEntitlements()` - Example function showing how to request entitlements for a user entity

## Key Concepts

### SDK Client Creation

```go
client, err := sdk.New(
    platformEndpoint,
    sdk.WithClientCredentials("opentdf", "secret", nil),
)
```

### Entity Identification

Entities can be identified using various methods:
- Email Address (as shown in the example)
- Username
- UUID
- Client ID
- JWT Token

### Authorization V2 API

The example uses the newer Authorization V2 API which provides:
- More flexible entity identification methods
- Support for entity chains
- Comprehensive hierarchy support

## Learn More

- [OpenTDF Platform](https://github.com/opentdf/platform)
- [OpenTDF Documentation](https://opentdf.io/)
- [SDK Authorization Guide](https://opentdf.io/sdks/authorization)

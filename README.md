# QRNG API - Go SDK

Official Go client library for the QRNG API.

## Installation

```bash
go get github.com/qrng-api/go-sdk
```

## Quick Start

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/qrng-api/go-sdk/qrng"
)

func main() {
	client := qrng.NewClient("your-api-key")
	
	result, err := client.Generate(&qrng.GenerateOptions{
		Bytes:  32,
		Format: "hex",
	})
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Random data: %s\n", result.Data)
	fmt.Printf("Proof ID: %s\n", result.ProofID)
}
```

## Usage

### Basic Generation

```go
client := qrng.NewClient("your-api-key")

// Generate 32 bytes of random hex data
result, err := client.Generate(&qrng.GenerateOptions{
	Bytes:  32,
	Format: "hex",
})
if err != nil {
	log.Fatal(err)
}

fmt.Println(result.Data)
```

### Quantum Methods

```go
result, err := client.Generate(&qrng.GenerateOptions{
	Bytes:  32,
	Format: "hex",
	Method: "photon", // photon, tunneling, vacuum, simulator
})
```

### Post-Quantum Signatures

```go
// Pro tier: Dilithium2
result, err := client.Generate(&qrng.GenerateOptions{
	Bytes:         32,
	Format:        "hex",
	SignatureType: "dilithium2",
})

// Enterprise tier: Dilithium3/5
result, err := client.Generate(&qrng.GenerateOptions{
	Bytes:         32,
	Format:        "hex",
	SignatureType: "dilithium3",
})
```

### Health Check

```go
health, err := client.Health()
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Status: %s\n", health.Status)
```

## API Reference

### NewClient

```go
func NewClient(apiKey string) *Client
```

Creates a new QRNG API client.

### Generate

```go
func (c *Client) Generate(opts *GenerateOptions) (*EntropyResult, error)
```

Generate random entropy.

**GenerateOptions:**
- `Bytes` (int): Number of bytes (1-1024)
- `Format` (string): Output format (`hex`, `base64`, `binary`, `uint8`, `uint32`)
- `Method` (string): Quantum method (`auto`, `photon`, `tunneling`, `vacuum`, `simulator`)
- `SignatureType` (string): Signature type (`ed25519`, `dilithium2`, `dilithium3`, `dilithium5`)

### Health

```go
func (c *Client) Health() (*HealthStatus, error)
```

Get system health status.

## License

MIT

## Support

- Documentation: https://qrngapi.com/docs
- GitHub: https://github.com/qrng-api/go-sdk
- Email: support@qrngapi.com

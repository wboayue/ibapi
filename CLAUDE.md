# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go implementation of the Interactive Brokers TWS (Trader Workstation) API, providing idiomatic Go interfaces to IB Gateway functionality. The library targets TWS API version 10.12 and is not backwards compatible with older versions.

Currently implements real-time market data functionality with plans to expand to order execution and account management.

## Key Architecture

The codebase follows a clean, layered architecture:

- **client.go**: Main API client with connection management and request routing
- **transport.go**: TCP socket communication layer handling raw message exchange with IB Gateway
- **encoders.go/decoders.go**: Protocol message encoding/decoding for IB's proprietary format
- **messages.go**: Message type constants and protocol definitions
- **models.go**: Data structures (Contract, Bar, TickData, etc.)
- **versions.go**: Server version compatibility constants

## Common Development Commands

### Build and Test
```bash
# Run all tests with coverage and race detection
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Run specific test
go test -v -run TestRealTimeBarsEncoder

# Build the example application
go build -o ibapi-example cmd/main.go

# Check for issues
go vet ./...
```

### Module Management
```bash
# Update dependencies
go mod tidy

# Download dependencies
go mod download

# Verify dependencies
go mod verify
```

### Release Process
```bash
# Create and push a new version tag
git tag vx.y.z
git push origin vx.y.z

# Publish to package repository
GOPROXY=proxy.golang.org go list -m github.com/wboayue/ibapi@vx.y.z
```

## Implementation Patterns

### Request/Response Pattern
All client methods follow a consistent pattern:
1. Generate unique request ID using `nextRequestId()`
2. Create response channel and register it in `client.channels` map
3. Encode and send request message
4. Handle responses asynchronously via goroutines
5. Clean up channels on context cancellation

### Thread Safety
- Use `sync.Mutex` for protecting shared state
- Separate mutexes for different concerns (requestIdMutex, contractDetailsMutex)
- Channels for safe concurrent message passing

### Error Handling
- Use `fmt.Errorf` with `%w` verb for error wrapping (Go 1.13+ standard)
- Return errors immediately, don't panic
- Check for IB error responses in message handlers

## Testing Approach

Tests use `github.com/stretchr/testify` for assertions. Test files follow Go convention with `*_test.go` suffix.

Key test patterns:
- Unit tests for encoders/decoders with known message formats
- Mock MessageBus interface for testing client logic without network
- Table-driven tests for multiple scenarios

## CI/CD

GitHub Actions workflow (`.github/workflows/ci.yml`) runs on all PRs and pushes to main:
- Tests with race detection enabled
- Coverage reporting to Codecov
- Go 1.17+ required
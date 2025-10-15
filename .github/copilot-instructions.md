# ResponsibleAPI-Go Copilot Instructions

## Project Overview
This is a Go-based JWT authentication library providing pluggable authentication providers (Basic Auth, API Key) with **storage-agnostic** persistence. The codebase follows a clean architecture with clear separation between authentication logic, providers, data models, and storage implementations.

## Architecture & Key Components

### Core Auth System (`auth/auth.go`)
- **AuthWrapper**: Main entry point containing provider, storage, and options
- **AuthInterface**: Contract for all auth providers with methods for token creation, validation, and credential decoding
- **AuthOptions**: Configuration struct with JWT settings, durations, and custom claims

### Provider Pattern (`service/`)
- **BasicAuth** (`service/basic_auth.go`): Base64-encoded username:password authentication
- **APIKeyAuth** (`service/api_key_auth.go`): API key-based authentication
- All providers implement `AuthInterface` and use injected storage for data operations

### Storage Interface (`storage/interface.go`)
- **UserStorage**: Interface for data operations - allows any storage backend
- **MySQL Implementation** (`storage/mysql/`): Production-ready GORM/MySQL implementation
- **In-Memory Implementation** (`storage/memory/`): Testing/simple apps implementation
- External apps implement `UserStorage` for custom storage (Redis, PostgreSQL, APIs, etc.)

### Token Management (`resource/access/`, `internal/`)
- **RToken** (`resource/access/rtoken.go`): Wrapper for JWT tokens with utility methods
- **Model/ResponseDTO** (`resource/access/model.go`): Data transfer objects for API responses
- **Validation** (`internal/validate.go`): JWT parsing and claims validation with custom `ClaimsGeneric`

## Development Workflows

### Build & Run Commands
```bash
# Dependencies
go mod tidy

# Run MySQL example (requires database)
go run cmd/api/main.go

# Run in-memory example (no database required)
go run examples/memory-storage/main.go

# Build binary
go build -o bin/api cmd/api/main.go

# Tests
go test ./...
```

### Storage Setup Options
1. **MySQL**: Execute `migration/schema.sql` + set DB_* environment variables
2. **In-Memory**: No setup required, includes sample data
3. **Custom**: Implement `storage.UserStorage` interface

## Project-Specific Patterns

### Storage-Agnostic Initialization Pattern
```go
// MySQL storage
db, _ := tools.NewDatabase()
storage := mysql.NewMySQLStorage(db)

// OR in-memory storage
storage := memory.NewInMemoryStorage()

// OR custom storage
storage := &YourCustomStorage{}

authService := auth.NewAuth(service.NewBasicAuth(), storage, auth.AuthOptions{...})
```

### Token Flow Pattern
1. **Decode credentials**: `Provider.Decode(encodedString)`
2. **Validate via storage**: `storage.FindUserByCredentials()` or `storage.FindUserByAPIKey()`
3. **Create tokens**: `Provider.CreateAccessToken()` and `Provider.CreateRefreshToken()`
4. **Build response**: Use `access.NewModel()` with builder methods (`WithAccessToken()`, etc.)

### Custom Storage Implementation
```go
type CustomStorage struct{}
func (c *CustomStorage) FindUserByCredentials(username, credentials string) (*user.User, error) {
    // Your storage logic (Redis, PostgreSQL, API calls, etc.)
}
// Implement other UserStorage interface methods...
```

## Key Integration Points

### Storage Interface Requirements
- `FindUserByCredentials(username, credentials string) (*user.User, error)`
- `FindUserByAPIKey(apiKey string) (*user.User, error)`
- `UpdateRefreshToken(userID, refreshToken string) error`
- `ValidateRefreshToken(refreshToken string) (*user.User, error)`

### Custom Claims Support
- `AuthOptions.CustomClaims` map allows arbitrary JWT claims
- `concerns.ClaimsGeneric` extends standard JWT claims with Role and Scopes
- Claims validation happens in `internal/validate.go` with configurable leeway

### Environment Configuration (MySQL only)
- Uses `godotenv` for .env file loading
- `envdecode` for struct-based environment variable binding
- Database configuration in `config/config.go` for MySQL implementation

## Critical Files for Understanding
- `storage/interface.go`: Storage abstraction contract
- `auth/auth.go`: Core interfaces and storage injection
- `storage/mysql/mysql.go`: Reference MySQL implementation
- `storage/memory/memory.go`: Simple in-memory implementation
- `examples/memory-storage/main.go`: Database-free usage example
- `cmd/api/main.go`: MySQL-based usage example
- `STORAGE.md`: Comprehensive storage implementation guide

## Migration & Adoption Notes
- **Breaking Change**: `auth.NewAuth()` now requires storage parameter
- **No Database Dependency**: Library is now storage-agnostic
- **Custom Storage**: External apps implement `UserStorage` for their specific needs
- **Backward Compatibility**: MySQL implementation maintains existing database schema
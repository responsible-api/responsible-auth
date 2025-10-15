# Storage Interface Documentation

## Overview

The ResponsibleAPI-Go library is **storage-agnostic**, meaning you can use any data storage backend by implementing the `UserStorage` interface. This eliminates the dependency on MySQL and allows integration with any database, in-memory storage, or external services.

## UserStorage Interface

```go
type UserStorage interface {
    // FindUserByCredentials retrieves a user by username/email and validates their credentials
    FindUserByCredentials(username, credentials string) (*user.User, error)
    
    // FindUserByAPIKey retrieves a user by their API key
    FindUserByAPIKey(apiKey string) (*user.User, error)
    
    // UpdateRefreshToken stores a refresh token for a user
    UpdateRefreshToken(userID string, refreshToken string) error
    
    // ValidateRefreshToken checks if a refresh token is valid for a user
    ValidateRefreshToken(refreshToken string) (*user.User, error)
}
```

## Usage

### With MySQL (Reference Implementation)

```go
package main

import (
    "github.com/responsible-api/responsible-auth/auth"
    "github.com/responsible-api/responsible-auth/service"
    "github.com/responsible-api/responsible-auth/storage/mysql"
    "github.com/responsible-api/responsible-auth/tools"
)

func main() {
    // Initialize database connection
    db, err := tools.NewDatabase()
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Create MySQL storage implementation
    storage := mysql.NewMySQLStorage(db)

    // Create auth service with storage
    authService := auth.NewAuth(service.NewBasicAuth(), storage, options)
}
```

### With In-Memory Storage

```go
package main

import (
    "github.com/responsible-api/responsible-auth/auth"
    "github.com/responsible-api/responsible-auth/service"
    "github.com/responsible-api/responsible-auth/storage/memory"
)

func main() {
    // Create in-memory storage (useful for testing)
    storage := memory.NewInMemoryStorage()

    // Create auth service with storage
    authService := auth.NewAuth(service.NewBasicAuth(), storage, options)
}
```

### Custom Storage Implementation

```go
package main

import (
    "github.com/responsible-api/responsible-auth/resource/user"
    "github.com/responsible-api/responsible-auth/storage"
)

// CustomStorage implements the UserStorage interface
type CustomStorage struct {
    // Your custom fields (Redis client, external API client, etc.)
}

func (c *CustomStorage) FindUserByCredentials(username, credentials string) (*user.User, error) {
    // Your custom implementation
    // Could query Redis, call external API, read from files, etc.
}

func (c *CustomStorage) FindUserByAPIKey(apiKey string) (*user.User, error) {
    // Your custom implementation
}

func (c *CustomStorage) UpdateRefreshToken(userID string, refreshToken string) error {
    // Your custom implementation
}

func (c *CustomStorage) ValidateRefreshToken(refreshToken string) (*user.User, error) {
    // Your custom implementation
}

func NewCustomStorage() storage.UserStorage {
    return &CustomStorage{}
}
```

## Available Implementations

### 1. MySQL Storage (`storage/mysql`)
- **Use case**: Production applications with MySQL database
- **Setup**: Requires MySQL database with schema from `migration/schema.sql`
- **Example**: See `cmd/api/main.go`

### 2. In-Memory Storage (`storage/memory`)
- **Use case**: Testing, development, simple applications
- **Setup**: No external dependencies required
- **Example**: See `examples/memory-storage/main.go`

## Implementation Guidelines

When implementing your own storage, ensure:

1. **Error Handling**: Return appropriate errors when users/tokens are not found
2. **Security**: Hash passwords appropriately, validate API keys securely
3. **Performance**: Implement efficient queries for your storage backend
4. **Consistency**: Maintain referential integrity between users and tokens

## Migration from Previous Versions

If you're upgrading from a version that had hardcoded MySQL dependency:

1. Replace `auth.NewAuth(provider, options)` with `auth.NewAuth(provider, storage, options)`
2. Choose a storage implementation or create your own
3. Remove direct database dependencies from your application code

## Examples

Complete working examples are available in:
- `cmd/api/main.go` - MySQL storage example
- `examples/memory-storage/main.go` - In-memory storage example

### Supporting storage types
Use MySQL (existing behavior)
- storage := mysql.NewMySQLStorage(db)

Use Redis
- storage := &RedisStorage{client: redisClient}

Use PostgreSQL
- storage := &PostgreSQLStorage{db: pgDb}

Use external APIs
- storage := &APIStorage{baseURL: "https://api.example.com"}

All with the same auth interface!
- authService := auth.NewAuth(service.NewBasicAuth(), storage, options)
# ResponsibleAPI-Go

A storage-agnostic JWT authentication library for Go applications with pluggable authentication providers and flexible data persistence.

## Features

- **Storage-Agnostic**: Works with any data storage (MySQL, PostgreSQL, Redis, in-memory, external APIs)
- **Pluggable Providers**: Basic Auth, API Key authentication (extensible for OAuth, LDAP, etc.)
- **JWT Tokens**: Access tokens and refresh tokens with custom claims
- **Clean Architecture**: Clear separation of concerns with dependency injection
- **Zero Database Dependencies**: Library core has no hardcoded storage requirements

## Installation

```bash
go get github.com/responsible-api/responsible-auth
```

## Quick Start

### Option 1: In-Memory Storage (No Database Required)

Perfect for testing, development, or simple applications:

```go
package main

import (
    "log"
    "time"
    
    "github.com/responsible-api/responsible-auth/auth"
    "github.com/responsible-api/responsible-auth/resource/access"
    "github.com/responsible-api/responsible-auth/service"
    "github.com/responsible-api/responsible-auth/storage/memory"
)

func main() {
    // 1. Create in-memory storage (includes sample data)
    storage := memory.NewInMemoryStorage()
    
    // 2. Initialize auth service
    authService := auth.NewAuth(service.NewBasicAuth(), storage, auth.AuthOptions{
        SecretKey:            "your-super-secure-secret-key-here", 
        TokenDuration:        5 * time.Hour,
        RefreshTokenDuration: 24 * 7 * time.Hour,
        TokenLeeway:          10 * time.Second,
        Issuer:               "https://your-app.com",
        Subject:              "your-app-user",
        
        // Add custom claims
        CustomClaims: map[string]interface{}{
            "organization": "your-org",
            "tier":         "premium",
        },
    })
    
    // 3. Authenticate user (test@example.com:password123)
    user, pass, err := authService.Provider.Decode("dGVzdEBleGFtcGxlLmNvbTpwYXNzd29yZDEyMw==")
    if err != nil {
        log.Fatalf("Failed to decode credentials: %v", err)
    }
    
    // 4. Generate access token
    accessToken, err := authService.Provider.CreateAccessToken(user, pass)
    if err != nil {
        log.Fatalf("Failed to create access token: %v", err)
    }
    
    // 5. Generate refresh token
    refreshToken, err := authService.Provider.CreateRefreshToken(user, pass)
    if err != nil {
        log.Fatalf("Failed to create refresh token: %v", err)
    }
    
    // 6. Build response
    expiry, _ := accessToken.GetExpirationTime()
    model := access.NewModel()
    model.WithAccessToken(accessToken.GetToken())
    model.WithRefreshToken(refreshToken.GetToken())
    model.WithExpiresIn(expiry.Unix())
    model.WithCreatedAt(time.Now().Unix())
    
    response := model.ToResponseDTO()
    
    log.Printf("🎉 Authentication successful!")
    log.Printf("Access Token: %s", response.AccessToken)
    log.Printf("Refresh Token: %s", response.RefreshToken)
    log.Printf("Expires In: %d seconds", response.ExpiresIn)
}
```

**Run it:**
```bash
go run examples/memory-storage/main.go
```

### Option 2: MySQL Storage (Production Ready)

For production applications with persistent storage:

```go
package main

import (
    "log"
    "time"
    
    "github.com/responsible-api/responsible-auth/auth"
    "github.com/responsible-api/responsible-auth/service"
    "github.com/responsible-api/responsible-auth/storage/mysql"
    "github.com/responsible-api/responsible-auth/tools"
)

func main() {
    // 1. Connect to MySQL database
    db, err := tools.NewDatabase()
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    
    // 2. Create MySQL storage implementation
    storage := mysql.NewMySQLStorage(db)
    
    // 3. Initialize auth service
    authService := auth.NewAuth(service.NewBasicAuth(), storage, auth.AuthOptions{
        SecretKey:            "your-super-secure-secret-key-here",
        TokenDuration:        5 * time.Hour,
        RefreshTokenDuration: 24 * 7 * time.Hour,
        TokenLeeway:          10 * time.Second,
        Issuer:               "https://your-app.com",
        Subject:              "your-app-user",
    })
    
    // 4. Same authentication flow as above...
    // (decode credentials, create tokens, build response)
}
```

**Setup MySQL:**
1. Execute `migration/schema.sql` to create database schema
2. Set environment variables:
   ```bash
   export DB_HOST="localhost"
   export DB_PORT="3306"
   export DB_USER="your_user"
   export DB_PASS="your_password"
   export DB_NAME="responsible_api"
   ```

**Run it:**
```bash
go run cmd/api/main.go
```

## Custom Storage Implementation

Create your own storage backend for Redis, PostgreSQL, external APIs, etc.:

```go
package main

import (
    "github.com/responsible-api/responsible-auth/resource/user"
    "github.com/responsible-api/responsible-auth/storage"
)

// RedisStorage implements UserStorage interface with Redis
type RedisStorage struct {
    client *redis.Client
}

func (r *RedisStorage) FindUserByCredentials(username, credentials string) (*user.User, error) {
    // Query Redis for user
    userData, err := r.client.HGetAll(ctx, "user:"+username).Result()
    if err != nil {
        return nil, err
    }
    
    // Validate credentials
    if userData["secret"] != credentials {
        return nil, errors.New("invalid credentials")
    }
    
    // Convert to user struct and return
    return &user.User{
        AccountID: parseUint64(userData["account_id"]),
        Name:      userData["name"],
        Mail:      userData["mail"],
        Secret:    userData["secret"],
        // ... other fields
    }, nil
}

func (r *RedisStorage) FindUserByAPIKey(apiKey string) (*user.User, error) {
    // Your Redis API key lookup logic
}

func (r *RedisStorage) UpdateRefreshToken(userID, refreshToken string) error {
    // Your Redis refresh token storage logic
}

func (r *RedisStorage) ValidateRefreshToken(refreshToken string) (*user.User, error) {
    // Your Redis refresh token validation logic
}

func NewRedisStorage(client *redis.Client) storage.UserStorage {
    return &RedisStorage{client: client}
}

// Usage
func main() {
    redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    storage := NewRedisStorage(redisClient)
    authService := auth.NewAuth(service.NewBasicAuth(), storage, options)
}
```

## API Key Authentication

Switch to API Key authentication instead of Basic Auth:

```go
// Use API Key provider instead of Basic Auth
authService := auth.NewAuth(service.NewApiKeyAuth(), storage, options)

// Authenticate with API key
token, err := authService.Provider.CreateAccessToken("", "your-api-key-here")
```

## Development Commands

```bash
# Install dependencies
go mod tidy

# Run tests
go test ./...

# Run MySQL example (requires database)
go run cmd/api/main.go

# Run in-memory example (no dependencies)
go run examples/memory-storage/main.go

# Build binary
go build -o bin/api cmd/api/main.go
```

## Project Structure

```
responsible-auth/
├── auth/                   # Core authentication logic
├── service/               # Authentication providers (Basic Auth, API Key)
├── storage/               # Storage interface and implementations
│   ├── interface.go       # UserStorage interface definition
│   ├── mysql/            # MySQL implementation
│   └── memory/           # In-memory implementation
├── resource/             # Data models and DTOs
├── internal/             # JWT token creation and validation
├── examples/             # Complete usage examples
├── migration/            # Database schema
└── tools/                # Database utilities
```

## Key Benefits

1. **Storage Flexibility**: Use any database or storage system
2. **Provider Extensibility**: Easy to add new authentication methods
3. **Zero Lock-in**: No vendor or database dependencies
4. **Production Ready**: Includes MySQL implementation with proper schema
5. **Developer Friendly**: In-memory option for quick testing
6. **Clean Architecture**: Clear separation of concerns

## Documentation

- **[Storage Guide](STORAGE.md)**: Comprehensive guide for implementing custom storage
- **[API Documentation](https://pkg.go.dev/github.com/responsible-api/responsible-auth)**: Complete API reference
- **Examples**: See `examples/` directory for complete implementations

## Migration from Previous Versions

If upgrading from MySQL-only versions:

```go
// Old
authService := auth.NewAuth(provider, options)

// New (with MySQL)
db, _ := tools.NewDatabase()
storage := mysql.NewMySQLStorage(db)
authService := auth.NewAuth(provider, storage, options)

// New (with in-memory for testing)
storage := memory.NewInMemoryStorage()
authService := auth.NewAuth(provider, storage, options)
```

## Contributing

Feedback and suggestions are welcome! Please open an issue for any bugs or feature requests.

## License

[MIT License](LICENSE)

Author of this project [@vince-scarpa](https://github.com/vince-scarpa)
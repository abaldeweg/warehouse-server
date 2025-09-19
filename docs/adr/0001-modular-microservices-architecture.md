# ADR-0001: Modular Microservices Architecture with Go Workspace

## Status

Accepted

## Context

The warehouse-server is a database management system for warehouse operations with multiple distinct concerns:

1. **Core API functionality** - Book inventory management, reservations, user authentication
2. **File serving** - Static file serving for uploads and assets  
3. **Log processing** - Import and web interface for application logs
4. **Shared framework** - Common utilities, middleware, and patterns

The system needs to:
- Scale different components independently
- Maintain clear separation of concerns
- Share common code efficiently
- Support different deployment strategies (containerized microservices)
- Provide consistent API patterns across services

## Decision

We will implement a **modular microservices architecture** using Go's workspace feature with the following structure:

```
warehouse-server/
├── go.work                 # Go workspace configuration
├── framework/              # Shared framework module
├── gateway/                # Main API gateway and core business logic
├── static/                 # Static file server microservice
├── logs_import/            # Log processing microservice  
├── logs_web/               # Log web interface microservice
```

### Key Architectural Components

#### 1. Go Workspace Structure
- Use `go.work` to manage multiple related modules in a single repository
- Each module has its own `go.mod` file for independent versioning
- Shared framework module provides common functionality

#### 2. Framework Module (`framework/`)
Provides shared utilities and patterns:
- **API Key Authentication** (`framework/apikey`) - Centralized API key management
- **CORS Handling** (`framework/cors`) - Cross-origin request middleware
- **Configuration** (`framework/config`) - Unified configuration loading with Viper
- **Router Setup** (`framework/router`) - Standard Gin router with middleware
- **Storage Abstraction** (`framework/storage`) - Pluggable storage (filesystem/cloud)

#### 3. Gateway Service (`gateway/`)
Main API service with layered architecture:
- **Controllers** - HTTP request handling and response formatting
- **Repository Pattern** - Data access abstraction with GORM
- **Models** - Domain entities with validation
- **Multi-tenancy** - Branch-based data isolation
- **Authentication** - Integration with external auth service
- **File Upload** - Image processing and cover management

#### 4. Microservices
- **Static Service** (`static/`) - Simple HTTP file server for assets
- **Logs Import** (`logs_import/`) - Background log processing with MongoDB
- **Logs Web** (`logs_web/`) - Web interface for log viewing

## Consequences

### Positive
- **Independent Deployment**: Each service can be deployed and scaled independently
- **Clear Boundaries**: Well-defined interfaces between components
- **Code Reuse**: Framework module eliminates duplication
- **Technology Flexibility**: Services can use different databases (SQLite/MySQL vs MongoDB)
- **Team Autonomy**: Teams can work on different services independently
- **Container Ready**: Each service has its own Dockerfile for containerization

### Negative
- **Complexity**: More complex than a monolith for simple operations
- **Network Latency**: Inter-service communication overhead
- **Distributed System Concerns**: Need to handle service discovery, circuit breakers, etc.
- **Testing Complexity**: Integration testing across services is more complex

### Risks and Mitigations
- **Service Communication**: Currently uses HTTP proxy - consider service mesh for production
- **Data Consistency**: Repository pattern provides transaction boundaries within services
- **Configuration Management**: Framework config module provides consistent environment handling
- **Authentication**: Centralized API key system handles service-to-service auth

## Implementation Notes

### Repository Pattern
Each service uses the repository pattern for data access:
```go
type ReservationRepository interface {
    FindAll(uint) ([]models.Reservation, error)
    FindOne(id uuid.UUID) (*models.Reservation, error)
    Create(reservation *models.Reservation) error
    Update(reservation *models.Reservation) error
    Delete(id uuid.UUID) error
}
```

### Middleware Architecture
Standard middleware stack using Gin:
- CORS handling
- API key authentication  
- Permission-based authorization
- Request/response logging

### Configuration Strategy
- Environment variables for service-specific config
- Viper for configuration file loading
- Framework config module for consistent patterns

### Multi-tenancy
- Branch-based data isolation in the gateway service
- User context passed through middleware
- Repository methods filter by branch ID

## Related Decisions

- [ADR-0002: Repository Pattern for Data Access](0002-repository-pattern.md) (Future)
- [ADR-0003: API Key Authentication System](0003-api-key-authentication.md) (Future)
- [ADR-0004: Storage Abstraction Layer](0004-storage-abstraction.md) (Future)
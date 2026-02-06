# Go Context-Based Logger

A lightweight, context-aware logging library for Go that provides structured JSON logging with Go's standard `context.Context`.

## Features

- 🎯 **Context-aware logging** - Automatically includes contextual information in log entries
- 📝 **Structured JSON output** - All logs are output as JSON for easy parsing and analysis
- 🔄 **Context propagation** - Uses Go's standard `context.Context` for passing log context
- 🏷️ **Flexible tagging** - Add tags for categorisation and filtering
- 📊 **Metadata support** - Include custom key-value pairs in logs
- 🎫 **Session tracking** - Track requests/operations with session IDs
- 🛡️ **Type-safe** - Fully typed with Go's type system
- ⚡ **Zero dependencies** - Only uses Go standard library

## Installation

```bash
cd golang
go mod tidy
```

## Quick Start

### Basic Usage

```go
package main

import "github.com/peterzzshi/context-based-logger"

func main() {
    // Simple logging at different levels
    logger.Info("Application started")
    logger.Warn("This is a warning message")
    logger.Debug("Debug information")
    logger.Error("Something went wrong")
}
```

### Context-Aware Logging

```go
package main

import (
    "context"
    "github.com/peterzzshi/context-based-logger/logger"
)

func main() {
    // Create a context with session ID, category, tags, and metadata
    ctx := context.Background()
    logCtx := logger.NewLogContext(logger.LogContextData{}).
        WithSessionID("req-123").
        WithCategory("api-request").
        WithTags("user-service", "authentication").
        WithMetadata(map[string]string{
            "userId":   "456",
            "endpoint": "/api/login",
        })
    
    ctx = logger.WithLogContext(ctx, logCtx)
    log := logger.New(ctx)
    
    // All logs will include the context data
    log.Info("Processing login request")
    log.Debug("Validating credentials")
    log.Info("Login successful")
}
```

## API Reference

### Logger Methods

The logger provides four log levels:

```go
log.Debug(args ...interface{})
log.Info(args ...interface{})
log.Warn(args ...interface{})
log.Error(args ...interface{})
```

### LogContext

#### Creating Contexts

```go
// Create a new empty context
logCtx := logger.NewLogContext(logger.LogContextData{})

// Or initialize with data
logCtx := logger.NewLogContext(logger.LogContextData{
    SessionID: "req-123",
    Category:  "api",
    Tags:      map[string]bool{"user-service": true},
    Metadata:  map[string]string{"userId": "456"},
})
```

#### Builder Methods (Immutable)

```go
// Add category
logCtx = logCtx.WithCategory("api-request")

// Add session ID
logCtx = logCtx.WithSessionID("req-123")

// Add tags
logCtx = logCtx.WithTags("api", "database")

// Remove tags
logCtx = logCtx.WithoutTags("database")

// Add metadata
logCtx = logCtx.WithMetadata(map[string]string{
    "userId": "456",
    "endpoint": "/api/users",
})

// Remove metadata
logCtx = logCtx.WithoutMetadata("endpoint")
```

#### Context Operations

```go
// Add LogContext to context.Context
ctx = logger.WithLogContext(ctx, logCtx)

// Retrieve LogContext from context.Context
logCtx = logger.GetLogContext(ctx)

// Create a logger with context
log := logger.New(ctx)
```

## Usage with HTTP Servers

Perfect for use with Go web frameworks:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    logCtx := logger.NewLogContext(logger.LogContextData{}).
        WithSessionID(r.Header.Get("X-Request-ID")).
        WithCategory("http").
        WithTags("api").
        WithMetadata(map[string]string{
            "method": r.Method,
            "path":   r.URL.Path,
        })
    
    ctx := logger.WithLogContext(r.Context(), logCtx)
    log := logger.New(ctx)
    
    log.Info("Processing request")
    
    // Pass context to other functions
    processRequest(ctx)
}
```

## Running Examples

```bash
cd examples
go run main.go
```

## Running Tests

```bash
cd logger
go test -v
go test -bench=.
```

## Example Output

```json
{
  "level": "info",
  "message": "Processing user request",
  "sessionId": "req-123",
  "details": {
    "category": "api-request",
    "tags": ["api", "user-service"],
    "metadata": {
      "userId": "456",
      "endpoint": "/api/users"
    },
    "timestamp": "2026-02-06T10:30:45Z"
  }
}
```

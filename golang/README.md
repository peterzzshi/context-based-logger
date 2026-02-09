# Go Context-Based Logger

Context-aware logging for Go using `context.Context` with structured JSON output.

## Usage

Copy the `logger/` folder into your project.

```go
import "yourproject/logger"

ctx := context.Background()
logCtx := logger.NewLogContext(logger.LogContextData{}).
    WithSessionID("req-123").
    WithTags("api").
    WithMetadata(map[string]string{"userId": "456"})

_, err := logger.WithLogContext(ctx, logCtx, func(ctx context.Context) (struct{}, error) {
    logger.Info(ctx, "Processing request")
    return struct{}{}, nil
})
```

## API

### Creating Context

```go
logCtx := logger.NewLogContext(logger.LogContextData{}).
    WithCategory("api").
    WithSessionID("req-123").
    WithTags("tag1", "tag2").
    WithMetadata(map[string]string{"key": "value"}).
    WithoutTags("old-tag").
    WithoutMetadata("old-key")
```

### Using Context

The callback-based approach ensures context is properly scoped:

```go
_, err := logger.WithLogContext(ctx, logCtx, func(ctx context.Context) (struct{}, error) {
    logger.Info(ctx, "Processing")
    
    // You can nest contexts for more specific operations
    enrichedLogCtx := logger.GetLogContext(ctx).WithTags("database")
    _, _ = logger.WithLogContext(ctx, enrichedLogCtx, func(ctx context.Context) (struct{}, error) {
        logger.Debug(ctx, "Database operation")
        return struct{}{}, nil
    })
    
    return struct{}{}, nil
})
```

**With return values:**
```go
result, err := logger.WithLogContext(ctx, logCtx, func(ctx context.Context) (int, error) {
    logger.Info(ctx, "Calculating")
    return 42, nil
})
```

**Retrieving context:**
```go
func nested(ctx context.Context) {
    logCtx := logger.GetLogContext(ctx)
    // Use logCtx to create new enriched context
    enriched := logCtx.WithTags("additional-tag")
    _, _ = logger.WithLogContext(ctx, enriched, func(ctx context.Context) (struct{}, error) {
        logger.Info(ctx, "Nested operation")
        return struct{}{}, nil
    })
}
```

### Log Levels

```go
logger.Debug(ctx, "Debug message")
logger.Info(ctx, "Info message")
logger.Warn(ctx, "Warning message")
logger.Error(ctx, "Error message")

// Multiple arguments
logger.Info(ctx, "User", "logged in")

// With error (stack trace extracted)
logger.Error(ctx, "Failed to process", err)
```

## Example

See `examples/main.go`:

```bash
go run examples/main.go
```

## Testing

```bash
go test ./logger -v
```

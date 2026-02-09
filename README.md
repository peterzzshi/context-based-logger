# Context-Based Logger

Lightweight context-aware logging for **TypeScript** and **Go** with structured JSON output.

## Features

- Context-aware logging with automatic propagation
- Structured JSON output
- Flexible tagging and metadata
- Session tracking
- Type-safe
- Zero external dependencies

## Usage

Copy the `logger/` folder into your project.

### TypeScript

```typescript
import { logger } from './logger/logger';
import { LogContext, withLogContext } from './logger/context';

const logCtx = LogContext.create()
  .withSessionId('req-123')
  .withMetadata({ userId: '456' });

withLogContext(logCtx, () => {
  logger.info('Processing request');
});
```

### Go

```go
import "yourproject/logger"

ctx := context.Background()
logCtx := logger.NewLogContext(logger.LogContextData{}).
    WithSessionID("req-123").
    WithMetadata(map[string]string{"userId": "456"})

_, err := logger.WithLogContext(ctx, logCtx, func(ctx context.Context) (struct{}, error) {
    logger.Info(ctx, "Processing request")
    return struct{}{}, nil
})
```

## Project Structure

```
context-based-logger/
├── typescript/
│   └── src/logger/     # Copy this folder
├── golang/
│   └── logger/         # Copy this folder
└── README.md
```

## Examples

See `typescript/examples/demo.ts` and `golang/examples/main.go` for complete examples.

## API

### Context Methods

Both implementations provide immutable context builders:

- `WithCategory(category)` / `withCategory(category)`
- `WithSessionID(id)` / `withSessionId(id)`
- `WithTags(...)` / `withTags(...)`
- `WithoutTags(...)` / `withoutTags(...)`
- `WithMetadata(map)` / `withMetadata(obj)`
- `WithoutMetadata(...)` / `withoutMetadata(...)`

### Logger Methods

- `Debug(msg)` / `logger.debug(msg)`
- `Info(msg)` / `logger.info(msg)`
- `Warn(msg)` / `logger.warn(msg)`
- `Error(msg)` / `logger.error(msg)`

## Output Format

```json
{
  "level": "info",
  "message": "Processing request",
  "sessionId": "req-123",
  "details": {
    "category": "api",
    "tags": ["user-service"],
    "metadata": {"userId": "456"},
    "timestamp": "2026-02-09T06:13:24Z"
  }
}
```



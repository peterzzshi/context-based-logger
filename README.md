# Context-Based Logger

A lightweight, context-aware logging library available in both **TypeScript** and **Go** that provides structured JSON logging with context management.

## Features

- 🎯 **Context-aware logging** - Automatically includes contextual information in log entries
- 📝 **Structured JSON output** - All logs are output as JSON for easy parsing and analysis
- 🔄 **Context propagation** - TypeScript uses `AsyncLocalStorage`, Go uses `context.Context`
- 🏷️ **Flexible tagging** - Add tags for categorization and filtering
- 📊 **Metadata support** - Include custom key-value pairs in logs
- 🎫 **Session tracking** - Track requests/operations with session IDs
- 🛡️ **Type-safe** - Fully typed implementations
- ⚡ **Zero external dependencies** - Only uses standard libraries

## Project Structure

```
context-based-logger/
├── typescript/          # TypeScript implementation
│   ├── src/
│   │   └── logger/     # Logger source code
│   ├── examples/       # Usage examples
│   └── README.md       # TypeScript-specific docs
├── golang/             # Go implementation
│   ├── logger/         # Logger package
│   ├── examples/       # Usage examples
│   └── README.md       # Go-specific docs
└── README.md           # This file
```

## Quick Start

### TypeScript

```bash
cd typescript
npm install
npm run dev
```

See [typescript/README.md](typescript/README.md) for detailed TypeScript documentation.

### Go

```bash
cd golang
go run examples/main.go
```

See [golang/README.md](golang/README.md) for detailed Go documentation.

## Language-Specific Examples

### TypeScript Example

```typescript
import { logger } from './logger/logger';
import { LogContext, withLogContext } from './logger/context';

// Basic logging
logger.info('Application started');

// Context-aware logging
const logCtx = LogContext.create({
  sessionId: 'req-123',
  category: 'api',
  tags: new Set(['user-service']),
  metadata: new Map([['userId', '456']])
});

withLogContext(logCtx, () => {
  logger.info('Processing request');
});
```

### Go Example

```go
package main

import (
    "context"
    "github.com/yourusername/context-based-logger/logger"
)

func main() {
    // Basic logging
    logger.Info("Application started")
    
    // Context-aware logging
    ctx := context.Background()
    logCtx := logger.NewLogContext(logger.LogContextData{}).
        WithSessionID("req-123").
        WithCategory("api").
        WithTags("user-service").
        WithMetadata(map[string]string{"userId": "456"})
    
    ctx = logger.WithLogContext(ctx, logCtx)
    log := logger.New(ctx)
    log.Info("Processing request")
}
```

## Key Concepts

### TypeScript Implementation
- Uses Node.js `AsyncLocalStorage` for automatic context propagation
- Context flows through async/await operations automatically
- Immutable context with builder methods

### Go Implementation  
- Uses Go's standard `context.Context` for explicit context passing
- Context must be passed through function parameters
- Immutable context with builder methods
- Perfect for HTTP handlers and microservices

## Example Output

Both implementations produce structured JSON logs:

```json
{
  "level": "info",
  "message": "Processing request",
  "sessionId": "req-123",
  "details": {
    "category": "api",
    "tags": ["user-service"],
    "metadata": {"userId": "456"},
    "timestamp": "2026-02-06T10:30:45Z"
  }
}
```

## Testing

### TypeScript
```bash
cd typescript
npm install
npm test
```

### Go
```bash
cd golang/logger
go test -v
go test -bench=.
```

## License

MIT

## Contributing

Contributions are welcome! Please see the language-specific README files for detailed documentation and development guidelines.

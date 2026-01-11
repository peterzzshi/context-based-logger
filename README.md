# Context-Based Logger

A lightweight, context-aware logging library for TypeScript that provides structured JSON logging with async context management using Node.js's `AsyncLocalStorage`.

## Features

- 🎯 **Context-aware logging** - Automatically includes contextual information in log entries
- 📝 **Structured JSON output** - All logs are output as JSON for easy parsing and analysis
- 🔄 **Async context preservation** - Context persists across async operations using `AsyncLocalStorage`
- 🏷️ **Flexible tagging** - Add tags for categorisation and filtering
- 📊 **Metadata support** - Include custom key-value pairs in logs
- 🎫 **Session tracking** - Track requests/operations with session IDs
- 🛡️ **Fully type-safe** - Written in TypeScript with complete type definitions
- ⚡ **Zero dependencies** - Only uses Node.js built-in modules

## Installation

```bash
npm install context-based-logger
```

Or clone for local development:

```bash
git clone https://github.com/peterzzshi/context-based-logger
cd context-based-logger
npm install
npm run build
```

## Quick Start

### Basic Usage

```typescript
import { logger } from 'context-based-logger/logger/logger';

// Simple logging at different levels
logger.info('Application started');
logger.warn('This is a warning message');
logger.debug('Debug information');

// Logging errors - pass message and Error object
logger.error('Something went wrong', new Error('Connection failed'));

// You can also pass just an Error
logger.error(new Error('Database connection failed'));
```

### Context-Aware Logging

```typescript
import { logger } from 'context-based-logger/logger/logger';
import { LogContext, withLogContext } from 'context-based-logger/logger/context';

// Create a context with session ID, category, tags, and metadata
const requestContext = LogContext.create({
  sessionId: 'req-123',
  category: 'api-request',
  tags: new Set(['user-service', 'authentication']),
  metadata: new Map([
    ['userId', '456'],
    ['endpoint', '/api/login']
  ])
});

// All logs within this context will include the context data
withLogContext(requestContext, () => {
  logger.info('Processing login request');
  logger.debug('Validating credentials');
  logger.info('Login successful');
});
```

## API Reference

### Logger Methods

The logger provides four log levels. Each method accepts either a single value or a message with an additional value (typically an Error):

```typescript
// Single argument - can be any value
logger.debug(message: unknown): void
logger.info(message: unknown): void
logger.warn(message: unknown): void
logger.error(message: unknown): void

// Two arguments - message with additional context (e.g., Error)
logger.debug(message: string, value: unknown): void
logger.info(message: string, value: unknown): void
logger.warn(message: string, value: unknown): void
logger.error(message: string, error: Error): void
```

### LogContext

#### Creating Contexts

```typescript
import { LogContext } from 'context-based-logger/logger/context';

// Empty context
const context = LogContext.create();

// Context with initial data
const context = LogContext.create({
  sessionId: 'session-123',
  category: 'api',
  tags: new Set(['database', 'user']),
  metadata: new Map([['operation', 'CREATE']])
});
```

#### Context Methods

All context methods are immutable and return a new `LogContext` instance:

```typescript
// Add or update session ID
context.withSessionId(sessionId: string): LogContext

// Set category (can be undefined to clear)
context.withCategory(category: string | undefined): LogContext

// Add tags (merges with existing tags)
context.withTags(...tags: string[]): LogContext

// Remove specific tags
context.withoutTags(...tags: string[]): LogContext

// Add metadata (merges with existing metadata)
context.withMetadata(metadata: Record<string, string>): LogContext

// Remove specific metadata keys
context.withoutMetadata(...keys: string[]): LogContext
```

#### Using Contexts

```typescript
import { logger } from 'context-based-logger/logger/logger';
import { LogContext, withLogContext, getLogContext } from 'context-based-logger/logger/context';

// Execute code within a context (works with sync and async functions)
withLogContext(context, () => {
  logger.info('This log includes context');
});

// Async example
await withLogContext(context, async () => {
  await someAsyncOperation();
  logger.info('Context is preserved across await');
});

// Get current context (returns empty context if none is set)
const currentContext = getLogContext();
```

### Building Context Fluently

You can chain context methods for fluent context building:

```typescript
const context = LogContext.create()
  .withSessionId(`req-${Date.now()}`)
  .withCategory('api')
  .withTags('user-service', 'database')
  .withMetadata({
    userId: '123',
    endpoint: '/api/user',
    method: 'GET'
  });
```

## Type Definitions

### LogContextData

```typescript
interface LogContextData {
  readonly tags: ReadonlySet<string>;
  readonly category: string | undefined;
  readonly metadata: ReadonlyMap<string, string>;
  readonly sessionId: string | undefined;
}
```

### LogOutput

The structure of JSON log output:

```typescript
interface LogOutput {
  readonly level: 'debug' | 'info' | 'warn' | 'error';
  readonly message?: unknown;
  readonly sessionId?: string;
  readonly details: {
    readonly tags?: string[];
    readonly category?: string;
    readonly metadata?: Record<string, string>;
    readonly stack?: string;  // Included when logging Errors
    readonly timestamp: string;
  };
}
```

### Logger

```typescript
interface Logger {
  debug(...parameters: LoggingParameters): void;
  info(...parameters: LoggingParameters): void;
  warn(...parameters: LoggingParameters): void;
  error(...parameters: LoggingParameters): void;
}
```

## Log Output Examples

### Basic log (no context)

```json
{
  "level": "info",
  "message": "Application started",
  "details": {
    "timestamp": "2025-01-20T10:30:45.123Z"
  }
}
```

### Log with context

```json
{
  "level": "info",
  "message": "User login successful",
  "sessionId": "req-123",
  "details": {
    "tags": ["authentication", "user-service"],
    "category": "api-request",
    "metadata": {
      "userId": "456",
      "endpoint": "/api/login"
    },
    "timestamp": "2025-01-20T10:30:45.123Z"
  }
}
```

### Error log with stack trace

```json
{
  "level": "error",
  "message": "Database query failed Error: Connection timeout",
  "sessionId": "req-123",
  "details": {
    "tags": ["database"],
    "category": "db-operation",
    "stack": "Error: Connection timeout\n    at ...",
    "timestamp": "2025-01-20T10:30:45.123Z"
  }
}
```

## Usage Examples

### Express.js Middleware

```typescript
import express from 'express';
import { logger } from 'context-based-logger/logger/logger';
import { LogContext, withLogContext, getLogContext } from 'context-based-logger/logger/context';

const app = express();

// Request logging middleware
app.use((req, res, next) => {
  const context = LogContext.create({
    sessionId: `req-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
    category: 'http-request',
    tags: new Set(['express', 'api']),
    metadata: new Map([
      ['method', req.method],
      ['url', req.url],
      ['userAgent', req.get('User-Agent') || 'unknown']
    ])
  });

  withLogContext(context, () => {
    logger.info(`${req.method} ${req.url}`);
    next();
  });
});

app.get('/users/:id', async (req, res) => {
  const userId = req.params.id;
  
  // Enrich current context with user-specific data
  const userContext = getLogContext()
    .withMetadata({ userId })
    .withTags('user-service');
  
  await withLogContext(userContext, async () => {
    logger.info('Fetching user data');
    
    try {
      const user = await getUserById(userId);
      logger.info('User data retrieved successfully');
      res.json(user);
    } catch (error) {
      logger.error('Failed to fetch user', error as Error);
      res.status(500).json({ error: 'Internal server error' });
    }
  });
});
```

### Database Operations

```typescript
async function executeQuery<T>(sql: string, params: unknown[]): Promise<T> {
  const dbContext = getLogContext()
    .withTags('database')
    .withMetadata({ 
      operation: sql.split(' ')[0].toUpperCase(),
      table: extractTableName(sql)
    });

  return withLogContext(dbContext, async () => {
    logger.debug('Executing database query');
    
    try {
      const result = await db.query<T>(sql, params);
      logger.info('Query executed successfully');
      return result;
    } catch (error) {
      logger.error('Database query failed', error as Error);
      throw error;
    }
  });
}
```

### Background Job Processing

```typescript
async function processJob(jobId: string, jobData: JobData): Promise<void> {
  const jobContext = LogContext.create()
    .withSessionId(jobId)
    .withCategory('background-job')
    .withTags('worker', 'async')
    .withMetadata({
      jobType: jobData.type,
      priority: String(jobData.priority)
    });

  await withLogContext(jobContext, async () => {
    logger.info('Starting job processing');
    
    try {
      logger.debug('Validating job data');
      await validateJobData(jobData);
      
      logger.debug('Executing job logic');
      await executeJobLogic(jobData);
      
      logger.info('Job completed successfully');
    } catch (error) {
      logger.error('Job processing failed', error as Error);
      throw error;
    }
  });
}
```

### Custom Context Factories

Create reusable context factories for your domain:

```typescript
import { Request } from 'express';
import { LogContext, getLogContext } from 'context-based-logger/logger/context';

// Factory for HTTP request contexts
export function createRequestContext(req: Request): LogContext {
  return LogContext.create({
    sessionId: (req.headers['x-request-id'] as string) || `req-${Date.now()}`,
    category: 'http',
    tags: new Set(['api', req.method.toLowerCase()]),
    metadata: new Map([
      ['method', req.method],
      ['path', req.path],
      ['ip', req.ip || 'unknown']
    ])
  });
}

// Factory for database operation contexts
export function createDbContext(operation: string, table: string): LogContext {
  return getLogContext()
    .withTags('database', operation.toLowerCase())
    .withMetadata({ operation, table });
}

// Factory for external service call contexts
export function createServiceContext(serviceName: string, operation: string): LogContext {
  return getLogContext()
    .withTags('external-service', serviceName)
    .withMetadata({ service: serviceName, operation });
}
```

## Development

### Running the Demo

```bash
npm install
npm run dev        # Run demo with tsx
npm run demo       # Build and run compiled demo
```

### Building

```bash
npm run build      # Compile TypeScript to dist/
npm run clean      # Remove dist/ directory
```

### Testing

```bash
npm test                # Run tests
npm run test:watch      # Run tests in watch mode
npm run test:coverage   # Run tests with coverage report
```

### Code Quality

```bash
npm run lint           # Run ESLint
npm run lint:fix       # Fix ESLint issues
npm run format         # Format with Prettier
npm run format:check   # Check formatting
npm run type-check     # TypeScript type checking
npm run validate       # Run all checks (type-check, lint, test)
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

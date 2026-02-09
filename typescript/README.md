# TypeScript Context-Based Logger

Context-aware logging for TypeScript using `AsyncLocalStorage` with structured JSON output.

## Usage

Copy the `src/logger/` folder into your project.

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

## API

### Creating Context

```typescript
const logCtx = LogContext.create()
  .withCategory('api')
  .withSessionId('req-123')
  .withTags('tag1', 'tag2')
  .withMetadata({ key: 'value' })
  .withoutTags('old-tag')
  .withoutMetadata('old-key');
```

### Using Context

```typescript
// With callback (returns value)
const result = withLogContext(logCtx, () => {
  logger.info('Processing');
  return 'done';
});

// Retrieve context in nested functions
function nested() {
  const logCtx = getLogContext();
  logger.info('Message');
}
```

### Log Levels

```typescript
logger.debug('Debug message');
logger.info('Info message');
logger.warn('Warning message');
logger.error('Error message', new Error('details'));
```

## Example

See `examples/demo.ts`:

```bash
npm install
npm run dev
```

## Testing

```bash
npm test
```

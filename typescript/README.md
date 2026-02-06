# TypeScript Context-Based Logger

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
npm install
npm run build
```

## Quick Start

### Basic Usage

```typescript
import { logger } from './logger/logger';

// Simple logging at different levels
logger.info('Application started');
logger.warn('This is a warning message');
logger.debug('Debug information');

// Logging errors - pass message and Error object
logger.error('Something went wrong', new Error('Connection failed'));
```

### Context-Aware Logging

```typescript
import { logger } from './logger/logger';
import { LogContext, withLogContext } from './logger/context';

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

## Running Examples

```bash
npm run dev
```

## Running Tests

```bash
npm test
```

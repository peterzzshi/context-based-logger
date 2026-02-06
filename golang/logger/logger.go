package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// contextLogger implements the Logger interface with context support
type contextLogger struct {
	ctx context.Context
}

// New creates a new Logger with the given context
func New(ctx context.Context) Logger {
	return &contextLogger{ctx: ctx}
}

// Debug logs a debug-level message
func (cl *contextLogger) Debug(args ...interface{}) {
	cl.log(LevelDebug, args...)
}

// Info logs an info-level message
func (cl *contextLogger) Info(args ...interface{}) {
	cl.log(LevelInfo, args...)
}

// Warn logs a warning-level message
func (cl *contextLogger) Warn(args ...interface{}) {
	cl.log(LevelWarn, args...)
}

// Error logs an error-level message
func (cl *contextLogger) Error(args ...interface{}) {
	cl.log(LevelError, args...)
}

// log creates and writes a log entry
func (cl *contextLogger) log(level LogLevel, args ...interface{}) {
	logContext := GetLogContext(cl.ctx)

	output := LogOutput{
		Level:   level,
		Details: make(map[string]interface{}),
	}

	// Add session ID if present
	if logContext.data.SessionID != "" {
		output.SessionID = logContext.data.SessionID
	}

	// Add tags if present
	if len(logContext.data.Tags) > 0 {
		tags := make([]string, 0, len(logContext.data.Tags))
		for tag := range logContext.data.Tags {
			tags = append(tags, tag)
		}
		sort.Strings(tags) // Sort for consistent output
		output.Details["tags"] = tags
	}

	// Add category if present
	if logContext.data.Category != "" {
		output.Details["category"] = logContext.data.Category
	}

	// Add metadata if present
	if len(logContext.data.Metadata) > 0 {
		output.Details["metadata"] = logContext.data.Metadata
	}

	// Add timestamp
	output.Details["timestamp"] = time.Now().UTC().Format(time.RFC3339)

	// Format message
	if len(args) > 0 {
		output.Message = fmt.Sprint(args...)
	}

	// Marshal to JSON and write to stdout
	jsonBytes, err := json.Marshal(output)
	if err != nil {
		// Fallback if JSON marshaling fails
		fmt.Fprintf(os.Stderr, "Failed to marshal log: %v\n", err)
		return
	}

	fmt.Fprintln(os.Stdout, string(jsonBytes))
}

// Default logger instance for convenience functions
var defaultLogger = New(context.Background())

// Debug logs a debug-level message using the default logger
func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

// Info logs an info-level message using the default logger
func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

// Warn logs a warning-level message using the default logger
func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

// Error logs an error-level message using the default logger
func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

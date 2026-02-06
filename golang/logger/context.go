package logger

import "context"

type contextKey string

const logContextKey contextKey = "logContext"

// LogContext wraps the context data and provides immutable operations
type LogContext struct {
	data LogContextData
}

// NewLogContext creates a new LogContext with the given data
func NewLogContext(data LogContextData) *LogContext {
	if data.Tags == nil {
		data.Tags = make(map[string]bool)
	}
	if data.Metadata == nil {
		data.Metadata = make(map[string]string)
	}
	return &LogContext{data: data}
}

// WithCategory returns a new LogContext with the specified category
func (lc *LogContext) WithCategory(category string) *LogContext {
	newData := lc.copyData()
	newData.Category = category
	return &LogContext{data: newData}
}

// WithSessionID returns a new LogContext with the specified session ID
func (lc *LogContext) WithSessionID(sessionID string) *LogContext {
	newData := lc.copyData()
	newData.SessionID = sessionID
	return &LogContext{data: newData}
}

// WithTags returns a new LogContext with additional tags
func (lc *LogContext) WithTags(tags ...string) *LogContext {
	newData := lc.copyData()
	for _, tag := range tags {
		newData.Tags[tag] = true
	}
	return &LogContext{data: newData}
}

// WithoutTags returns a new LogContext with specified tags removed
func (lc *LogContext) WithoutTags(tags ...string) *LogContext {
	newData := lc.copyData()
	for _, tag := range tags {
		delete(newData.Tags, tag)
	}
	return &LogContext{data: newData}
}

// WithMetadata returns a new LogContext with additional metadata
func (lc *LogContext) WithMetadata(metadata map[string]string) *LogContext {
	newData := lc.copyData()
	for k, v := range metadata {
		newData.Metadata[k] = v
	}
	return &LogContext{data: newData}
}

// WithoutMetadata returns a new LogContext with specified metadata keys removed
func (lc *LogContext) WithoutMetadata(keys ...string) *LogContext {
	newData := lc.copyData()
	for _, key := range keys {
		delete(newData.Metadata, key)
	}
	return &LogContext{data: newData}
}

// copyData creates a deep copy of the LogContextData
func (lc *LogContext) copyData() LogContextData {
	tags := make(map[string]bool, len(lc.data.Tags))
	for k, v := range lc.data.Tags {
		tags[k] = v
	}
	metadata := make(map[string]string, len(lc.data.Metadata))
	for k, v := range lc.data.Metadata {
		metadata[k] = v
	}
	return LogContextData{
		Tags:      tags,
		Category:  lc.data.Category,
		Metadata:  metadata,
		SessionID: lc.data.SessionID,
	}
}

// GetLogContext retrieves the LogContext from the given context, or returns an empty one
func GetLogContext(ctx context.Context) *LogContext {
	if lc, ok := ctx.Value(logContextKey).(*LogContext); ok {
		return lc
	}
	return NewLogContext(LogContextData{})
}

// WithLogContext adds the LogContext to the given context
func WithLogContext(ctx context.Context, lc *LogContext) context.Context {
	return context.WithValue(ctx, logContextKey, lc)
}

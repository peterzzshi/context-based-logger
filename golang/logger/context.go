package logger

import "context"

type contextKey string

const logContextKey contextKey = "logContext"

type LogContext struct {
	data LogContextData
}

func NewLogContext(data LogContextData) *LogContext {
	if data.Tags == nil {
		data.Tags = make(map[string]bool)
	}
	if data.Metadata == nil {
		data.Metadata = make(map[string]string)
	}
	return &LogContext{data: data}
}

func (lc *LogContext) WithCategory(category string) *LogContext {
	newData := lc.copyData()
	newData.Category = category
	return &LogContext{data: newData}
}

func (lc *LogContext) WithSessionID(sessionID string) *LogContext {
	newData := lc.copyData()
	newData.SessionID = sessionID
	return &LogContext{data: newData}
}

func (lc *LogContext) WithTags(tags ...string) *LogContext {
	newData := lc.copyData()
	for _, tag := range tags {
		newData.Tags[tag] = true
	}
	return &LogContext{data: newData}
}

func (lc *LogContext) WithoutTags(tags ...string) *LogContext {
	newData := lc.copyData()
	for _, tag := range tags {
		delete(newData.Tags, tag)
	}
	return &LogContext{data: newData}
}

func (lc *LogContext) WithMetadata(metadata map[string]string) *LogContext {
	newData := lc.copyData()
	for k, v := range metadata {
		newData.Metadata[k] = v
	}
	return &LogContext{data: newData}
}

func (lc *LogContext) WithoutMetadata(keys ...string) *LogContext {
	newData := lc.copyData()
	for _, key := range keys {
		delete(newData.Metadata, key)
	}
	return &LogContext{data: newData}
}

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

func GetLogContext(ctx context.Context) *LogContext {
	if lc, ok := ctx.Value(logContextKey).(*LogContext); ok {
		return lc
	}
	return NewLogContext(LogContextData{})
}

// WithLogContext executes a callback with an enriched context containing the log context.
// Returns the result and error from the callback.
func WithLogContext[T any](ctx context.Context, logContext *LogContext, callback func(context.Context) (T, error)) (T, error) {
	enrichedCtx := context.WithValue(ctx, logContextKey, logContext)
	return callback(enrichedCtx)
}

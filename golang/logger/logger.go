package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

func Debug(ctx context.Context, args ...interface{}) {
	log(ctx, LevelDebug, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	log(ctx, LevelInfo, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	log(ctx, LevelWarn, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	log(ctx, LevelError, args...)
}

func log(ctx context.Context, level LogLevel, args ...interface{}) {
	logContext := GetLogContext(ctx)

	output := LogOutput{
		Level:   level,
		Details: make(map[string]interface{}),
	}

	if logContext.data.SessionID != "" {
		output.SessionID = logContext.data.SessionID
	}

	if len(logContext.data.Tags) > 0 {
		tags := make([]string, 0, len(logContext.data.Tags))
		for tag := range logContext.data.Tags {
			tags = append(tags, tag)
		}
		sort.Strings(tags)
		output.Details["tags"] = tags
	}

	if logContext.data.Category != "" {
		output.Details["category"] = logContext.data.Category
	}

	if len(logContext.data.Metadata) > 0 {
		metadata := make(map[string]string, len(logContext.data.Metadata))
		for k, v := range logContext.data.Metadata {
			metadata[k] = v
		}
		output.Details["metadata"] = metadata
	}

	if len(args) > 0 {
		message, stack := extractMessageAndStack(args...)
		if message != "" {
			output.Message = message
		}
		if stack != "" {
			output.Details["stack"] = stack
		}
	}

	output.Details["timestamp"] = time.Now().UTC().Format(time.RFC3339)

	jsonBytes, err := json.Marshal(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal log: %v\n", err)
		return
	}

	fmt.Fprintln(os.Stdout, string(jsonBytes))
}

func extractMessageAndStack(args ...interface{}) (message string, stack string) {
	if len(args) == 0 {
		return "", ""
	}

	// Single argument case
	if len(args) == 1 {
		if err, ok := args[0].(error); ok {
			message = err.Error()
			stack = fmt.Sprintf("%+v", err)
		} else {
			message = fmt.Sprint(args[0])
		}
		return message, stack
	}

	// Multiple arguments case
	firstArg := fmt.Sprint(args[0])
	lastArg := args[len(args)-1]

	// Check if last argument is an error
	if err, ok := lastArg.(error); ok {
		if len(args) == 2 {
			message = fmt.Sprintf("%s %s", firstArg, err.Error())
		} else {
			middle := fmt.Sprint(args[1 : len(args)-1]...)
			message = fmt.Sprintf("%s%s %s", firstArg, middle, err.Error())
		}
		stack = fmt.Sprintf("%+v", err)
	} else {
		message = fmt.Sprint(args...)
	}

	return message, stack
}

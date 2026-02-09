package logger

import (
	"context"
	"testing"
)

func TestLogContext_WithCategory(t *testing.T) {
	lc := NewLogContext(LogContextData{})
	lc2 := lc.WithCategory("test-category")

	if lc.data.Category != "" {
		t.Error("Original context should not be modified")
	}
	if lc2.data.Category != "test-category" {
		t.Errorf("Expected category 'test-category', got '%s'", lc2.data.Category)
	}
}

func TestLogContext_WithSessionID(t *testing.T) {
	lc := NewLogContext(LogContextData{})
	lc2 := lc.WithSessionID("session-123")

	if lc.data.SessionID != "" {
		t.Error("Original context should not be modified")
	}
	if lc2.data.SessionID != "session-123" {
		t.Errorf("Expected session ID 'session-123', got '%s'", lc2.data.SessionID)
	}
}

func TestLogContext_WithTags(t *testing.T) {
	lc := NewLogContext(LogContextData{})
	lc2 := lc.WithTags("tag1", "tag2")
	lc3 := lc2.WithTags("tag3")

	if len(lc.data.Tags) != 0 {
		t.Error("Original context should not be modified")
	}
	if len(lc2.data.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(lc2.data.Tags))
	}
	if len(lc3.data.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(lc3.data.Tags))
	}
}

func TestLogContext_WithoutTags(t *testing.T) {
	lc := NewLogContext(LogContextData{})
	lc2 := lc.WithTags("tag1", "tag2", "tag3")
	lc3 := lc2.WithoutTags("tag2")

	if len(lc3.data.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(lc3.data.Tags))
	}
	if lc3.data.Tags["tag2"] {
		t.Error("tag2 should have been removed")
	}
	if !lc3.data.Tags["tag1"] || !lc3.data.Tags["tag3"] {
		t.Error("tag1 and tag3 should still be present")
	}
}

func TestLogContext_WithMetadata(t *testing.T) {
	lc := NewLogContext(LogContextData{})
	lc2 := lc.WithMetadata(map[string]string{"key1": "value1"})
	lc3 := lc2.WithMetadata(map[string]string{"key2": "value2"})

	if len(lc.data.Metadata) != 0 {
		t.Error("Original context should not be modified")
	}
	if len(lc2.data.Metadata) != 1 {
		t.Errorf("Expected 1 metadata entry, got %d", len(lc2.data.Metadata))
	}
	if len(lc3.data.Metadata) != 2 {
		t.Errorf("Expected 2 metadata entries, got %d", len(lc3.data.Metadata))
	}
}

func TestLogContext_WithoutMetadata(t *testing.T) {
	lc := NewLogContext(LogContextData{})
	lc2 := lc.WithMetadata(map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	})
	lc3 := lc2.WithoutMetadata("key2")

	if len(lc3.data.Metadata) != 2 {
		t.Errorf("Expected 2 metadata entries, got %d", len(lc3.data.Metadata))
	}
	if _, ok := lc3.data.Metadata["key2"]; ok {
		t.Error("key2 should have been removed")
	}
	if lc3.data.Metadata["key1"] != "value1" || lc3.data.Metadata["key3"] != "value3" {
		t.Error("key1 and key3 should still be present")
	}
}

func TestGetLogContext(t *testing.T) {
	ctx := context.Background()
	lc := GetLogContext(ctx)

	if lc == nil {
		t.Fatal("Expected non-nil LogContext")
	}
	if len(lc.data.Tags) != 0 || len(lc.data.Metadata) != 0 {
		t.Error("Empty context should have no tags or metadata")
	}
}

func TestWithLogContext(t *testing.T) {
	ctx := context.Background()
	logCtx := NewLogContext(LogContextData{
		SessionID: "test-session",
		Category:  "test-category",
	})

	_, err := WithLogContext(ctx, logCtx, func(ctx context.Context) (struct{}, error) {
		retrieved := GetLogContext(ctx)

		if retrieved.data.SessionID != "test-session" {
			t.Errorf("Expected session ID 'test-session', got '%s'", retrieved.data.SessionID)
		}
		if retrieved.data.Category != "test-category" {
			t.Errorf("Expected category 'test-category', got '%s'", retrieved.data.Category)
		}
		return struct{}{}, nil
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestWithLogContext_CallbackString(t *testing.T) {
	ctx := context.Background()
	logCtx := NewLogContext(LogContextData{
		SessionID: "test-123",
		Category:  "test",
	}).WithTags("unit-test")

	result, err := WithLogContext(ctx, logCtx, func(ctx context.Context) (string, error) {
		lc := GetLogContext(ctx)
		if lc.data.SessionID != "test-123" {
			t.Error("Context not properly propagated")
		}
		return "success", nil
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != "success" {
		t.Errorf("Expected result 'success', got '%s'", result)
	}
}

func TestWithLogContext_CallbackInt(t *testing.T) {
	ctx := context.Background()
	logCtx := NewLogContext(LogContextData{SessionID: "calc-123"})

	result, err := WithLogContext(ctx, logCtx, func(ctx context.Context) (int, error) {
		lc := GetLogContext(ctx)
		if lc.data.SessionID != "calc-123" {
			t.Error("Context not properly propagated")
		}
		return 42, nil
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 42 {
		t.Errorf("Expected result 42, got %d", result)
	}
}

func TestLogger_ConvenienceFunctions(t *testing.T) {
	ctx := context.Background()
	logCtx := NewLogContext(LogContextData{
		SessionID: "req-123",
		Category:  "test",
	}).WithTags("api", "test").WithMetadata(map[string]string{
		"userId": "456",
	})

	_, _ = WithLogContext(ctx, logCtx, func(ctx context.Context) (struct{}, error) {
		Debug(ctx, "Debug message")
		Info(ctx, "Info message")
		Warn(ctx, "Warning message")
		Error(ctx, "Error message")
		return struct{}{}, nil
	})
}

func BenchmarkLogger_NoContext(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info(ctx, "Benchmark message")
	}
}

func BenchmarkLogger_WithContext(b *testing.B) {
	ctx := context.Background()
	logCtx := NewLogContext(LogContextData{
		SessionID: "req-123",
		Category:  "benchmark",
	}).WithTags("api", "test", "benchmark").WithMetadata(map[string]string{
		"userId":   "456",
		"endpoint": "/api/test",
		"method":   "GET",
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = WithLogContext(ctx, logCtx, func(ctx context.Context) (struct{}, error) {
			Info(ctx, "Benchmark message")
			return struct{}{}, nil
		})
	}
}

func BenchmarkLogContext_Copy(b *testing.B) {
	lc := NewLogContext(LogContextData{}).
		WithTags("tag1", "tag2", "tag3").
		WithMetadata(map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lc.WithCategory("test")
	}
}

func ExampleInfo() {
	ctx := context.Background()
	logCtx := NewLogContext(LogContextData{}).WithSessionID("req-123").
		WithCategory("api").
		WithTags("user-service").
		WithMetadata(map[string]string{
			"userId": "456",
		})

	_, _ = WithLogContext(ctx, logCtx, func(ctx context.Context) (struct{}, error) {
		Info(ctx, "Processing user request")
		return struct{}{}, nil
	})
}

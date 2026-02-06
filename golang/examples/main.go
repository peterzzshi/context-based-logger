package main

import (
	"context"
	"fmt"
	"time"

	"github.com/peterzzshi/context-based-logger/logger"
)

func main() {
	fmt.Println("=== Basic Logging ===")
	// Basic usage without context
	logger.Info("Application started")
	logger.Warn("This is a warning message")
	logger.Error("Something went wrong")

	fmt.Println("\n=== Context Logging ===")
	// Create a context with logging information
	ctx := context.Background()
	logCtx := logger.NewLogContext(logger.LogContextData{}).
		WithSessionID("req-123").
		WithTags("api", "user-service").
		WithCategory("http-request").
		WithMetadata(map[string]string{
			"userId":   "456",
			"endpoint": "/api/users",
		})

	ctx = logger.WithLogContext(ctx, logCtx)
	log := logger.New(ctx)

	log.Info("Processing user request")

	// Nested context with additional tags
	enrichedCtx := logger.WithLogContext(
		ctx,
		logCtx.WithTags("database").WithMetadata(map[string]string{
			"operation": "SELECT",
		}),
	)
	enrichedLog := logger.New(enrichedCtx)
	enrichedLog.Debug("Executing database query")

	log.Info("Request completed successfully")

	fmt.Println("\n=== API Request Simulation ===")
	// Simulate an API request handler
	handleUserRequest(ctx, "user-123")
}

// handleUserRequest simulates handling an API request with context
func handleUserRequest(parentCtx context.Context, userID string) {
	// Create a new context for this request
	logCtx := logger.NewLogContext(logger.LogContextData{}).
		WithSessionID(fmt.Sprintf("req-%d", time.Now().Unix())).
		WithCategory("api").
		WithTags("user-service", "database").
		WithMetadata(map[string]string{
			"userId":   userID,
			"endpoint": "/api/user",
			"method":   "GET",
		})

	ctx := logger.WithLogContext(parentCtx, logCtx)
	log := logger.New(ctx)

	log.Info(fmt.Sprintf("Fetching user data for user %s", userID))

	// Simulate validation step
	validatePermissions(ctx, userID)

	// Simulate database query
	time.Sleep(50 * time.Millisecond)

	log.Info("User data retrieved successfully")
}

// validatePermissions demonstrates context propagation to nested functions
func validatePermissions(ctx context.Context, userID string) {
	log := logger.New(ctx)
	log.Debug("Validating user permissions")

	// The log will automatically include all context from the parent
	time.Sleep(20 * time.Millisecond)

	log.Debug("Permissions validated")
}

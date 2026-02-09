package main

import (
	"context"
	"fmt"
	"time"

	"github.com/peterzzshi/context-based-logger/logger"
)

func main() {
	fmt.Println("=== Basic Logging ===")
	ctx := context.Background()
	logger.Info(ctx, "Application started")
	logger.Warn(ctx, "This is a warning message")
	logger.Error(ctx, "Something went wrong")

	fmt.Println("\n=== Context Logging with Callback ===")
	logCtx := logger.NewLogContext(logger.LogContextData{}).
		WithSessionID("req-123").
		WithTags("api", "user-service").
		WithCategory("http-request").
		WithMetadata(map[string]string{
			"userId":   "456",
			"endpoint": "/api/users",
		})

	_, err := logger.WithLogContext(ctx, logCtx, func(ctx context.Context) (struct{}, error) {
		logger.Info(ctx, "Processing user request")

		// Create enriched context for database operation
		enrichedLogCtx := logger.GetLogContext(ctx).
			WithTags("database").
			WithMetadata(map[string]string{"operation": "SELECT"})

		_, _ = logger.WithLogContext(ctx, enrichedLogCtx, func(ctx context.Context) (struct{}, error) {
			logger.Debug(ctx, "Executing database query")
			return struct{}{}, nil
		})

		logger.Info(ctx, "Request completed successfully")
		return struct{}{}, nil
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("\n=== API Request Simulation ===")
	handleUserRequest(ctx, "user-123")

	fmt.Println("\n=== Order Processing ===")
	orderResult := processOrder(ctx, "order-456")
	fmt.Printf("Order processing result: %s\n", orderResult)
}

func processOrder(parentCtx context.Context, orderID string) string {
	logCtx := logger.NewLogContext(logger.LogContextData{}).
		WithSessionID(fmt.Sprintf("order-%d", time.Now().Unix())).
		WithCategory("order-processing").
		WithTags("ecommerce", "payment").
		WithMetadata(map[string]string{
			"orderId": orderID,
			"region":  "us-west",
		})

	_, err := logger.WithLogContext(parentCtx, logCtx, func(ctx context.Context) (struct{}, error) {
		logger.Info(ctx, "Starting order processing")

		time.Sleep(30 * time.Millisecond)
		logger.Debug(ctx, "Validating payment method")

		time.Sleep(20 * time.Millisecond)
		logger.Info(ctx, "Order processed successfully")

		return struct{}{}, nil
	})

	if err != nil {
		return "FAILED"
	}
	return "SUCCESS"
}

func handleUserRequest(parentCtx context.Context, userID string) {
	logCtx := logger.NewLogContext(logger.LogContextData{}).
		WithSessionID(fmt.Sprintf("req-%d", time.Now().Unix())).
		WithCategory("api").
		WithTags("user-service", "database").
		WithMetadata(map[string]string{
			"userId":   userID,
			"endpoint": "/api/user",
			"method":   "GET",
		})

	_, _ = logger.WithLogContext(parentCtx, logCtx, func(ctx context.Context) (struct{}, error) {
		logger.Info(ctx, fmt.Sprintf("Fetching user data for user %s", userID))
		logger.Debug(ctx, "Validating user permissions")

		time.Sleep(50 * time.Millisecond)
		logger.Info(ctx, "User data retrieved successfully")

		return struct{}{}, nil
	})
}

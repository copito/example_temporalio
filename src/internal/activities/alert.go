package activities

import (
	"context"
	"log/slog"

	"go.temporal.io/sdk/activity"
)

// SendAlert sends an email notification.
func SendAlert(ctx context.Context, email string, value []float64) error {
	latestValue := value[len(value)-1]

	logger := activity.GetLogger(ctx)

	logger.Info("Sending alert - Condition failed", slog.String("email", email), slog.Float64("value", latestValue))
	// Implement email sending logic here
	return nil
}

func SendKafkaAlert(ctx context.Context, value []float64, isAlert bool) error {
	latestValue := value[len(value)-1]

	logger := activity.GetLogger(ctx)

	logger.Info("Sending alert to Kafka", slog.Bool("is_alert", isAlert), slog.Float64("value", latestValue))
	// Implement email sending logic here
	return nil
}

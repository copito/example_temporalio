package activities

import (
	"context"
	"errors"

	"github.com/copito/quality/src/internal/entities"
	"go.temporal.io/sdk/activity"
)

// EvaluateCondition checks if the condition is met.
func EvaluateCondition(ctx context.Context, value []float64, condition entities.Condition) (bool, error) {
	// Example condition: ">= 25"
	logger := activity.GetLogger(ctx)
	logger.Info("Starting evaluating expression")

	latestValue := value[len(value)-1]

	switch condition.Op {
	case entities.Equal:
		return latestValue == condition.Value, nil
	case entities.NotEqual:
		return latestValue != condition.Value, nil
	case entities.GreaterThan:
		return latestValue > condition.Value, nil
	case entities.GreaterThanOrEqual:
		return latestValue >= condition.Value, nil
	case entities.LessThan:
		return latestValue < condition.Value, nil
	case entities.LessThanOrEqual:
		return latestValue <= condition.Value, nil
	default:
		return false, errors.New("Invalid Operation")
	}
}

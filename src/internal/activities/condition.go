package activities

import (
	"context"
	"strconv"
	"strings"
)

// EvaluateCondition checks if the condition is met.
func EvaluateCondition(ctx context.Context, value []float64, condition string) (bool, error) {
	// Example condition: ">= 25"
	parts := strings.Fields(condition) // ["=>", "25"]
	if len(parts) != 2 {
		return false, nil
	}

	threshold, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return false, err
	}

	latestValue := value[len(value)-1]

	switch parts[0] {
	case ">=":
		return latestValue >= threshold, nil
	case "<=":
		return latestValue <= threshold, nil
	case ">":
		return latestValue > threshold, nil
	case "<":
		return latestValue < threshold, nil
	case "==":
		return latestValue == threshold, nil
	}

	return false, nil
}

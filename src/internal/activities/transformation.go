package activities

import (
	"context"
	"errors"
)

// ApplyTransformation calculates the transformation (e.g., DIFF).
func ApplyTransformation(ctx context.Context, data []float64, transformation string) ([]float64, error) {
	switch transformation {
	case "DIFF":
		return diffTransform(data)
	case "ROLLING_SUM":
		return rollingSumTransform(data)
	case "SUM":
		return sumTransform(data)
	default:
		return data, nil
	}
}

func diffTransform(data []float64) ([]float64, error) {
	if len(data) < 2 {
		return []float64{}, errors.New("NoEnoughDataPoints")
	}

	diffs := make([]float64, 0, len(data)-1)
	for i := 1; i < len(data); i++ {
		diffs = append(diffs, data[i]-data[i-1])
	}
	return diffs, nil
}

func rollingSumTransform(data []float64) ([]float64, error) {
	if len(data) == 0 {
		return nil, errors.New("NoDataPoints")
	}

	sums := make([]float64, len(data))
	sums[0] = data[0]
	for i := 1; i < len(data); i++ {
		sums[i] = sums[i-1] + data[i]
	}
	return sums, nil
}

func sumTransform(data []float64) ([]float64, error) {
	var total float64
	for _, v := range data {
		total += v
	}
	return []float64{total}, nil
}

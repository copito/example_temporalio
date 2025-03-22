package activities

import "context"

// FetchMetricData retrieves the last 300 data points for a given metric.
func FetchMetricData(ctx context.Context, metric string) ([]float64, error) {
	// Simulated: Call an external API or database to get metric data
	data := []float64{0, 50, 50, 200, 210, 200, 334, 420, 420, 420, 520} // Example data
	return data, nil
}

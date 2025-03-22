package workflows

import (
	"time"

	"github.com/copito/quality/src/internal/activities"

	"go.temporal.io/sdk/workflow"
)

// WorkflowInput defines the parameters for this workflow
type WorkflowInput struct {
	Schedule       string // e.g., "0 20 * * *"
	Metric         string // e.g., "row_count"
	Transformation string // e.g., "DIFF"
	Condition      string // e.g., ">= 25"
	AlertEmail     string
}

// MetricCheckWorkflow is the main workflow function
func MetricCheckWorkflow(ctx workflow.Context, input WorkflowInput) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("MyWorkflow")

	// Set activity options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		TaskQueue:           "metrics-task-queue",
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Step 1: Fetch metric data (last 300 points)
	var metricData []float64
	err := workflow.ExecuteActivity(ctx, activities.FetchMetricData, input.Metric).Get(ctx, &metricData)
	if err != nil {
		return err
	}

	// Step 2: Apply transformation
	var transformedValue []float64
	err = workflow.ExecuteActivity(ctx, activities.ApplyTransformation, metricData, input.Transformation).Get(ctx, &transformedValue)
	if err != nil {
		return err
	}

	// Example to wait
	err = workflow.Sleep(ctx, time.Minute)
	if err != nil {
		return err
	}

	// Step 3: Check condition
	var alertNeeded bool
	err = workflow.ExecuteActivity(ctx, activities.EvaluateCondition, transformedValue, input.Condition).Get(ctx, &alertNeeded)
	if err != nil {
		return err
	}

	// Step 4: Send alert if needed
	_ = workflow.ExecuteActivity(ctx, activities.SendKafkaAlert, transformedValue, alertNeeded)
	if alertNeeded {
		err = workflow.ExecuteActivity(ctx, activities.SendAlert, input.AlertEmail, transformedValue).Get(ctx, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

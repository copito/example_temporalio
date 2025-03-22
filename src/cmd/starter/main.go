package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"go.temporal.io/sdk/client"
	temp_logger "go.temporal.io/sdk/log"
)

// Copied over from "github.com/copito/quality/src/internal/workflows" (can be part of idl)
type WorkflowInput struct {
	Schedule       string // e.g., "0 20 * * *"
	Metric         string // e.g., "row_count"
	Transformation string // e.g., "DIFF"
	Condition      string // e.g., ">= 25"
	AlertEmail     string
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Started work request (starter/client)...")

	input := WorkflowInput{
		Schedule:       "*/5 * * * *",
		Metric:         "table.row_count",
		Transformation: "DIFF",
		Condition:      ">= 25",
		AlertEmail:     "user@example.com",
	}

	c, err := client.Dial(client.Options{
		Namespace: "data_quality",
		HostPort:  "localhost:7233",
		Identity:  "client_blah",
		// Credentials: internal.NewAPIKeyStaticCredentials(apiKey string),
		Logger: temp_logger.NewStructuredLogger(logger),
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:           fmt.Sprintf("metric-check-workflow-%v", "lemon-pie"),
		TaskQueue:    "metrics-task-queue",
		CronSchedule: input.Schedule, // User-defined cron expression
	}

	_, err = c.ExecuteWorkflow(
		context.Background(),
		options,
		"MetricCheckWorkflow", // github.com/copito/quality/src/internal/workflows.MetricCheckWorkflow
		input,
	)
	if err != nil {
		log.Fatalf("Unable to start workflow: %v", err)
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/copito/quality/src/internal/entities"
	"go.temporal.io/sdk/client"
	temp_logger "go.temporal.io/sdk/log"
)

// Copied over from "github.com/copito/quality/src/internal/workflows" (can be part of idl)
// type WorkflowInput struct {
// 	Schedule       string // e.g., "0 20 * * *"
// 	Metric         string // e.g., "row_count"
// 	Transformation string // e.g., "DIFF"
// 	Condition      string // e.g., ">= 25"
// 	AlertEmail     string
// }

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Started work request (starter/client)...")

	input := entities.WorkflowInput{
		Trigger: entities.Trigger{
			Type: entities.ScheduleTrigger,
			ScheduleTrigger: &entities.ScheduleTriggerType{
				Schedule: "*/5 * * * *",
				Timezone: "UTC",
			},
		},
		Metric:         "table.row_count",
		Transformation: "DIFF",
		Condition: entities.Condition{
			Op:    entities.GreaterThan,
			Value: 25,
		},
		AlertEmail: "user@example.com",
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

	if input.Trigger.Type == entities.ScheduleTrigger {
		options := client.StartWorkflowOptions{
			ID:           fmt.Sprintf("metric-check-workflow-%v", "jello"),
			TaskQueue:    "metrics-task-queue",
			CronSchedule: fmt.Sprintf("CRON_TZ=%v %v", input.Trigger.ScheduleTrigger.Timezone, input.Trigger.ScheduleTrigger.Schedule), // User-defined cron expression
		}
		r, err := c.ExecuteWorkflow(
			context.Background(),
			options,
			"MetricCheckWorkflow", // github.com/copito/quality/src/internal/workflows.MetricCheckWorkflow
			input,
		)
		if err != nil {
			log.Fatalf("Unable to start workflow: %v", err)
		}

		logger.Info("Workflow requested", slog.String("id", r.GetID()), slog.String("run_id", r.GetRunID()))
	} else {
		logger.Error("Not implemented trigger based...")
	}
}

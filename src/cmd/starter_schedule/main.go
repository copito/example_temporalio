package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/copito/quality/src/internal/entities"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	temp_logger "go.temporal.io/sdk/log"
	// "go.temporal.io/sdk/schedulepb"
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

		// params := []entities.WorkflowInput{input}
		// scheduleCronExpression := fmt.Sprintf("CRON_TZ=%v %v", input.Trigger.ScheduleTrigger.Timezone, input.Trigger.ScheduleTrigger.Schedule)

		sc := c.ScheduleClient()

		options := client.ScheduleOptions{
			ID: fmt.Sprintf("metric-check-workflow-%v-scheduled", "jello"),
			Spec: client.ScheduleSpec{
				CronExpressions: strings.Split(input.Trigger.ScheduleTrigger.Schedule, " "),
			},
			Action: &client.ScheduleWorkflowAction{
				ID:        fmt.Sprintf("metric-check-workflow-%v", "jello-scheduled"),
				Workflow:  "MetricCheckWorkflow",
				TaskQueue: "metrics-task-queue",
				Args:      []interface{},
			},
			Paused:         false,
			PauseOnFailure: false,
			CatchupWindow:  4 * time.Hour, // if the server is down it will try to catch up if at least 4 hours aways
			Overlap:        enums.SCHEDULE_OVERLAP_POLICY_BUFFER_ONE,
		}

		scheduleHandler, err := sc.Create(context.Background(), options)
		if err != nil {
			logger.Error("Scheduled Workflow failed to be created...", slog.Any("err", err))
			return
		}

		logger.Info("Workflow requested", slog.String("id", scheduleHandler.GetID()))
	} else {
		logger.Error("Not implemented trigger based...")
	}
}

package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/copito/quality/src/internal/activities"
	"github.com/copito/quality/src/internal/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"

	temp_logger "go.temporal.io/sdk/log"
)

const (
	NAMESPACE string = "data_quality"
	QUEUE     string = "metrics-task-queue"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Started worker agent...")

	sigs := make(chan os.Signal, 1)
	doneChan := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Temporal client
	c, err := client.Dial(client.Options{
		Namespace: NAMESPACE,
		HostPort:  client.DefaultHostPort, // "localhost:7233",
		// Credentials: internal.NewAPIKeyStaticCredentials(apiKey string),
		Logger: temp_logger.NewStructuredLogger(logger),
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Create worker
	w := worker.New(c, QUEUE, worker.Options{})

	// Register workflow & activities
	w.RegisterWorkflowWithOptions(workflows.MetricCheckWorkflow, workflow.RegisterOptions{
		Name: "MetricCheckWorkflow", // This is so the starter/client can call it without being in same repo
	})
	w.RegisterActivity(activities.FetchMetricData)
	w.RegisterActivity(activities.ApplyTransformation)
	w.RegisterActivity(activities.EvaluateCondition)
	w.RegisterActivity(activities.SendAlert)

	// Start worker
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("Unable to start worker: %v", err)
	}

	// Wait for termination signal
	go func(logger *slog.Logger) {
		sig := <-sigs
		logger.Info("Closing down based on signal...", slog.Any("signal", sig))
		doneChan <- true
	}(logger)

	<-doneChan
}

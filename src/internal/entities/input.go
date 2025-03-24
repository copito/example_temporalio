package entities

type TriggerType string

type Operator string

const (
	Equal              Operator = "EQUAL"
	NotEqual           Operator = "NOT_EQUAL"
	GreaterThan        Operator = "GREATER_THAN"
	GreaterThanOrEqual Operator = "GREATER_THAN_OR_EQUAL"
	LessThan           Operator = "LESS_THAN"
	LessThanOrEqual    Operator = "LESS_THAN_OR_EQUAL"
)

const (
	ScheduleTrigger TriggerType = "SCHEDULE"
	EventTrigger    TriggerType = "EVENT"
)

type ScheduleTriggerType struct {
	Schedule string // e.g., "0 20 * * *"
	Timezone string
}

type EventTriggerType struct {
	Source string // e.g., kafka
	Topic  string // e.g., "db.debezium.lemon_pie"
}

type Trigger struct {
	Type            TriggerType          `json:"type"`
	ScheduleTrigger *ScheduleTriggerType `json:"schedule_type,omitempty"`
	EventTrigger    *EventTriggerType    `json:"event_type,omitempty"`
}

type Condition struct {
	Op    Operator `json:"op"`
	Value float64  `json:"value"`
}

// Copied over from "github.com/copito/quality/src/internal/workflows" (can be part of idl)
type WorkflowInput struct {
	Trigger        Trigger   // e.g., "0 20 * * *"
	Metric         string    // e.g., "table.row_count"
	Transformation string    // e.g., "DIFF"
	Condition      Condition // e.g., ">= 25"
	AlertEmail     string
}

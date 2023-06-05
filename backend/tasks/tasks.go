package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const (
	RunPipelineTask = "pipeline:run"
)

type RunPipelinePayload struct {
	RunID     uint
	StepIndex uint
}

func NewRunPipelineTask(runID uint, stepIndex uint) (*asynq.Task, error) {
	payload, err := json.Marshal(RunPipelinePayload{RunID: runID, StepIndex: stepIndex})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(RunPipelineTask, payload), nil
}

func HandleRunPipelineTask(ctx context.Context, t *asynq.Task) error {
	var p RunPipelinePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// Image resizing code ...
	return nil
}

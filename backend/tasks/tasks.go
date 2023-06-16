package tasks

import (
	"context"
	"di/steps"
	"encoding/json"
	"fmt"

	"github.com/dominikbraun/graph"
	"github.com/hibiken/asynq"
)

const (
	RunPipelineTask = "pipeline:run"
)

type RunPipelinePayload struct {
	Graph     graph.Graph[int, steps.Step]
	StepIndex uint
}

func NewRunPipelineTask(graph graph.Graph[int, steps.Step], stepIndex uint) (*asynq.Task, error) {
	payload, err := json.Marshal(RunPipelinePayload{Graph: graph, StepIndex: stepIndex})
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

	return nil
}

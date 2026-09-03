package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/inbound"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

type MatchingWorker struct {
	processMatching inbound.ProcessMatchingInt
}

func NewMatchingWorker(processMatching inbound.ProcessMatchingInt) *MatchingWorker {
	return &MatchingWorker{
		processMatching: processMatching,
	}
}

func (w *MatchingWorker) HandleStartMatching(ctx context.Context, task *asynq.Task) error {
	var payload StartMatchingPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		logger.Error("asynq matching worker failed to decode task payload", zap.Error(err))
		return fmt.Errorf("decode start matching payload: %w", err)
	}

	if payload.RunID == "" {
		logger.Error("asynq matching worker received task with missing run ID")
		return fmt.Errorf("missing matching run ID")
	}

	logger.Info(
		"asynq matching worker picked up task",
		zap.String("task_type", task.Type()),
		zap.String("run_id", payload.RunID),
	)

	if err := w.processMatching.Execute(ctx, payload.RunID); err != nil {
		logger.Error(
			"asynq matching worker failed to execute matching",
			zap.Error(err),
			zap.String("run_id", payload.RunID),
		)
		return fmt.Errorf("execute matching: %w", err)
	}

	logger.Info(
		"asynq matching worker finished task successfully",
		zap.String("run_id", payload.RunID),
	)
	return nil
}

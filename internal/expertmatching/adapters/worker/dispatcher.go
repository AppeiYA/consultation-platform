package worker

import (
	"context"
	"encoding/json"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
	"github.com/hibiken/asynq"
)

type MatchingJobDispatcher struct {
	client *asynq.Client
}

func NewMatchingJobDispatcher(client *asynq.Client) *MatchingJobDispatcher {
	return &MatchingJobDispatcher{
		client: client,
	}
}

func (d *MatchingJobDispatcher) DispatchMatching(ctx context.Context, runID string) error {
	payload, err := json.Marshal(StartMatchingPayload{
		RunID: runID,
	})
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeStartMatching, payload)
	_, err = d.client.EnqueueContext(ctx, task)
	return err
}

var _ outbound.MatchingJobDispatcher = (*MatchingJobDispatcher)(nil)

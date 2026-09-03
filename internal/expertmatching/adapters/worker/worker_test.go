package worker

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"
)

type mockProcessMatching struct {
	ExecuteFn func(ctx context.Context, runID string) error
}

func (m *mockProcessMatching) Execute(ctx context.Context, runID string) error {
	if m.ExecuteFn != nil {
		return m.ExecuteFn(ctx, runID)
	}
	return nil
}

func TestMatchingWorker_HandleStartMatching(t *testing.T) {
	ctx := context.Background()

	t.Run("successfully processes start matching task", func(t *testing.T) {
		var executedRunID string
		mockProc := &mockProcessMatching{
			ExecuteFn: func(ctx context.Context, runID string) error {
				executedRunID = runID
				return nil
			},
		}

		worker := NewMatchingWorker(mockProc)

		payload, err := json.Marshal(StartMatchingPayload{RunID: "mrun_12345"})
		require.NoError(t, err)

		task := asynq.NewTask(TypeStartMatching, payload)
		err = worker.HandleStartMatching(ctx, task)

		require.NoError(t, err)
		require.Equal(t, "mrun_12345", executedRunID)
	})

	t.Run("returns error on invalid JSON payload", func(t *testing.T) {
		mockProc := &mockProcessMatching{}
		worker := NewMatchingWorker(mockProc)

		task := asynq.NewTask(TypeStartMatching, []byte("invalid-json"))
		err := worker.HandleStartMatching(ctx, task)

		require.Error(t, err)
		require.Contains(t, err.Error(), "decode start matching payload")
	})

	t.Run("returns error on missing run_id in payload", func(t *testing.T) {
		mockProc := &mockProcessMatching{}
		worker := NewMatchingWorker(mockProc)

		payload, _ := json.Marshal(StartMatchingPayload{RunID: ""})
		task := asynq.NewTask(TypeStartMatching, payload)
		err := worker.HandleStartMatching(ctx, task)

		require.Error(t, err)
		require.Contains(t, err.Error(), "missing matching run ID")
	})

	t.Run("returns error when process matching use case fails", func(t *testing.T) {
		mockProc := &mockProcessMatching{
			ExecuteFn: func(ctx context.Context, runID string) error {
				return errors.New("candidate generation failed")
			},
		}
		worker := NewMatchingWorker(mockProc)

		payload, _ := json.Marshal(StartMatchingPayload{RunID: "mrun_fail"})
		task := asynq.NewTask(TypeStartMatching, payload)
		err := worker.HandleStartMatching(ctx, task)

		require.Error(t, err)
		require.Contains(t, err.Error(), "execute matching: candidate generation failed")
	})
}

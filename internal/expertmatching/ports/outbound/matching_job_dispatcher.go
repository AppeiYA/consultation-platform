package outbound

import "context"

type MatchingJobDispatcher interface {
	DispatchMatching(ctx context.Context, runID string) error
}

package outbound

import "context"

type MatchingJob struct {
	RunID  string
	CaseID string
}

type MatchingJobEnqueuer interface {
	Enqueue(ctx context.Context, job MatchingJob) error
}

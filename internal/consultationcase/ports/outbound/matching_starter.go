package outbound

import "context"

type MatchingStarter interface {
	StartMatching(ctx context.Context, caseID string) error
}

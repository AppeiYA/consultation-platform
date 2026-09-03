package inbound

import "context"

type ProcessMatchingInt interface {
	Execute(ctx context.Context, runID string) error
}

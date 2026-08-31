package inbound

import "context"

type RemoveExpertiseInt interface {
	Execute(ctx context.Context, userID string, expertiseID string) error
}

package inbound

import "context"

type DeactivateAvailabilityInt interface {
	Execute(ctx context.Context, userID string, availabilityID string) error
}
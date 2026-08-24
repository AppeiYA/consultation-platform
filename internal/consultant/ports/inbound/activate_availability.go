package inbound

import "context"

type ActivateAvailabilityInt interface {
	Execute(ctx context.Context, userID string, availabilityID string) error
}
package outbound

import "context"

type EventPublisher interface {
    Publish(ctx context.Context, events ...any) error
}
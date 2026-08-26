package inbound

import "context"

type DeleteCaseInt interface {
	Execute(ctx context.Context, clientID string, caseID string) error
}
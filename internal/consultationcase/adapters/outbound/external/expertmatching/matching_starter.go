package consultationCaseExpertMatchingAdapter

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching"
)

type MatchingStarterAdapter struct {
	expertMatchingModule *expertmatching.Module
}

func NewMatchingStarterAdapter(
	expertMatchingModule *expertmatching.Module,
) *MatchingStarterAdapter {
	return &MatchingStarterAdapter{
		expertMatchingModule: expertMatchingModule,
	}
}

func (a *MatchingStarterAdapter) StartMatching(ctx context.Context, caseID string) error {
	if a.expertMatchingModule == nil || a.expertMatchingModule.StartMatching == nil {
		return nil
	}
	_, err := a.expertMatchingModule.StartMatching.Execute(ctx, caseID)
	return err
}

var _ outbound.MatchingStarter = (*MatchingStarterAdapter)(nil)

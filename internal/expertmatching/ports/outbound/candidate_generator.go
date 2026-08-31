package outbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
)

type CandidateGenerator interface {
	GenerateCandidates(ctx context.Context, criteria domain.CandidateGenerationCriteria) (domain.CandidatePool, error)
}

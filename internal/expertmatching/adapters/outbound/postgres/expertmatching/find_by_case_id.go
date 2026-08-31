package expertMatchingRepo

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
)

var findByCaseIDQuery = `
	SELECT 
		id,
		case_id,
		status,
		ranking_version,
		failure_reason,
		cancellation_reason,
		created_at,
		completed_at
	FROM matching_runs
	WHERE case_id = $1
	ORDER BY created_at DESC
`

func (r *ExpertMatchingRepository) FindByCaseID(ctx context.Context, caseID string) ([]*domain.MatchingRun, error) {
	var runModels []MatchingRunModel
	err := r.repository.Executor(ctx).SelectContext(ctx, &runModels, findByCaseIDQuery, caseID)
	if err != nil {
		return nil, err
	}

	results := make([]*domain.MatchingRun, 0, len(runModels))
	for _, m := range runModels {
		candidates, err := r.findCandidatesByRunID(ctx, m.ID)
		if err != nil {
			return nil, err
		}
		domainRun, err := m.ToDomain(candidates)
		if err != nil {
			return nil, err
		}
		results = append(results, domainRun)
	}

	return results, nil
}
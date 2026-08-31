package expertMatchingRepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
)

var findLatestByCaseIDQuery = `
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
	LIMIT 1
`
func (r *ExpertMatchingRepository) FindLatestByCaseID(ctx context.Context, caseID string) (*domain.MatchingRun, error) {
	

	var runModel MatchingRunModel
	err := r.repository.Executor(ctx).GetContext(ctx, &runModel, findLatestByCaseIDQuery, caseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrMatchingRunNotFound
		}
		return nil, err
	}

	candidates, err := r.findCandidatesByRunID(ctx, runModel.ID)
	if err != nil {
		return nil, err
	}

	return runModel.ToDomain(candidates)
}
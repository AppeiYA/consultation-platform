package expertMatchingRepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
)

var findByIDquery = `
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
	WHERE id = $1
`
var findByRunIDQuery = `
	SELECT 
		id,
		run_id,
		consultant_id,
		rank_position,
		score,
		reasons::text AS reasons
	FROM matching_run_candidates
	WHERE run_id = $1
	ORDER BY rank_position ASC
`
func (r *ExpertMatchingRepository) FindByID(ctx context.Context, id string) (*domain.MatchingRun, error) {
	var runModel MatchingRunModel
	err := r.repository.Executor(ctx).GetContext(ctx, &runModel, findByIDquery, id)
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

func (r *ExpertMatchingRepository) findCandidatesByRunID(ctx context.Context, runID string) ([]CandidateModel, error) {
	var candidates []CandidateModel
	err := r.repository.Executor(ctx).SelectContext(ctx, &candidates, findByRunIDQuery, runID)
	if err != nil {
		return nil, err
	}

	return candidates, nil
}
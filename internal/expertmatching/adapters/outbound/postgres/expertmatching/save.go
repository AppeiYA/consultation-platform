package expertMatchingRepo

import (
	"context"
	"encoding/json"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
)

var saveRunQuery = `
		INSERT INTO matching_runs (
			id,
			case_id,
			status,
			ranking_version,
			failure_reason,
			cancellation_reason,
			created_at,
			completed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			failure_reason = EXCLUDED.failure_reason,
			cancellation_reason = EXCLUDED.cancellation_reason,
			completed_at = EXCLUDED.completed_at
`
var deleteMatchingRunCandidatesQuery = `DELETE FROM matching_run_candidates WHERE run_id = $1
`
var insertMatchingRunCandidateQuery = `
	INSERT INTO matching_run_candidates (
		run_id,
		consultant_id,
		rank_position,
		score,
		reasons
	) VALUES ($1, $2, $3, $4, $5)
`
func (r *ExpertMatchingRepository) Save(ctx context.Context, run *domain.MatchingRun) error {
	var failureReason *string
	if len(run.FailureReason()) > 0 {
		f := run.FailureReason()
		failureReason = &f
	}

	var cancellationReason *string
	if len(run.CancellationReason()) > 0 {
		c := run.CancellationReason()
		cancellationReason = &c
	}

	executor := r.repository.Executor(ctx)
	_, err := executor.ExecContext(
		ctx,
		saveRunQuery,
		run.ID(),
		run.CaseID(),
		string(run.Status()),
		run.RankingVersion().Value(),
		failureReason,
		cancellationReason,
		run.CreatedAt(),
		run.CompletedAt(),
	)
	if err != nil {
		return err
	}

	// Delete and re-insert candidates for this run to keep it idempotent
	if len(run.Candidates()) > 0 {
		if _, err := executor.ExecContext(ctx, deleteMatchingRunCandidatesQuery, run.ID()); err != nil {
			return err
		}

		for _, cand := range run.Candidates() {
			reasons := cand.Reasons()
			reasonJSONs := make([]ReasonJSON, 0, len(reasons))
			for _, re := range reasons {
				reasonJSONs = append(reasonJSONs, ReasonJSON{
					Factor: re.Factor().String(),
					Detail: re.Detail(),
				})
			}
			reasonsBytes, err := json.Marshal(reasonJSONs)
			if err != nil {
				return err
			}

			_, err = executor.ExecContext(
				ctx,
				insertMatchingRunCandidateQuery,
				run.ID(),
				cand.ConsultantID(),
				cand.Rank().Value(),
				cand.Score().Value(),
				string(reasonsBytes),
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}


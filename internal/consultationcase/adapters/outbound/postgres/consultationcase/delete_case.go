package consultationCaseRepo

import "context"

var deleteCaseQuery = `
	DELETE FROM consultation_cases WHERE id = $1
`
func (r *ConsultationCaseRepository) DeleteCase(ctx context.Context, caseID string) error {
	executor := r.repository.Executor(ctx)
	_, err := executor.ExecContext(ctx, deleteCaseQuery, caseID)
	if err != nil {
		return err
	}
	return nil
}
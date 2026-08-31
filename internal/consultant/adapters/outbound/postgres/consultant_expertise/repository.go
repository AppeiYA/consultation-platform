package consultant_expertise

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/shared/db"
)

type ConsultantExpertiseRepository struct {
	repository *db.Repository
}

func NewConsultantExpertiseRepository(repository *db.Repository) *ConsultantExpertiseRepository {
	return &ConsultantExpertiseRepository{
		repository: repository,
	}
}

func (r *ConsultantExpertiseRepository) SaveMany(ctx context.Context, expertises []*domain.Expertise) error {
	if len(expertises) == 0 {
		return nil
	}

	query := `
		INSERT INTO consultant_expertises (id, consultant_id, name, created_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		ON CONFLICT (consultant_id, name) DO NOTHING
	`

	executor := r.repository.Executor(ctx)
	for _, exp := range expertises {
		if _, err := executor.ExecContext(ctx, query, exp.ID(), exp.ConsultantID(), exp.Name()); err != nil {
			return err
		}
	}
	return nil
}

func (r *ConsultantExpertiseRepository) Add(ctx context.Context, expertise *domain.Expertise) error {
	query := `
		INSERT INTO consultant_expertises (id, consultant_id, name, created_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		ON CONFLICT (consultant_id, name) DO NOTHING
	`
	_, err := r.repository.Executor(ctx).ExecContext(ctx, query, expertise.ID(), expertise.ConsultantID(), expertise.Name())
	return err
}

func (r *ConsultantExpertiseRepository) FindByConsultantID(ctx context.Context, consultantID string) ([]*domain.Expertise, error) {
	query := `
		SELECT id, consultant_id, name, created_at
		FROM consultant_expertises
		WHERE consultant_id = $1
		ORDER BY name ASC
	`

	var models []ExpertiseModel
	err := r.repository.Executor(ctx).SelectContext(ctx, &models, query, consultantID)
	if err != nil {
		return nil, err
	}

	results := make([]*domain.Expertise, 0, len(models))
	for _, m := range models {
		domainExp, err := m.ToDomain()
		if err != nil {
			return nil, err
		}
		results = append(results, domainExp)
	}

	return results, nil
}

func (r *ConsultantExpertiseRepository) Delete(ctx context.Context, consultantID string, expertiseID string) error {
	query := `
		DELETE FROM consultant_expertises
		WHERE id = $1 AND consultant_id = $2
	`
	_, err := r.repository.Executor(ctx).ExecContext(ctx, query, expertiseID, consultantID)
	return err
}

func (r *ConsultantExpertiseRepository) ReplaceAll(ctx context.Context, consultantID string, expertises []*domain.Expertise) error {
	executor := r.repository.Executor(ctx)

	deleteQuery := `DELETE FROM consultant_expertises WHERE consultant_id = $1`
	if _, err := executor.ExecContext(ctx, deleteQuery, consultantID); err != nil {
		return err
	}

	if len(expertises) == 0 {
		return nil
	}

	insertQuery := `
		INSERT INTO consultant_expertises (id, consultant_id, name, created_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		ON CONFLICT (consultant_id, name) DO NOTHING
	`

	for _, exp := range expertises {
		if _, err := executor.ExecContext(ctx, insertQuery, exp.ID(), consultantID, exp.Name()); err != nil {
			return err
		}
	}

	return nil
}

var _ outbound.ExpertiseRepository = (*ConsultantExpertiseRepository)(nil)

package consultant_adapter

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/shared/db"
	"github.com/lib/pq"
)

type CandidateGeneratorAdapter struct {
	repository *db.Repository
}

func NewCandidateGeneratorAdapter(repository *db.Repository) *CandidateGeneratorAdapter {
	return &CandidateGeneratorAdapter{
		repository: repository,
	}
}

type candidateRow struct {
	ConsultantID       string         `db:"id"`
	Profession         string         `db:"profession"`
	YearsExperience    int            `db:"years_experience"`
	Bio                string         `db:"bio"`
	IsAcceptingClients bool           `db:"is_accepting_clients"`
	IsVerified         bool           `db:"is_verified"`
	Expertises         pq.StringArray `db:"expertises"`
}

func (a *CandidateGeneratorAdapter) GenerateCandidates(
	ctx context.Context,
	criteria domain.CandidateGenerationCriteria,
) (domain.CandidatePool, error) {
	query := `
		SELECT 
			c.id,
			p.name AS profession,
			c.years_experience,
			c.bio,
			c.is_accepting_clients,
			COALESCE(bool_or(v.status = 'APPROVED'), false) AS is_verified,
			COALESCE(ARRAY_AGG(DISTINCT ce.name) FILTER (WHERE ce.name IS NOT NULL), '{}') AS expertises
		FROM consultants c
		JOIN professions p ON c.profession_id = p.id
		LEFT JOIN consultant_verifications v ON c.id = v.consultant_id
		LEFT JOIN consultant_expertises ce ON c.id = ce.consultant_id
		WHERE c.is_accepting_clients = true
		  AND (
		      $1 = '' 
		      OR p.id = $1
		      OR UPPER(p.name) = UPPER($1)
		      OR REPLACE(UPPER(p.name), '_', ' ') = REPLACE(UPPER($1), '_', ' ')
		      OR UPPER(p.name) = UPPER(REPLACE($1, ' ', '_'))
		      OR p.name ILIKE '%' || $1 || '%'
		      OR EXISTS (
		          SELECT 1 FROM consultant_expertises ce2
		          WHERE ce2.consultant_id = c.id
		            AND (
		                UPPER(ce2.name) = UPPER($1)
		                OR ce2.name ILIKE '%' || $1 || '%'
		                OR $1 ILIKE '%' || ce2.name || '%'
		            )
		      )
		  )
		GROUP BY c.id, p.name, c.years_experience, c.bio, c.is_accepting_clients
	`

	var rows []candidateRow
	err := a.repository.Executor(ctx).SelectContext(ctx, &rows, query, criteria.Category().Value())
	if err != nil {
		return domain.CandidatePool{}, err
	}

	profiles := make([]domain.CandidateProfile, 0, len(rows))
	for _, r := range rows {
		if criteria.RequiredVerifiedOnly() && !r.IsVerified {
			continue
		}
		profCategory, err := domain.NewMatchingCategory(r.Profession)
		if err != nil {
			continue
		}

		expertiseList := make([]domain.Expertise, 0, len(r.Expertises))
		for _, expName := range r.Expertises {
			if expObj, err := domain.NewExpertise(expName); err == nil {
				expertiseList = append(expertiseList, expObj)
			}
		}

		profile, err := domain.NewCandidateProfile(
			r.ConsultantID,
			profCategory,
			r.Profession,
			expertiseList,
			r.YearsExperience,
			r.Bio,
		)
		if err != nil {
			continue
		}
		profiles = append(profiles, profile)
	}

	return domain.NewCandidatePool(profiles, domain.DefaultMaxCandidatePoolSize)
}

var _ outbound.CandidateGenerator = (*CandidateGeneratorAdapter)(nil)

package domain

type CandidateGenerationCriteria struct {
	category             MatchingCategory
	requiredVerifiedOnly bool
}

func NewCandidateGenerationCriteria(
	category MatchingCategory,
	requiredVerifiedOnly bool,
) CandidateGenerationCriteria {
	return CandidateGenerationCriteria{
		category:             category,
		requiredVerifiedOnly: requiredVerifiedOnly,
	}
}

func (c CandidateGenerationCriteria) Category() MatchingCategory {
	return c.category
}

func (c CandidateGenerationCriteria) RequiredVerifiedOnly() bool {
	return c.requiredVerifiedOnly
}

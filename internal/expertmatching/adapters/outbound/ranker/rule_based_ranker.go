package ranker

import (
	"context"
	"math"
	"sort"
	"strings"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
)

type RuleBasedRanker struct{}

func NewRuleBasedRanker() *RuleBasedRanker {
	return &RuleBasedRanker{}
}

type scoredCandidate struct {
	profile domain.CandidateProfile
	score   float64
	reasons []domain.MatchReason
}

func (r *RuleBasedRanker) Rank(
	ctx context.Context,
	req outbound.RankingRequest,
) ([]domain.RankedCandidate, error) {
	candidates := req.CandidatePool.Candidates()
	if len(candidates) == 0 {
		return []domain.RankedCandidate{}, nil
	}

	targetCategory := req.CaseDetails.Category.Value()

	scored := make([]scoredCandidate, 0, len(candidates))
	for _, c := range candidates {
		var totalScore float64
		var reasons []domain.MatchReason

		// 1. Profession / Category match (0.50 max)
		if strings.EqualFold(c.Profession(), targetCategory) || c.Category().Equals(req.CaseDetails.Category) {
			totalScore += 0.50
			reason, _ := domain.NewMatchReason(domain.MatchFactorCategory, "Matches requested consultation category")
			reasons = append(reasons, reason)
		} else if strings.Contains(strings.ToLower(c.Profession()), strings.ToLower(targetCategory)) {
			totalScore += 0.35
			reason, _ := domain.NewMatchReason(domain.MatchFactorCategory, "Partially matches requested category")
			reasons = append(reasons, reason)
		}

		// 2. Experience score (0.50 max; scaled up to 15 years)
		expRatio := math.Min(float64(c.YearsExperience())/15.0, 1.0)
		expScore := 0.50 * expRatio
		totalScore += expScore
		if c.YearsExperience() >= 5 {
			reason, _ := domain.NewMatchReason(domain.MatchFactorExperience, "High years of domain experience")
			reasons = append(reasons, reason)
		}

		// Clamp between 0.0 and 1.0
		if totalScore > 1.0 {
			totalScore = 1.0
		}
		if totalScore < 0.0 {
			totalScore = 0.0
		}

		scored = append(scored, scoredCandidate{
			profile: c,
			score:   totalScore,
			reasons: reasons,
		})
	}

	// Sort descending by score; tie-break by ConsultantID for determinism
	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].score != scored[j].score {
			return scored[i].score > scored[j].score
		}
		return scored[i].profile.ConsultantID() < scored[j].profile.ConsultantID()
	})

	results := make([]domain.RankedCandidate, 0, len(scored))
	for idx, sc := range scored {
		rankVal, err := domain.NewRank(idx + 1)
		if err != nil {
			return nil, err
		}
		scoreVal, err := domain.NewMatchScore(sc.score)
		if err != nil {
			return nil, err
		}

		cand, err := domain.NewRankedCandidate(sc.profile.ConsultantID(), rankVal, scoreVal, sc.reasons)
		if err != nil {
			return nil, err
		}
		results = append(results, cand)
	}

	return results, nil
}

var _ outbound.CandidateRanker = (*RuleBasedRanker)(nil)

package geminiranker

import (
	"fmt"
	"sort"
	"strings"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
)

func mapRankingResponse(
	req outbound.RankingRequest,
	response rankingResponse,
) ([]domain.RankedCandidate, error) {
	candidateIDs := make(map[string]struct{})

	for _, candidate := range req.CandidatePool.Candidates() {
		candidateIDs[candidate.ConsultantID()] = struct{}{}
	}

	if len(response.Candidates) != len(candidateIDs) {
		return nil, fmt.Errorf(
			"gemini returned %d candidates, expected %d",
			len(response.Candidates),
			len(candidateIDs),
		)
	}

	seen := make(map[string]struct{})

	type scoredCandidate struct {
		consultantID string
		score        domain.MatchScore
		reasons      []domain.MatchReason
	}

	scored := make([]scoredCandidate, 0, len(response.Candidates))

	for _, item := range response.Candidates {
		if _, exists := candidateIDs[item.ConsultantID]; !exists {
			return nil, fmt.Errorf(
				"gemini returned unknown consultant %q",
				item.ConsultantID,
			)
		}

		if _, exists := seen[item.ConsultantID]; exists {
			return nil, fmt.Errorf(
				"gemini returned duplicate consultant %q",
				item.ConsultantID,
			)
		}

		seen[item.ConsultantID] = struct{}{}

		score, err := domain.NewMatchScore(item.Score)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid score for consultant %q: %w",
				item.ConsultantID,
				err,
			)
		}

		reasons, err := mapReasons(item.Reasons)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid reasons for consultant %q: %w",
				item.ConsultantID,
				err,
			)
		}

		scored = append(scored, scoredCandidate{
			consultantID: item.ConsultantID,
			score:        score,
			reasons:      reasons,
		})
	}

	sort.SliceStable(scored, func(i, j int) bool {
		return scored[i].score.Value() > scored[j].score.Value()
	})

	result := make([]domain.RankedCandidate, 0, len(scored))

	for index, candidate := range scored {
		rank, err := domain.NewRank(index + 1)
		if err != nil {
			return nil, fmt.Errorf("create candidate rank: %w", err)
		}

		rankedCandidate, err := domain.NewRankedCandidate(
			candidate.consultantID,
			rank,
			candidate.score,
			candidate.reasons,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"create ranked candidate %q: %w",
				candidate.consultantID,
				err,
			)
		}

		result = append(result, rankedCandidate)
	}

	return result, nil
}

func mapReasons(
	input []reasonResponse,
) ([]domain.MatchReason, error) {

	result := make([]domain.MatchReason, 0, len(input))

	for _, reason := range input {
		factor, err := mapMatchFactor(reason.Factor)
		if err != nil {
			return nil, err
		}

		matchReason, err := domain.NewMatchReason(
			factor,
			reason.Detail,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, matchReason)
	}

	return result, nil
}


func mapMatchFactor(value string) (domain.MatchFactor, error) {
	switch strings.ToUpper(strings.TrimSpace(value)) {

	case "EXPERTISE":
		return domain.MatchFactorExpertise, nil

	case "EXPERIENCE":
		return domain.MatchFactorExperience, nil

	case "CATEGORY":
		return domain.MatchFactorCategory, nil

	case "LANGUAGE":
		return domain.MatchFactorLanguage, nil

	default:
		return "", fmt.Errorf(
			"unsupported match factor %q",
			value,
		)
	}
}
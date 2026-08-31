package geminiranker

import (
	"context"
	"encoding/json"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
	geminiclient "github.com/AppeiYA/consultation-platform/internal/shared/adapters/gemini"
	"google.golang.org/genai"
)

type CandidateRanker struct {
	client *geminiclient.Client
}

func NewCandidateRanker(client *geminiclient.Client) *CandidateRanker {
	return &CandidateRanker{
		client: client,
	}
}

func (r *CandidateRanker) Rank(
	ctx context.Context,
	req outbound.RankingRequest,
) ([]domain.RankedCandidate, error) {
	if len(req.CandidatePool.Candidates()) == 0 {
		return []domain.RankedCandidate{}, nil
	}

	prompt := buildRankingPrompt(req)

	response, err := r.client.GenerateContent(
		ctx, 
		[]*genai.Content{
			genai.NewContentFromText(prompt, genai.RoleUser),
		}, 
		rankingConfig(),
	)

	if err != nil {
		return nil, ErrGeneratingCandidateRanking("error generating candidates ranking", err)
	}

	text := response.Text()

	var result rankingResponse

	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, ErrGeneratingCandidateRanking("decode gemini ranking response", err)
	}

	return mapRankingResponse(req, result)
}

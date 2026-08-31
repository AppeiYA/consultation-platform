package geminiranker

type rankingResponse struct {
	Candidates []rankedCandidateResponse `json:"candidates"`
}

type rankedCandidateResponse struct {
	ConsultantID string           `json:"consultant_id"`
	Score        float64          `json:"score"`
	Reasons      []reasonResponse `json:"reasons"`
}

type reasonResponse struct {
	Factor string `json:"factor"`
	Detail string `json:"detail"`
}
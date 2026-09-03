package app

import (
	"github.com/AppeiYA/consultation-platform/internal/expertmatching"
	expertmatching_http "github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/inbound/http"
	consultant_adapter "github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/outbound/external/consultant"
	consultationcase_adapter "github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/outbound/external/consultationcase"
	expertMatchingRepo "github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/outbound/postgres/expertmatching"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/outbound/ranker"
	geminiranker "github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/outbound/ranker/gemini"
	matching_worker "github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/worker"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
	geminiclient "github.com/AppeiYA/consultation-platform/internal/shared/adapters/gemini"
	"github.com/AppeiYA/consultation-platform/internal/shared/db"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
	"github.com/hibiken/asynq"
)

func (a *App) registerExpertMatchingModule(
	repository db.Repository,
	asynqClient *asynq.Client,
	clock shared_outbound.Clock,
	idGenerator shared_outbound.IdentifierGenerator,
	geminiClient *geminiclient.Client,
) {
	caseReader := consultationcase_adapter.NewCaseReaderAdapter(&repository)
	candidateGen := consultant_adapter.NewCandidateGeneratorAdapter(&repository)

	var candidateRanker outbound.CandidateRanker = ranker.NewRuleBasedRanker()
	if geminiClient != nil {
		candidateRanker = geminiranker.NewCandidateRanker(geminiClient)
	}

	runRepo := expertMatchingRepo.NewExpertMatchingRepository(&repository)
	dispatcher := matching_worker.NewMatchingJobDispatcher(asynqClient)

	a.expertMatchingModule = expertmatching.NewModule(
		caseReader,
		candidateGen,
		candidateRanker,
		runRepo,
		dispatcher,
		idGenerator,
		clock,
	)
	a.expertMatchingHandler = expertmatching_http.NewExpertMatchingHandler(
		a.expertMatchingModule,
	)
}

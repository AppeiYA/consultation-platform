package app

import (
	"github.com/AppeiYA/consultation-platform/internal/expertmatching"
	expertmatching_http "github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/inbound/http"
	consultant_adapter "github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/outbound/external/consultant"
	consultationcase_adapter "github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/outbound/external/consultationcase"
	expertMatchingRepo "github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/outbound/postgres/expertmatching"
	// "github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/outbound/ranker"
	geminiranker "github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/outbound/ranker/gemini"
	redis_adapter "github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/outbound/redis"
	geminiclient "github.com/AppeiYA/consultation-platform/internal/shared/adapters/gemini"
	"github.com/AppeiYA/consultation-platform/internal/shared/db"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/shared/redis"
)

func (a *App) registerExpertMatchingModule(
	repository db.Repository,
	redisClient *redis.Redis,
	clock shared_outbound.Clock,
	idGenerator shared_outbound.IdentifierGenerator,
	geminiClient *geminiclient.Client,
) {
	caseReader := consultationcase_adapter.NewCaseReaderAdapter(&repository)
	candidateGen := consultant_adapter.NewCandidateGeneratorAdapter(&repository)
	geminiRanker := geminiranker.NewCandidateRanker(geminiClient)
	// candidateRanker := ranker.NewRuleBasedRanker()
	runRepo := expertMatchingRepo.NewExpertMatchingRepository(&repository)
	jobEnqueuer := redis_adapter.NewRedisMatchingJobEnqueuer(redisClient)

	a.expertMatchingModule = expertmatching.NewModule(
		caseReader,
		candidateGen,
		geminiRanker,
		runRepo,
		jobEnqueuer,
		idGenerator,
		clock,
	)
	a.expertMatchingHandler = expertmatching_http.NewExpertMatchingHandler(
		a.expertMatchingModule,
	)
}

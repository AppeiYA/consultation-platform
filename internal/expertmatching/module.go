package expertmatching

import (
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/inbound"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/usecase"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
)

type Module struct {
	StartMatching     inbound.StartMatchingInt
	ExecuteMatching   inbound.ExecuteMatchingInt
	GetMatchingResult inbound.GetMatchingResultInt
}

func NewModule(
	caseReader outbound.CaseReader,
	candidateGenerator outbound.CandidateGenerator,
	candidateRanker outbound.CandidateRanker,
	runRepository outbound.MatchingRunRepository,
	jobEnqueuer outbound.MatchingJobEnqueuer,
	idGenerator shared_outbound.IdentifierGenerator,
	clock shared_outbound.Clock,
) *Module {
	return &Module{
		StartMatching: usecase.NewStartMatchingUsecase(
			caseReader,
			runRepository,
			jobEnqueuer,
			idGenerator,
			clock,
		),
		ExecuteMatching: usecase.NewExecuteMatchingUsecase(
			caseReader,
			candidateGenerator,
			candidateRanker,
			runRepository,
			clock,
		),
		GetMatchingResult: usecase.NewGetMatchingResultUsecase(runRepository),
	}
}

package expertmatching

import (
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/inbound"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/usecase"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
)

type Module struct {
	StartMatching     inbound.StartMatchingInt
	ProcessMatching   inbound.ProcessMatchingInt
	GetMatchingResult inbound.GetMatchingResultInt
}

func NewModule(
	caseReader outbound.CaseReader,
	candidateGenerator outbound.CandidateGenerator,
	candidateRanker outbound.CandidateRanker,
	runRepository outbound.MatchingRunRepository,
	dispatcher outbound.MatchingJobDispatcher,
	idGenerator shared_outbound.IdentifierGenerator,
	clock shared_outbound.Clock,
) *Module {
	return &Module{
		StartMatching: usecase.NewStartMatchingUsecase(
			caseReader,
			runRepository,
			dispatcher,
			idGenerator,
			clock,
		),
		ProcessMatching: usecase.NewProcessMatchingUsecase(
			caseReader,
			candidateGenerator,
			candidateRanker,
			runRepository,
			clock,
		),
		GetMatchingResult: usecase.NewGetMatchingResultUsecase(runRepository),
	}
}

package consultationcase

import (
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/ports/inbound"
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/usecase"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
)

type Module struct {
	SaveCase inbound.SaveCaseInt
	GetCase  inbound.GetCaseInt
	ListCase inbound.ListCaseInt
	DeleteCase inbound.DeleteCaseInt
	UpdateCase inbound.UpdateCaseInt
}

func NewModule(
	caseRepo outbound.CaseRepository,
	clientVerifier outbound.ClientVerifier,
	idGenerator shared_outbound.IdentifierGenerator,
	clock shared_outbound.Clock,
) *Module {
	return &Module{
		SaveCase: usecase.NewCreateCaseUsecase(
			caseRepo,
			idGenerator,
			clientVerifier,
			clock,
		),
		GetCase: usecase.NewGetCaseUsecase(caseRepo),
		ListCase: usecase.NewListCasesUsecase(caseRepo),
		DeleteCase: usecase.NewDeleteCaseUsecase(caseRepo),
		UpdateCase: usecase.NewUpdateCaseUsecase(caseRepo, clock),
	}
}
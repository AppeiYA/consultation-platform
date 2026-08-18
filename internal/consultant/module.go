package consultant

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/inbound"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase"
)

type Module struct {
	RegisterConsultant inbound.RegisterConsultantInt
	GetConsultant      inbound.GetConsultantInt
	UpdateConsultant inbound.UpdateConsultantInt
	SubmitVerification inbound.SubmitVerificationInt
}

func NewModule(
	consultantRepo outbound.ConsultantRepository,
	verificationService outbound.VerificationService,
	verificationRepository outbound.VerificationRepository,
	idGenerator shared_outbound.IdentifierGenerator,
	clock shared_outbound.Clock,
) *Module {
	return &Module{
		RegisterConsultant: usecase.NewRegisterConsultantUsecase(
			consultantRepo, 
			idGenerator, 
			clock,
		),
		GetConsultant: usecase.NewGetConsultantUsecase(
			consultantRepo,
		),
		UpdateConsultant: usecase.NewUpdateConsultantUsecase(
			consultantRepo,
			clock,
		),
		SubmitVerification: usecase.NewSubmitVerificationUsecase(
			consultantRepo,
			idGenerator,
			clock,
			verificationService, 
			verificationRepository,
		),
	}
}
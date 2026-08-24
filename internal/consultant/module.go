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
	UpdateConsultant   inbound.UpdateConsultantInt
	SubmitVerification inbound.SubmitVerificationInt
	CreateAvailability inbound.CreateAvailabilityInt
	GetAvailability    inbound.GetAvailabilityInt
	UpdateAvailability inbound.UpdateAvailabilityInt
	ActivateAvailability inbound.ActivateAvailabilityInt
	DeactivateAvailability inbound.DeactivateAvailabilityInt
	GetProfession      inbound.GetProfessionInt
	ListProfessions     inbound.ListProfessionsInt
}

func NewModule(
	consultantRepo outbound.ConsultantRepository,
	verificationService outbound.VerificationService,
	verificationRepository outbound.VerificationRepository,
	availabilityRepository outbound.AvailabilityRepository,
	professionRepo outbound.ProfessionRepository,
	identityAdapter outbound.RoleAssigner,
	idGenerator shared_outbound.IdentifierGenerator,
	clock shared_outbound.Clock,
) *Module {
	return &Module{
		RegisterConsultant: usecase.NewRegisterConsultantUsecase(
			consultantRepo, 
			professionRepo,
			identityAdapter,
			idGenerator, 
			clock,
		),
		GetConsultant: usecase.NewGetConsultantUsecase(
			consultantRepo,
			professionRepo,
		),
		UpdateConsultant: usecase.NewUpdateConsultantUsecase(
			consultantRepo,
			professionRepo,
			clock,
		),
		SubmitVerification: usecase.NewSubmitVerificationUsecase(
			consultantRepo,
			idGenerator,
			clock,
			verificationService, 
			verificationRepository,
		),
		CreateAvailability: usecase.NewCreateAvailabilityUsecase(
			availabilityRepository,
			idGenerator,
			consultantRepo,
			clock,
		),
		GetAvailability: usecase.NewGetAvailabilityUsecase(
			availabilityRepository,
		),
		UpdateAvailability: usecase.NewUpdateAvailabilityUsecase(
			availabilityRepository,
			consultantRepo,
			clock,
		),
		ActivateAvailability: usecase.NewActivateAvailabilityUsecase(
			availabilityRepository,
			consultantRepo,
			clock,
		),
		DeactivateAvailability: usecase.NewDeactivateAvailabilityUsecase(
			availabilityRepository,
			consultantRepo,
			clock,
		),
		GetProfession: usecase.NewGetProfessionUsecase(
			professionRepo,
		),
		ListProfessions: usecase.NewListProfessionsUsecase(
			professionRepo,
		),
	}
}
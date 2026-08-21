package app

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant"
	consultant_http "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http"

	consultantRepo "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/outbound/postgres/consultant"
	consultantAvailabilityRepo "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/outbound/postgres/consultant_availability"
	consultantVerificationRepo "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/outbound/postgres/consultant_verification"
	professionRepo "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/outbound/postgres/profession"
	"github.com/AppeiYA/consultation-platform/internal/consultant/adapters/outbound/verification"
	"github.com/AppeiYA/consultation-platform/internal/shared/adapters/id/uuid"
	system "github.com/AppeiYA/consultation-platform/internal/shared/adapters/outbound/clock"
	"github.com/AppeiYA/consultation-platform/internal/shared/db"
)

func (a *App) registerConsultantModule(
	db db.Repository,
	clock *system.SystemClock,
	idGenerator *uuid.Generator,
) {
	verificationService := verification.UnavailableVerificationService{}

	consultantRepo := consultantRepo.NewConsultantRepository(db, clock)
	verificationRepo := consultantVerificationRepo.NewVerificationRepository(db, clock)
	availabilityRepo := consultantAvailabilityRepo.NewAvailabilityRepository(db, clock)
	professionRepo := professionRepo.NewProfessionRepository(db, clock)

	consultantModule := consultant.NewModule(
		consultantRepo,
		&verificationService,
		verificationRepo,
		availabilityRepo,
		professionRepo,
		idGenerator,
		clock,
	)

	consultantHandler := consultant_http.NewConsultantHandler(*consultantModule)
	a.consultantModule = consultantModule
	a.consultantHandler = consultantHandler
}
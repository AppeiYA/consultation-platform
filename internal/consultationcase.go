package app

import (
	"github.com/AppeiYA/consultation-platform/internal/consultationcase"
	consultationcase_http "github.com/AppeiYA/consultation-platform/internal/consultationcase/adapters/inbound/http"
	consultationCaseIdentityAdapter "github.com/AppeiYA/consultation-platform/internal/consultationcase/adapters/outbound/external/identity"
	consultationCaseRepo "github.com/AppeiYA/consultation-platform/internal/consultationcase/adapters/outbound/postgres/consultationcase"
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/shared/adapters/id/uuid"
	system "github.com/AppeiYA/consultation-platform/internal/shared/adapters/outbound/clock"
	"github.com/AppeiYA/consultation-platform/internal/shared/db"
)

func (a *App) registerConsultationCaseModule(
	db db.Repository,
	clientVerifier *consultationCaseIdentityAdapter.ClientVerifier,
	matchingStarter outbound.MatchingStarter,
	clock *system.SystemClock,
	idGenerator *uuid.Generator,
) {
	caseRepository := consultationCaseRepo.NewConsultationCaseRepository(&db)

	consultationCaseModule := consultationcase.NewModule(
		caseRepository,
		clientVerifier,
		matchingStarter,
		idGenerator,
		clock,
	)

	consultationCaseHandler := consultationcase_http.NewConsultationCaseHandler(
		consultationCaseModule,
	)

	a.consultationCaseModule = consultationCaseModule
	a.consultationCaseHandler = consultationCaseHandler
}

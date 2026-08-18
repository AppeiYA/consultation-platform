package app

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant"
	consultant_http "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http"
	consultant_postgres "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/outbound/postgres"
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

	consultantRepo := consultant_postgres.NewConsultantRepository(db, clock)
	verificationRepo := consultant_postgres.NewVerificationRepository(db, clock)

	consultantModule := consultant.NewModule(
		consultantRepo,
		&verificationService,
		verificationRepo,
		idGenerator,
		clock,
	)

	consultantHandler := consultant_http.NewConsultantHandler(*consultantModule)
	a.consultantModule = consultantModule
	a.consultantHandler = consultantHandler
}
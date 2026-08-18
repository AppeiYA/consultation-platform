package usecase

import (
	"context"
	"errors"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
	shared_outbound "github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
)

type SubmitVerificationUsecase struct {
	consultantRepository   outbound.ConsultantRepository
	idGenerator            shared_outbound.IdentifierGenerator
	clock                  shared_outbound.Clock
	verificationService    outbound.VerificationService
	verificationRepository outbound.VerificationRepository
}

func NewSubmitVerificationUsecase(
	consultantRepository outbound.ConsultantRepository,
	idGenerator shared_outbound.IdentifierGenerator,
	clock shared_outbound.Clock,
	verificationService outbound.VerificationService,
	verificationRepository outbound.VerificationRepository,
) *SubmitVerificationUsecase {
	return &SubmitVerificationUsecase{
		consultantRepository:   consultantRepository,
		idGenerator:            idGenerator,
		clock:                  clock,
		verificationService:    verificationService,
		verificationRepository: verificationRepository,
	}
}

func (u *SubmitVerificationUsecase) Execute(ctx context.Context, userID string) (*dto.SubmitVerificationResponse, error) {
	consultant, err := u.consultantRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if consultant == nil {
		return nil, domain.ErrConsultantNotFound
	}

	verification, err := u.verificationRepository.FindByConsultantID(ctx, consultant.ID())
	if err != nil && !errors.Is(err, domain.ErrConsultantVerificationNotFound) {
		return nil, err
	}

	if verification != nil {
		if err := u.evaluateExistingVerification(verification); err != nil {
			return nil, err
		}
	}

	res, err := u.verificationService.CreateInquiry(ctx, consultant.ID())
	if err != nil {
		return nil, err
	}

	verificationID, err := u.idGenerator.Generate(domain.VerificationIDPrefix)
	if err != nil {
		return nil, err
	}

	newVerification := domain.NewVerification(
		verificationID,
		consultant.ID(),
		res.Provider,
		res.ProviderReference,
		domain.VerificationStatusPending,
		u.clock.Now(),
	)

	if err := u.verificationRepository.Save(ctx, newVerification); err != nil {
		return nil, err
	}

	return &dto.SubmitVerificationResponse{
		VerificationID:    verificationID,
		ProviderReference: res.ProviderReference,
		VerificationURL:   res.VerificationURL,
		Status:            string(domain.VerificationStatusPending),
	}, nil
}

func (u *SubmitVerificationUsecase) evaluateExistingVerification(
	verification *domain.ConsultantVerification,
) error {
	switch verification.Status() {
	case domain.VerificationStatusPending:
		return domain.ErrVerificationPending

	case domain.VerificationStatusReviewing:
		return domain.ErrVerificationInReview

	case domain.VerificationStatusApproved:
		return domain.ErrVerificationAlreadyApproved

	case domain.VerificationStatusRejected:
		return nil

	default:
		return domain.ErrInvalidVerificationStatus
	}
}
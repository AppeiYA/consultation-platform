package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/mocks"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
	shared_mocks "github.com/AppeiYA/consultation-platform/internal/shared/mocks"
)

type testSubmitVerification struct {
	consultantRepo   *mocks.MockConsultantRepository
	verificationRepo *mocks.MockVerificationRepository
	verificationServ *mocks.MockVerificationService
	idGenerator      *shared_mocks.MockIDGenerator
	clock            *shared_mocks.MockClock

	sut *SubmitVerificationUsecase
}

func setUpSubmitVerification(t *testing.T) *testSubmitVerification {
	t.Helper()

	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	consultant, _ := domain.ReconstitueConsultant(
		"con_123",
		"user_123",
		"SOFTWARE_ENGINEER",
		"Jane Doe Tech",
		"Experienced software engineer.",
		10,
		true,
		now,
		now,
	)

	consultantRepo := &mocks.MockConsultantRepository{
		FindByUserIDFn: func(ctx context.Context, userID string) (*domain.Consultant, error) {
			return consultant, nil
		},
	}

	verificationRepo := &mocks.MockVerificationRepository{
		FindByConsultantIDFn: func(ctx context.Context, consultantID string) (*domain.ConsultantVerification, error) {
			return nil, domain.ErrConsultantVerificationNotFound
		},
		SaveFn: func(ctx context.Context, verification *domain.ConsultantVerification) error {
			return nil
		},
	}

	verificationServ := &mocks.MockVerificationService{
		CreateInquiryFn: func(ctx context.Context, consultantID string) (*outbound.CreateInquiryResult, error) {
			return &outbound.CreateInquiryResult{
				Provider:          "persona",
				ProviderReference: "inq_999",
				VerificationURL:   "https://withpersona.com/verify?inquiry=inq_999",
			}, nil
		},
	}

	idGenerator := &shared_mocks.MockIDGenerator{
		GenerateFn: func(prefix string) (string, error) {
			return "ver_456", nil
		},
	}

	clock := &shared_mocks.MockClock{
		NowFn: func() time.Time {
			return now
		},
	}

	sut := NewSubmitVerificationUsecase(
		consultantRepo,
		idGenerator,
		clock,
		verificationServ,
		verificationRepo,
	)

	return &testSubmitVerification{
		consultantRepo:   consultantRepo,
		verificationRepo: verificationRepo,
		verificationServ: verificationServ,
		idGenerator:      idGenerator,
		clock:            clock,
		sut:              sut,
	}
}

func TestSubmitVerification_Execute(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	userID := "user_123"

	t.Run("should submit verification successfully when no previous verification exists", func(t *testing.T) {
		tc := setUpSubmitVerification(t)

		var savedVerification *domain.ConsultantVerification
		tc.verificationRepo.SaveFn = func(ctx context.Context, verification *domain.ConsultantVerification) error {
			savedVerification = verification
			return nil
		}

		res, err := tc.sut.Execute(context.Background(), userID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if res == nil {
			t.Fatal("expected non-nil response")
		}

		if res.VerificationID != "ver_456" {
			t.Errorf("expected VerificationID ver_456, got %s", res.VerificationID)
		}
		if res.ProviderReference != "inq_999" {
			t.Errorf("expected ProviderReference inq_999, got %s", res.ProviderReference)
		}
		if res.VerificationURL != "https://withpersona.com/verify?inquiry=inq_999" {
			t.Errorf("expected VerificationURL https://withpersona.com/verify?inquiry=inq_999, got %s", res.VerificationURL)
		}
		if res.Status != string(domain.VerificationStatusPending) {
			t.Errorf("expected Status %s, got %s", domain.VerificationStatusPending, res.Status)
		}

		if savedVerification == nil {
			t.Fatal("expected verification to be saved")
		}
		if savedVerification.ID() != "ver_456" {
			t.Errorf("expected saved ID ver_456, got %s", savedVerification.ID())
		}
		if savedVerification.ConsultantID() != "con_123" {
			t.Errorf("expected ConsultantID con_123, got %s", savedVerification.ConsultantID())
		}
		if savedVerification.Provider() != "persona" {
			t.Errorf("expected Provider persona, got %s", savedVerification.Provider())
		}
		if savedVerification.ProviderReference() != "inq_999" {
			t.Errorf("expected ProviderReference inq_999, got %s", savedVerification.ProviderReference())
		}
		if savedVerification.Status() != domain.VerificationStatusPending {
			t.Errorf("expected Status %s, got %s", domain.VerificationStatusPending, savedVerification.Status())
		}
		if savedVerification.SubmittedAt() == nil || !savedVerification.SubmittedAt().Equal(now) {
			t.Errorf("expected SubmittedAt %v, got %v", now, savedVerification.SubmittedAt())
		}
	})

	t.Run("should submit verification successfully when previous verification was rejected", func(t *testing.T) {
		tc := setUpSubmitVerification(t)

		rejectedVerification, _ := domain.ReconstitueConsultantVerification(
			"ver_old",
			"con_123",
			"persona",
			"inq_old",
			domain.VerificationStatusRejected,
			&now,
			&now,
			now,
			now,
		)

		tc.verificationRepo.FindByConsultantIDFn = func(ctx context.Context, consultantID string) (*domain.ConsultantVerification, error) {
			return rejectedVerification, nil
		}

		res, err := tc.sut.Execute(context.Background(), userID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if res == nil {
			t.Fatal("expected non-nil response")
		}
		if res.VerificationID != "ver_456" {
			t.Errorf("expected VerificationID ver_456, got %s", res.VerificationID)
		}
	})

	t.Run("should fail when consultant is not found", func(t *testing.T) {
		tc := setUpSubmitVerification(t)
		tc.consultantRepo.FindByUserIDFn = func(ctx context.Context, userID string) (*domain.Consultant, error) {
			return nil, nil
		}

		res, err := tc.sut.Execute(context.Background(), userID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if res != nil {
			t.Errorf("expected nil response, got %v", res)
		}
		if !errors.Is(err, domain.ErrConsultantNotFound) {
			t.Errorf("expected ErrConsultantNotFound, got %v", err)
		}
	})

	t.Run("should fail when consultant repository returns error", func(t *testing.T) {
		tc := setUpSubmitVerification(t)
		expectedErr := errors.New("db error")
		tc.consultantRepo.FindByUserIDFn = func(ctx context.Context, userID string) (*domain.Consultant, error) {
			return nil, expectedErr
		}

		res, err := tc.sut.Execute(context.Background(), userID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected %v, got %v", expectedErr, err)
		}
		if res != nil {
			t.Errorf("expected nil response, got %v", res)
		}
	})

	t.Run("should fail when verification repository returns unexpected error", func(t *testing.T) {
		tc := setUpSubmitVerification(t)
		expectedErr := errors.New("verification db error")
		tc.verificationRepo.FindByConsultantIDFn = func(ctx context.Context, consultantID string) (*domain.ConsultantVerification, error) {
			return nil, expectedErr
		}

		res, err := tc.sut.Execute(context.Background(), userID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected %v, got %v", expectedErr, err)
		}
		if res != nil {
			t.Errorf("expected nil response, got %v", res)
		}
	})

	t.Run("should fail when existing verification is pending", func(t *testing.T) {
		tc := setUpSubmitVerification(t)
		pendingVerification, _ := domain.ReconstitueConsultantVerification(
			"ver_pending",
			"con_123",
			"persona",
			"inq_pending",
			domain.VerificationStatusPending,
			&now,
			nil,
			now,
			now,
		)
		tc.verificationRepo.FindByConsultantIDFn = func(ctx context.Context, consultantID string) (*domain.ConsultantVerification, error) {
			return pendingVerification, nil
		}

		res, err := tc.sut.Execute(context.Background(), userID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrVerificationPending) {
			t.Errorf("expected ErrVerificationPending, got %v", err)
		}
		if res != nil {
			t.Errorf("expected nil response, got %v", res)
		}
	})

	t.Run("should fail when existing verification is in review", func(t *testing.T) {
		tc := setUpSubmitVerification(t)
		reviewVerification, _ := domain.ReconstitueConsultantVerification(
			"ver_review",
			"con_123",
			"persona",
			"inq_review",
			domain.VerificationStatusReviewing,
			&now,
			nil,
			now,
			now,
		)
		tc.verificationRepo.FindByConsultantIDFn = func(ctx context.Context, consultantID string) (*domain.ConsultantVerification, error) {
			return reviewVerification, nil
		}

		res, err := tc.sut.Execute(context.Background(), userID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrVerificationInReview) {
			t.Errorf("expected ErrVerificationInReview, got %v", err)
		}
		if res != nil {
			t.Errorf("expected nil response, got %v", res)
		}
	})

	t.Run("should fail when existing verification is already approved", func(t *testing.T) {
		tc := setUpSubmitVerification(t)
		approvedVerification, _ := domain.ReconstitueConsultantVerification(
			"ver_appr",
			"con_123",
			"persona",
			"inq_appr",
			domain.VerificationStatusApproved,
			&now,
			&now,
			now,
			now,
		)
		tc.verificationRepo.FindByConsultantIDFn = func(ctx context.Context, consultantID string) (*domain.ConsultantVerification, error) {
			return approvedVerification, nil
		}

		res, err := tc.sut.Execute(context.Background(), userID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrVerificationAlreadyApproved) {
			t.Errorf("expected ErrVerificationAlreadyApproved, got %v", err)
		}
		if res != nil {
			t.Errorf("expected nil response, got %v", res)
		}
	})

	t.Run("should fail when verification service returns error", func(t *testing.T) {
		tc := setUpSubmitVerification(t)
		expectedErr := errors.New("persona api down")
		tc.verificationServ.CreateInquiryFn = func(ctx context.Context, consultantID string) (*outbound.CreateInquiryResult, error) {
			return nil, expectedErr
		}

		res, err := tc.sut.Execute(context.Background(), userID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected %v, got %v", expectedErr, err)
		}
		if res != nil {
			t.Errorf("expected nil response, got %v", res)
		}
	})

	t.Run("should fail when id generator returns error", func(t *testing.T) {
		tc := setUpSubmitVerification(t)
		expectedErr := errors.New("id generator error")
		tc.idGenerator.GenerateFn = func(prefix string) (string, error) {
			return "", expectedErr
		}

		res, err := tc.sut.Execute(context.Background(), userID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected %v, got %v", expectedErr, err)
		}
		if res != nil {
			t.Errorf("expected nil response, got %v", res)
		}
	})

	t.Run("should fail when verification repository save returns error", func(t *testing.T) {
		tc := setUpSubmitVerification(t)
		expectedErr := errors.New("failed to save verification")
		tc.verificationRepo.SaveFn = func(ctx context.Context, verification *domain.ConsultantVerification) error {
			return expectedErr
		}

		res, err := tc.sut.Execute(context.Background(), userID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected %v, got %v", expectedErr, err)
		}
		if res != nil {
			t.Errorf("expected nil response, got %v", res)
		}
	})
}

package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	consultant_dto "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/consultant/mocks"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/shared/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestSubmitVerification_Integration(t *testing.T) {
	validConsultantReq := consultant_dto.RegisterConsultantDTO{
		ProfessionID:    "prof_9ee432d7-b672-40ae-b03f-c1f1fb696621",
		DisplayName:     "Jane Doe Tech",
		Bio:             "Experienced software engineer with 10 years of experience.",
		YearsExperience: 10,
	}

	mockVerificationService := &mocks.MockVerificationService{
		CreateInquiryFn: func(ctx context.Context, consultantID string) (*outbound.CreateInquiryResult, error) {
			return &outbound.CreateInquiryResult{
				Provider:          "persona",
				ProviderReference: "inq_test_12345",
				VerificationURL:   "https://withpersona.com/verify?inquiry=inq_test_12345",
			}, nil
		},
	}

	t.Run("should submit verification successfully for registered consultant", func(t *testing.T) {
		harness := setUpConsultantApp(t, WithVerificationService(mockVerificationService))
		cookieHeader := registerAndLoginUser(t, harness, "verify.success@example.com")

		// 1. Register consultant
		regResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			validConsultantReq,
			cookieHeader,
		)
		require.NoError(t, err)
		regResp.Body.Close()
		require.Equal(t, http.StatusCreated, regResp.StatusCode)

		// 2. Submit verification
		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/verification",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		res := testhelpers.DecodeResponse(t, resp)
		require.True(t, res.Success)
		require.Equal(t, "Verification submitted successfully", res.Message)

		dataBytes, err := json.Marshal(res.Data)
		require.NoError(t, err)

		var verifyResp consultant_dto.SubmitVerificationResponseDTO
		err = json.Unmarshal(dataBytes, &verifyResp)
		require.NoError(t, err)

		require.NotEmpty(t, verifyResp.VerificationID)
		require.Equal(t, "inq_test_12345", verifyResp.ProviderReference)
		require.Equal(t, "https://withpersona.com/verify?inquiry=inq_test_12345", verifyResp.VerificationURL)
		require.Equal(t, "PENDING", verifyResp.Status)
	})

	t.Run("should fail when verification is already pending", func(t *testing.T) {
		harness := setUpConsultantApp(t, WithVerificationService(mockVerificationService))
		cookieHeader := registerAndLoginUser(t, harness, "verify.duplicate@example.com")

		// 1. Register consultant
		regResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			validConsultantReq,
			cookieHeader,
		)
		require.NoError(t, err)
		regResp.Body.Close()
		require.Equal(t, http.StatusCreated, regResp.StatusCode)

		// 2. First submission - succeeds
		resp1, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/verification",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		resp1.Body.Close()
		require.Equal(t, http.StatusOK, resp1.StatusCode)

		// 3. Second submission - fails because verification is pending
		resp2, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/verification",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp2.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp2.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp2)
		require.False(t, errRes.Success)
		require.Equal(t, "verification is pending", errRes.Message)
	})

	t.Run("should fail when consultant profile does not exist for authenticated user", func(t *testing.T) {
		harness := setUpConsultantApp(t, WithVerificationService(mockVerificationService))
		cookieHeader := registerAndLoginUser(t, harness, "verify.noconsultant@example.com")

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/verification",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "consultant not found", errRes.Message)
	})

	t.Run("should fail when unauthenticated (no session cookie)", func(t *testing.T) {
		harness := setUpConsultantApp(t, WithVerificationService(mockVerificationService))

		resp, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/verification",
			nil,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "User is unauthorized", errRes.Message)
	})

	t.Run("should fail when verification service is unavailable", func(t *testing.T) {
		// Default harness uses UnavailableVerificationService
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "verify.unavailable@example.com")

		// 1. Register consultant
		regResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			validConsultantReq,
			cookieHeader,
		)
		require.NoError(t, err)
		regResp.Body.Close()
		require.Equal(t, http.StatusCreated, regResp.StatusCode)

		// 2. Submit verification
		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/verification",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "verification service unavailable", errRes.Message)
	})
}

package integration

import (
	"net/http"
	"strings"
	"testing"

	case_dto "github.com/AppeiYA/consultation-platform/internal/consultationcase/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestCreateConsultationCase_Integration(t *testing.T) {
	validCaseReq := case_dto.CreateConsultationCaseDTO{
		Title:       "Corporate Tax Optimization",
		Description: "Seeking professional consultation for corporate tax restructuring.",
		Category:    "FINANCE",
	}

	t.Run("should create consultation case successfully for authenticated user", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "client.create@example.com")

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultation-cases",
			validCaseReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		res := testhelpers.DecodeResponse(t, resp)
		require.True(t, res.Success)
		require.Equal(t, "Case created successfully", res.Message)
	})

	t.Run("should fail when unauthenticated (no session cookie)", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)

		resp, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodPost,
			"/test/v1/consultation-cases",
			validCaseReq,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "User is unauthorized", errRes.Message)
	})

	t.Run("should fail when session cookie is invalid", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultation-cases",
			validCaseReq,
			"test_session_id=invalid_token_12345678901234567890123456789012",
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "User is unauthorized", errRes.Message)
	})

	t.Run("should fail when request body is invalid JSON", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "badbody@example.com")

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultation-cases",
			"{invalid json payload",
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail when title is empty", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "notitle@example.com")

		badReq := validCaseReq
		badReq.Title = ""

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultation-cases",
			badReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "title is required", errRes.Message)
	})

	t.Run("should fail when title exceeds 100 characters", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "longtitle@example.com")

		badReq := validCaseReq
		badReq.Title = strings.Repeat("x", 101)

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultation-cases",
			badReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "Case title is too long", errRes.Message)
	})

	t.Run("should fail when description is empty", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "nodesc@example.com")

		badReq := validCaseReq
		badReq.Description = ""

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultation-cases",
			badReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "description is required", errRes.Message)
	})

	t.Run("should fail when description exceeds 500 characters", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "longdesc@example.com")

		badReq := validCaseReq
		badReq.Description = strings.Repeat("y", 501)

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultation-cases",
			badReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "Case description is too long", errRes.Message)
	})

	t.Run("should fail when category is empty", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "nocat@example.com")

		badReq := validCaseReq
		badReq.Category = ""

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultation-cases",
			badReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "category is required", errRes.Message)
	})
}

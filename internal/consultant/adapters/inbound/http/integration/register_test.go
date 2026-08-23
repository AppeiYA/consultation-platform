package integration

import (
	"net/http"
	"testing"

	consultant_dto "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	identity_dto "github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/testhelpers"
	"github.com/stretchr/testify/require"
)

func registerAndLoginUser(t *testing.T, harness *TestHarness, email string) string {
	t.Helper()

	// 1. Register user
	regBody := identity_dto.RegisterRequest{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     email,
		Password:  "SecurePassword123!",
	}

	regResp, err := testhelpers.PerformRequest(
		harness.App,
		http.MethodPost,
		"/test/v1/identity/register",
		regBody,
	)
	require.NoError(t, err)
	regResp.Body.Close()
	require.Equal(t, http.StatusCreated, regResp.StatusCode)

	// 2. Login user to retrieve session cookie
	loginBody := identity_dto.LoginRequest{
		Email:    email,
		Password: "SecurePassword123!",
	}

	loginResp, err := testhelpers.PerformRequest(
		harness.App,
		http.MethodPost,
		"/test/v1/identity/login",
		loginBody,
	)
	require.NoError(t, err)
	defer loginResp.Body.Close()
	require.Equal(t, http.StatusOK, loginResp.StatusCode)

	cookieHeader := testhelpers.ExtractCookieHeader(loginResp)
	require.NotEmpty(t, cookieHeader)

	return cookieHeader
}

func TestRegisterConsultant_Integration(t *testing.T) {
	validConsultantReq := consultant_dto.RegisterConsultantDTO{
		ProfessionID:    "prof_9ee432d7-b672-40ae-b03f-c1f1fb696621",
		DisplayName:     "Jane Doe Tech",
		Bio:             "Experienced software engineer with 10 years of experience.",
		YearsExperience: 10,
	}

	t.Run("should register consultant successfully", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "consultant.jane@example.com")

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			validConsultantReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		res := testhelpers.DecodeResponse(t, resp)
		require.True(t, res.Success)
		require.Equal(t, "Consultant profile created. Please verify", res.Message)
	})

	t.Run("should fail when unauthenticated (no session cookie)", func(t *testing.T) {
		harness := setUpConsultantApp(t)

		resp, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			validConsultantReq,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "User is unauthorized", errRes.Message)
	})

	t.Run("should fail when session cookie is invalid", func(t *testing.T) {
		harness := setUpConsultantApp(t)

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			validConsultantReq,
			"test_session_id=invalid_token_12345678901234567890123456789012",
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "User is unauthorized", errRes.Message)
	})

	t.Run("should fail when request body is invalid", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "invalid.body@example.com")

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			"{invalid json payload",
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail when profession is invalid", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "invalid.prof@example.com")

		invalidReq := validConsultantReq
		invalidReq.ProfessionID = "INVALID_PROFESSION"

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			invalidReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "profession is not in system", errRes.Message)
	})

	t.Run("should fail when display name is invalid", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "invalid.name@example.com")

		invalidReq := validConsultantReq
		invalidReq.DisplayName = "Short" // len < 6

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			invalidReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "Display name length must be > 6 and < 20", errRes.Message)
	})

	t.Run("should fail when bio is empty", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "empty.bio@example.com")

		invalidReq := validConsultantReq
		invalidReq.Bio = ""

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			invalidReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "Bio is empty", errRes.Message)
	})

	t.Run("should fail when years of experience is invalid", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "invalid.exp@example.com")

		invalidReq := validConsultantReq
		invalidReq.YearsExperience = 0

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			invalidReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "Years of experience must be > 0 and < 70", errRes.Message)
	})

	t.Run("should fail when consultant profile already exists for user", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "duplicate.consultant@example.com")

		// First registration - should succeed
		resp1, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			validConsultantReq,
			cookieHeader,
		)
		require.NoError(t, err)
		resp1.Body.Close()
		require.Equal(t, http.StatusCreated, resp1.StatusCode)

		// Second registration - should fail with duplicate error
		resp2, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			validConsultantReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp2.Body.Close()

		require.Equal(t, http.StatusConflict, resp2.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp2)
		require.False(t, errRes.Success)
		require.Equal(t, "consultant already exists", errRes.Message)
	})
}

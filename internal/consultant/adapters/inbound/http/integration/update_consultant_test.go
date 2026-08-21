package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	consultant_dto "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestUpdateConsultant_Integration(t *testing.T) {
	initialConsultantReq := consultant_dto.RegisterConsultantDTO{
		ProfessionID:    "prof_9ee432d7-b672-40ae-b03f-c1f1fb696621",
		DisplayName:     "Jane Doe Tech",
		Bio:             "Experienced software engineer with 10 years of experience.",
		YearsExperience: 10,
	}

	validUpdateReq := consultant_dto.UpdateConsultantModel{
		ProfessionID:    "prof_12d965f5-e1f5-49aa-ac57-856772d236ce",
		DisplayName:     "Jane Doe Legal",
		Bio:             "Senior corporate lawyer with over 12 years of trial experience.",
		YearsExperience: 12,
	}

	t.Run("should update consultant profile successfully", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "update.success@example.com")

		// 1. Register consultant
		regResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			initialConsultantReq,
			cookieHeader,
		)
		require.NoError(t, err)
		regResp.Body.Close()
		require.Equal(t, http.StatusCreated, regResp.StatusCode)

		// 2. Update consultant profile via PUT /consultants/profile
		updateResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPut,
			"/test/v1/consultants/profile",
			validUpdateReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer updateResp.Body.Close()

		require.Equal(t, http.StatusOK, updateResp.StatusCode)

		res := testhelpers.DecodeResponse(t, updateResp)
		require.True(t, res.Success)
		require.Equal(t, "Consultant profile updated successfully", res.Message)

		// 3. Fetch consultant profile to verify updated details
		getResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/consultants/user",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer getResp.Body.Close()

		require.Equal(t, http.StatusOK, getResp.StatusCode)

		getRes := testhelpers.DecodeResponse(t, getResp)
		require.True(t, getRes.Success)

		consultantBytes, err := json.Marshal(getRes.Data)
		require.NoError(t, err)

		var consultantResp consultant_dto.PrivateConsultantResponseDTO
		err = json.Unmarshal(consultantBytes, &consultantResp)
		require.NoError(t, err)

		require.Equal(t, "LAWYER", consultantResp.Profession)
		require.Equal(t, "Jane Doe Legal", consultantResp.DisplayName)
		require.Equal(t, "Senior corporate lawyer with over 12 years of trial experience.", consultantResp.Bio)
		require.Equal(t, 12, consultantResp.YearsExperience)
	})

	t.Run("should fail when unauthenticated (no session cookie)", func(t *testing.T) {
		harness := setUpConsultantApp(t)

		resp, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodPut,
			"/test/v1/consultants/profile",
			validUpdateReq,
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
			http.MethodPut,
			"/test/v1/consultants/profile",
			validUpdateReq,
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
		cookieHeader := registerAndLoginUser(t, harness, "update.invalidbody@example.com")

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPut,
			"/test/v1/consultants/profile",
			"{invalid json payload",
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "Invalid request body", errRes.Message)
	})

	t.Run("should fail when consultant profile does not exist for authenticated user", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "update.noconsultant@example.com")

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPut,
			"/test/v1/consultants/profile",
			validUpdateReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "consultant not found", errRes.Message)
	})

	t.Run("should fail when profession is invalid", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "update.invalidprof@example.com")

		// Register consultant first
		regResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			initialConsultantReq,
			cookieHeader,
		)
		require.NoError(t, err)
		regResp.Body.Close()
		require.Equal(t, http.StatusCreated, regResp.StatusCode)

		invalidReq := validUpdateReq
		invalidReq.ProfessionID = "INVALID_PROFESSION"

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPut,
			"/test/v1/consultants/profile",
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
		cookieHeader := registerAndLoginUser(t, harness, "update.invalidname@example.com")

		// Register consultant first
		regResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			initialConsultantReq,
			cookieHeader,
		)
		require.NoError(t, err)
		regResp.Body.Close()
		require.Equal(t, http.StatusCreated, regResp.StatusCode)

		invalidReq := validUpdateReq
		invalidReq.DisplayName = "Short" // len < 6

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPut,
			"/test/v1/consultants/profile",
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
		cookieHeader := registerAndLoginUser(t, harness, "update.emptybio@example.com")

		// Register consultant first
		regResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			initialConsultantReq,
			cookieHeader,
		)
		require.NoError(t, err)
		regResp.Body.Close()
		require.Equal(t, http.StatusCreated, regResp.StatusCode)

		invalidReq := validUpdateReq
		invalidReq.Bio = ""

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPut,
			"/test/v1/consultants/profile",
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
		cookieHeader := registerAndLoginUser(t, harness, "update.invalidexp@example.com")

		// Register consultant first
		regResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			initialConsultantReq,
			cookieHeader,
		)
		require.NoError(t, err)
		regResp.Body.Close()
		require.Equal(t, http.StatusCreated, regResp.StatusCode)

		invalidReq := validUpdateReq
		invalidReq.YearsExperience = 0

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPut,
			"/test/v1/consultants/profile",
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
}

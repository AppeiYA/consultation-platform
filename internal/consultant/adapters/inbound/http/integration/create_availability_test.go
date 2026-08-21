package integration

import (
	"net/http"
	"testing"

	consultant_dto "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestCreateAvailability_Integration(t *testing.T) {
	validConsultantReq := consultant_dto.RegisterConsultantDTO{
		ProfessionID:    "prof_9ee432d7-b672-40ae-b03f-c1f1fb696621",
		DisplayName:     "Jane Doe Tech",
		Bio:             "Experienced software engineer with 10 years of experience.",
		YearsExperience: 10,
	}

	validAvailabilityReq := consultant_dto.CreateAvailabilityRequest{
		DayOfWeek: 1, // Monday
		StartTime: "09:00",
		EndTime:   "11:00",
	}

	t.Run("should create availability successfully for registered consultant", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "avail.success@example.com")

		// 1. Register consultant profile
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

		// 2. Create availability
		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			validAvailabilityReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		res := testhelpers.DecodeResponse(t, resp)
		require.True(t, res.Success)
		require.Equal(t, "Availability created successfully", res.Message)
	})

	t.Run("should create multiple non-overlapping availabilities on the same day", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "avail.multi@example.com")

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

		// 2. Create first slot: 09:00 - 11:00
		resp1, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			consultant_dto.CreateAvailabilityRequest{
				DayOfWeek: 1,
				StartTime: "09:00",
				EndTime:   "11:00",
			},
			cookieHeader,
		)
		require.NoError(t, err)
		resp1.Body.Close()
		require.Equal(t, http.StatusCreated, resp1.StatusCode)

		// 3. Create second slot: 11:00 - 13:00 (adjacent after)
		resp2, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			consultant_dto.CreateAvailabilityRequest{
				DayOfWeek: 1,
				StartTime: "11:00",
				EndTime:   "13:00",
			},
			cookieHeader,
		)
		require.NoError(t, err)
		resp2.Body.Close()
		require.Equal(t, http.StatusCreated, resp2.StatusCode)

		// 4. Create third slot: 14:00 - 16:00
		resp3, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			consultant_dto.CreateAvailabilityRequest{
				DayOfWeek: 1,
				StartTime: "14:00",
				EndTime:   "16:00",
			},
			cookieHeader,
		)
		require.NoError(t, err)
		resp3.Body.Close()
		require.Equal(t, http.StatusCreated, resp3.StatusCode)
	})

	t.Run("should create availability on different days of the week", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "avail.days@example.com")

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

		// Monday slot
		respMon, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			consultant_dto.CreateAvailabilityRequest{
				DayOfWeek: 1,
				StartTime: "09:00",
				EndTime:   "11:00",
			},
			cookieHeader,
		)
		require.NoError(t, err)
		respMon.Body.Close()
		require.Equal(t, http.StatusCreated, respMon.StatusCode)

		// Tuesday slot with same time
		respTue, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			consultant_dto.CreateAvailabilityRequest{
				DayOfWeek: 2,
				StartTime: "09:00",
				EndTime:   "11:00",
			},
			cookieHeader,
		)
		require.NoError(t, err)
		defer respTue.Body.Close()
		require.Equal(t, http.StatusCreated, respTue.StatusCode)
	})

	t.Run("should fail when unauthenticated (no session cookie)", func(t *testing.T) {
		harness := setUpConsultantApp(t)

		resp, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			validAvailabilityReq,
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
			"/test/v1/consultants/availability",
			validAvailabilityReq,
			"test_session_id=invalid_token_12345678901234567890123456789012",
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "User is unauthorized", errRes.Message)
	})

	t.Run("should fail when consultant profile is not found for user", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "noconsultant@example.com")

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			validAvailabilityReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "consultant not found", errRes.Message)
	})

	t.Run("should fail when request body is invalid JSON", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "invalidbody@example.com")

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			"{invalid json body",
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "Invalid request body", errRes.Message)
	})

	t.Run("should fail when start time format is invalid", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "badstarttime@example.com")

		// Register consultant
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

		badReq := validAvailabilityReq
		badReq.StartTime = "invalid"

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			badReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "invalid time format, expected HH:MM", errRes.Message)
	})

	t.Run("should fail when end time format is invalid", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "badendtime@example.com")

		// Register consultant
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

		badReq := validAvailabilityReq
		badReq.EndTime = "24:00"

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			badReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "invalid hour", errRes.Message)
	})

	t.Run("should fail when start time is equal to end time", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "equaltimes@example.com")

		// Register consultant
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

		badReq := validAvailabilityReq
		badReq.StartTime = "10:00"
		badReq.EndTime = "10:00"

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			badReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "invalid time range", errRes.Message)
	})

	t.Run("should fail when start time is after end time", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "reversedtimes@example.com")

		// Register consultant
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

		badReq := validAvailabilityReq
		badReq.StartTime = "12:00"
		badReq.EndTime = "10:00"

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			badReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "invalid time range", errRes.Message)
	})

	t.Run("should fail when availability overlaps with existing availability on same day", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "overlap@example.com")

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

		// 2. Create first availability (Monday 09:00 - 11:00)
		resp1, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			validAvailabilityReq,
			cookieHeader,
		)
		require.NoError(t, err)
		resp1.Body.Close()
		require.Equal(t, http.StatusCreated, resp1.StatusCode)

		// 3. Attempt to create overlapping availability (Monday 10:00 - 12:00)
		overlapReq := consultant_dto.CreateAvailabilityRequest{
			DayOfWeek: 1,
			StartTime: "10:00",
			EndTime:   "12:00",
		}

		resp2, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			overlapReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp2.Body.Close()

		require.Equal(t, http.StatusConflict, resp2.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp2)
		require.False(t, errRes.Success)
		require.Equal(t, "availability overlaps with existing availability", errRes.Message)
	})
}

package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	consultant_dto "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestActivateDeactivateAvailability_Integration(t *testing.T) {
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

	t.Run("should deactivate and activate availability slot successfully", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "activate.lifecycle@example.com")

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

		// 2. Fetch consultant ID
		getPriv, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodGet, "/test/v1/consultants/user", nil, cookieHeader)
		require.NoError(t, err)
		getPrivRes := testhelpers.DecodeResponse(t, getPriv)
		getPriv.Body.Close()
		privBytes, _ := json.Marshal(getPrivRes.Data)
		var cMap map[string]interface{}
		json.Unmarshal(privBytes, &cMap)
		consultantID := cMap["id"].(string)

		// 3. Create active availability slot
		createResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			validAvailabilityReq,
			cookieHeader,
		)
		require.NoError(t, err)
		createResp.Body.Close()

		// 4. Retrieve created slot ID
		getAvail1, err := testhelpers.PerformRequest(harness.App, http.MethodGet, "/test/v1/consultants/"+consultantID+"/availability", nil)
		require.NoError(t, err)
		getAvailRes1 := testhelpers.DecodeResponse(t, getAvail1)
		getAvail1.Body.Close()
		d1, _ := json.Marshal(getAvailRes1.Data)
		var avails1 []consultant_dto.GetAvailabilityResponse
		json.Unmarshal(d1, &avails1)
		require.Len(t, avails1, 1)
		availID := avails1[0].AvailabilityID
		require.NotEmpty(t, availID)

		// 5. Deactivate the availability slot
		deactResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPatch,
			"/test/v1/consultants/availability/"+availID+"/deactivate",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer deactResp.Body.Close()

		require.Equal(t, http.StatusOK, deactResp.StatusCode)
		deactRes := testhelpers.DecodeResponse(t, deactResp)
		require.True(t, deactRes.Success)
		require.Equal(t, "Availability deactivated successfully", deactRes.Message)

		// 6. Verify GET availability no longer includes deactivated slot (since it filters active slots)
		getAvail2, err := testhelpers.PerformRequest(harness.App, http.MethodGet, "/test/v1/consultants/"+consultantID+"/availability", nil)
		require.NoError(t, err)
		getAvailRes2 := testhelpers.DecodeResponse(t, getAvail2)
		getAvail2.Body.Close()
		d2, _ := json.Marshal(getAvailRes2.Data)
		var avails2 []consultant_dto.GetAvailabilityResponse
		json.Unmarshal(d2, &avails2)
		require.Empty(t, avails2)

		// 7. Deactivating already deactivated slot returns 409 Conflict
		deactAgainResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPatch,
			"/test/v1/consultants/availability/"+availID+"/deactivate",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer deactAgainResp.Body.Close()

		require.Equal(t, http.StatusConflict, deactAgainResp.StatusCode)
		errResDeact := testhelpers.DecodeErrorResponse(t, deactAgainResp)
		require.False(t, errResDeact.Success)
		require.Equal(t, "availability is already deactivated", errResDeact.Message)

		// 8. Reactivate the availability slot
		actResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPatch,
			"/test/v1/consultants/availability/"+availID+"/activate",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer actResp.Body.Close()

		require.Equal(t, http.StatusOK, actResp.StatusCode)
		actRes := testhelpers.DecodeResponse(t, actResp)
		require.True(t, actRes.Success)
		require.Equal(t, "Availability activated successfully", actRes.Message)

		// 9. Verify GET availability includes the reactivated slot again
		getAvail3, err := testhelpers.PerformRequest(harness.App, http.MethodGet, "/test/v1/consultants/"+consultantID+"/availability", nil)
		require.NoError(t, err)
		getAvailRes3 := testhelpers.DecodeResponse(t, getAvail3)
		getAvail3.Body.Close()
		d3, _ := json.Marshal(getAvailRes3.Data)
		var avails3 []consultant_dto.GetAvailabilityResponse
		json.Unmarshal(d3, &avails3)
		require.Len(t, avails3, 1)
		require.Equal(t, availID, avails3[0].AvailabilityID)

		// 10. Activating already active slot returns 409 Conflict
		actAgainResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPatch,
			"/test/v1/consultants/availability/"+availID+"/activate",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer actAgainResp.Body.Close()

		require.Equal(t, http.StatusConflict, actAgainResp.StatusCode)
		errResAct := testhelpers.DecodeErrorResponse(t, actAgainResp)
		require.False(t, errResAct.Success)
		require.Equal(t, "availability is already activated", errResAct.Message)
	})

	t.Run("should fail when unauthenticated (no session cookie)", func(t *testing.T) {
		harness := setUpConsultantApp(t)

		// Activate unauthenticated
		actResp, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodPatch,
			"/test/v1/consultants/availability/conav_123/activate",
			nil,
		)
		require.NoError(t, err)
		defer actResp.Body.Close()
		require.Equal(t, http.StatusUnauthorized, actResp.StatusCode)

		// Deactivate unauthenticated
		deactResp, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodPatch,
			"/test/v1/consultants/availability/conav_123/deactivate",
			nil,
		)
		require.NoError(t, err)
		defer deactResp.Body.Close()
		require.Equal(t, http.StatusUnauthorized, deactResp.StatusCode)
	})

	t.Run("should fail when session cookie is invalid", func(t *testing.T) {
		harness := setUpConsultantApp(t)

		// Activate with invalid cookie
		actResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPatch,
			"/test/v1/consultants/availability/conav_123/activate",
			nil,
			"test_session_id=invalid_token_12345678901234567890",
		)
		require.NoError(t, err)
		defer actResp.Body.Close()
		require.Equal(t, http.StatusNotFound, actResp.StatusCode)

		// Deactivate with invalid cookie
		deactResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPatch,
			"/test/v1/consultants/availability/conav_123/deactivate",
			nil,
			"test_session_id=invalid_token_12345678901234567890",
		)
		require.NoError(t, err)
		defer deactResp.Body.Close()
		require.Equal(t, http.StatusNotFound, deactResp.StatusCode)
	})

	t.Run("should fail when consultant profile is not found for authenticated user", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "noconsultant.act@example.com")

		// Activate without consultant profile
		actResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPatch,
			"/test/v1/consultants/availability/conav_123/activate",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer actResp.Body.Close()
		require.Equal(t, http.StatusNotFound, actResp.StatusCode)
		errResAct := testhelpers.DecodeErrorResponse(t, actResp)
		require.Equal(t, "consultant not found", errResAct.Message)

		// Deactivate without consultant profile
		deactResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPatch,
			"/test/v1/consultants/availability/conav_123/deactivate",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer deactResp.Body.Close()
		require.Equal(t, http.StatusNotFound, deactResp.StatusCode)
		errResDeact := testhelpers.DecodeErrorResponse(t, deactResp)
		require.Equal(t, "consultant not found", errResDeact.Message)
	})

	t.Run("should fail when availability_id does not exist", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "nonexistent.act@example.com")

		// Register consultant
		regResp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/register", validConsultantReq, cookieHeader)
		require.NoError(t, err)
		regResp.Body.Close()

		// Activate non-existent slot
		actResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPatch,
			"/test/v1/consultants/availability/conav_nonexistent_999/activate",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer actResp.Body.Close()
		require.Equal(t, http.StatusNotFound, actResp.StatusCode)
		errResAct := testhelpers.DecodeErrorResponse(t, actResp)
		require.Equal(t, "availability not found", errResAct.Message)

		// Deactivate non-existent slot
		deactResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPatch,
			"/test/v1/consultants/availability/conav_nonexistent_999/deactivate",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer deactResp.Body.Close()
		require.Equal(t, http.StatusNotFound, deactResp.StatusCode)
		errResDeact := testhelpers.DecodeErrorResponse(t, deactResp)
		require.Equal(t, "availability not found", errResDeact.Message)
	})

	t.Run("should fail when activating or deactivating slot belonging to another consultant", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookie1 := registerAndLoginUser(t, harness, "consultant1.act@example.com")
		cookie2 := registerAndLoginUser(t, harness, "consultant2.act@example.com")

		// 1. Register Consultant 1
		reg1, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/register", validConsultantReq, cookie1)
		require.NoError(t, err)
		reg1.Body.Close()

		// 2. Register Consultant 2
		req2 := validConsultantReq
		req2.DisplayName = "Consultant Beta"
		reg2, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/register", req2, cookie2)
		require.NoError(t, err)
		reg2.Body.Close()

		// 3. Get Consultant 1 ID
		getPriv1, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodGet, "/test/v1/consultants/user", nil, cookie1)
		require.NoError(t, err)
		getPrivRes1 := testhelpers.DecodeResponse(t, getPriv1)
		getPriv1.Body.Close()
		b1, _ := json.Marshal(getPrivRes1.Data)
		var cMap1 map[string]interface{}
		json.Unmarshal(b1, &cMap1)
		cID1 := cMap1["id"].(string)

		// 4. Create slot for Consultant 1
		createResp1, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/availability", validAvailabilityReq, cookie1)
		require.NoError(t, err)
		createResp1.Body.Close()

		// 5. Get slot ID for Consultant 1
		getAvail1, err := testhelpers.PerformRequest(harness.App, http.MethodGet, "/test/v1/consultants/"+cID1+"/availability", nil)
		require.NoError(t, err)
		getAvailRes1 := testhelpers.DecodeResponse(t, getAvail1)
		getAvail1.Body.Close()
		d1, _ := json.Marshal(getAvailRes1.Data)
		var avails1 []consultant_dto.GetAvailabilityResponse
		json.Unmarshal(d1, &avails1)
		slotID1 := avails1[0].AvailabilityID

		// 6. Consultant 2 attempts to deactivate Consultant 1's slot
		deactResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPatch,
			"/test/v1/consultants/availability/"+slotID1+"/deactivate",
			nil,
			cookie2,
		)
		require.NoError(t, err)
		defer deactResp.Body.Close()
		require.Equal(t, http.StatusNotFound, deactResp.StatusCode)
		errResDeact := testhelpers.DecodeErrorResponse(t, deactResp)
		require.Equal(t, "availability not found", errResDeact.Message)

		// 7. Consultant 1 deactivates own slot
		deactOwnResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPatch,
			"/test/v1/consultants/availability/"+slotID1+"/deactivate",
			nil,
			cookie1,
		)
		require.NoError(t, err)
		deactOwnResp.Body.Close()
		require.Equal(t, http.StatusOK, deactOwnResp.StatusCode)

		// 8. Consultant 2 attempts to activate Consultant 1's inactive slot
		actResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPatch,
			"/test/v1/consultants/availability/"+slotID1+"/activate",
			nil,
			cookie2,
		)
		require.NoError(t, err)
		defer actResp.Body.Close()
		require.Equal(t, http.StatusNotFound, actResp.StatusCode)
		errResAct := testhelpers.DecodeErrorResponse(t, actResp)
		require.Equal(t, "availability not found", errResAct.Message)
	})
}

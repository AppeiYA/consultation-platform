package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	consultant_dto "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestUpdateAvailability_Integration(t *testing.T) {
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

	t.Run("should update availability successfully for registered consultant", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "update.avail.success@example.com")

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

		// 2. Fetch consultant ID
		getPrivateResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/consultants/user",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		getPrivateRes := testhelpers.DecodeResponse(t, getPrivateResp)
		getPrivateResp.Body.Close()
		privBytes, err := json.Marshal(getPrivateRes.Data)
		require.NoError(t, err)
		var consultantMap map[string]interface{}
		require.NoError(t, json.Unmarshal(privBytes, &consultantMap))
		consultantID := consultantMap["id"].(string)

		// 3. Create initial availability slot (Monday 09:00 - 11:00)
		createResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			validAvailabilityReq,
			cookieHeader,
		)
		require.NoError(t, err)
		createResp.Body.Close()
		require.Equal(t, http.StatusCreated, createResp.StatusCode)

		// 4. Retrieve created availability to get its ID
		getAvailResp, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodGet,
			"/test/v1/consultants/"+consultantID+"/availability",
			nil,
		)
		require.NoError(t, err)
		getAvailRes := testhelpers.DecodeResponse(t, getAvailResp)
		getAvailResp.Body.Close()
		require.True(t, getAvailRes.Success)

		dataBytes, err := json.Marshal(getAvailRes.Data)
		require.NoError(t, err)
		var avails []consultant_dto.GetAvailabilityResponse
		require.NoError(t, json.Unmarshal(dataBytes, &avails))
		require.Len(t, avails, 1)
		availID := avails[0].AvailabilityID
		require.NotEmpty(t, availID)

		// 5. Update availability to Tuesday 14:00 - 16:00
		updateReq := consultant_dto.UpdateAvailabilityRequest{
			AvailabilityID: availID,
			DayOfWeek:      2, // Tuesday
			StartTime:      "14:00",
			EndTime:        "16:00",
		}

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPut,
			"/test/v1/consultants/availability",
			updateReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		res := testhelpers.DecodeResponse(t, resp)
		require.True(t, res.Success)
		require.Equal(t, "Availability updated successfully", res.Message)

		respDataBytes, err := json.Marshal(res.Data)
		require.NoError(t, err)
		var updatedSlot consultant_dto.GetAvailabilityResponse
		require.NoError(t, json.Unmarshal(respDataBytes, &updatedSlot))

		require.Equal(t, availID, updatedSlot.AvailabilityID)
		require.Equal(t, "Tuesday", updatedSlot.DayOfWeek)
		require.Equal(t, "14:00", updatedSlot.StartTime)
		require.Equal(t, "16:00", updatedSlot.EndTime)

		// 6. Verify via GET availability endpoint that changes are persisted
		getAvailResp2, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodGet,
			"/test/v1/consultants/"+consultantID+"/availability",
			nil,
		)
		require.NoError(t, err)
		getAvailRes2 := testhelpers.DecodeResponse(t, getAvailResp2)
		getAvailResp2.Body.Close()

		dataBytes2, err := json.Marshal(getAvailRes2.Data)
		require.NoError(t, err)
		var avails2 []consultant_dto.GetAvailabilityResponse
		require.NoError(t, json.Unmarshal(dataBytes2, &avails2))
		require.Len(t, avails2, 1)
		require.Equal(t, availID, avails2[0].AvailabilityID)
		require.Equal(t, "Tuesday", avails2[0].DayOfWeek)
		require.Equal(t, "14:00", avails2[0].StartTime)
		require.Equal(t, "16:00", avails2[0].EndTime)
	})

	t.Run("should update availability on the same day without changing day of week", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "update.sameday@example.com")

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
		getPrivateResp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodGet, "/test/v1/consultants/user", nil, cookieHeader)
		require.NoError(t, err)
		getPrivateRes := testhelpers.DecodeResponse(t, getPrivateResp)
		getPrivateResp.Body.Close()
		privBytes, _ := json.Marshal(getPrivateRes.Data)
		var consultantMap map[string]interface{}
		json.Unmarshal(privBytes, &consultantMap)
		consultantID := consultantMap["id"].(string)

		// 3. Create initial slot (Monday 09:00 - 11:00)
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
		getAvailResp, err := testhelpers.PerformRequest(harness.App, http.MethodGet, "/test/v1/consultants/"+consultantID+"/availability", nil)
		require.NoError(t, err)
		getAvailRes := testhelpers.DecodeResponse(t, getAvailResp)
		getAvailResp.Body.Close()
		dataBytes, _ := json.Marshal(getAvailRes.Data)
		var avails []consultant_dto.GetAvailabilityResponse
		json.Unmarshal(dataBytes, &avails)
		availID := avails[0].AvailabilityID

		// 5. Update slot on same Monday to 10:00 - 12:00
		updateReq := consultant_dto.UpdateAvailabilityRequest{
			AvailabilityID: availID,
			DayOfWeek:      1, // Monday
			StartTime:      "10:00",
			EndTime:        "12:00",
		}

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPut,
			"/test/v1/consultants/availability",
			updateReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		res := testhelpers.DecodeResponse(t, resp)
		require.True(t, res.Success)

		respDataBytes, _ := json.Marshal(res.Data)
		var updatedSlot consultant_dto.GetAvailabilityResponse
		json.Unmarshal(respDataBytes, &updatedSlot)
		require.Equal(t, "Monday", updatedSlot.DayOfWeek)
		require.Equal(t, "10:00", updatedSlot.StartTime)
		require.Equal(t, "12:00", updatedSlot.EndTime)
	})

	t.Run("should update slot to adjacent position next to another existing slot", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "update.adjacent@example.com")

		// 1. Register consultant
		regResp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/register", validConsultantReq, cookieHeader)
		require.NoError(t, err)
		regResp.Body.Close()

		// 2. Fetch consultant ID
		getPriv, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodGet, "/test/v1/consultants/user", nil, cookieHeader)
		require.NoError(t, err)
		getPrivRes := testhelpers.DecodeResponse(t, getPriv)
		getPriv.Body.Close()
		b, _ := json.Marshal(getPrivRes.Data)
		var cMap map[string]interface{}
		json.Unmarshal(b, &cMap)
		consultantID := cMap["id"].(string)

		// 3. Create Slot 1 (Monday 08:00 - 10:00) and Slot 2 (Monday 14:00 - 16:00)
		slot1Req := consultant_dto.CreateAvailabilityRequest{DayOfWeek: 1, StartTime: "08:00", EndTime: "10:00"}
		r1, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/availability", slot1Req, cookieHeader)
		require.NoError(t, err)
		r1.Body.Close()

		slot2Req := consultant_dto.CreateAvailabilityRequest{DayOfWeek: 1, StartTime: "14:00", EndTime: "16:00"}
		r2, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/availability", slot2Req, cookieHeader)
		require.NoError(t, err)
		r2.Body.Close()

		// 4. Get slot IDs
		getAvail, err := testhelpers.PerformRequest(harness.App, http.MethodGet, "/test/v1/consultants/"+consultantID+"/availability", nil)
		require.NoError(t, err)
		getAvailRes := testhelpers.DecodeResponse(t, getAvail)
		getAvail.Body.Close()
		dBytes, _ := json.Marshal(getAvailRes.Data)
		var avails []consultant_dto.GetAvailabilityResponse
		json.Unmarshal(dBytes, &avails)
		require.Len(t, avails, 2)

		// Identify slot 2 ID (the one at 14:00)
		var slot2ID string
		for _, a := range avails {
			if a.StartTime == "14:00" {
				slot2ID = a.AvailabilityID
			}
		}
		require.NotEmpty(t, slot2ID)

		// 5. Update slot 2 to Monday 10:00 - 12:00 (adjacent right after slot 1)
		updateReq := consultant_dto.UpdateAvailabilityRequest{
			AvailabilityID: slot2ID,
			DayOfWeek:      1,
			StartTime:      "10:00",
			EndTime:        "12:00",
		}

		resp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPut, "/test/v1/consultants/availability", updateReq, cookieHeader)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		res := testhelpers.DecodeResponse(t, resp)
		require.True(t, res.Success)
	})

	t.Run("should fail when unauthenticated (no session cookie)", func(t *testing.T) {
		harness := setUpConsultantApp(t)

		updateReq := consultant_dto.UpdateAvailabilityRequest{
			AvailabilityID: "conav_nonexistent",
			DayOfWeek:      1,
			StartTime:      "09:00",
			EndTime:        "11:00",
		}

		resp, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodPut,
			"/test/v1/consultants/availability",
			updateReq,
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

		updateReq := consultant_dto.UpdateAvailabilityRequest{
			AvailabilityID: "conav_nonexistent",
			DayOfWeek:      1,
			StartTime:      "09:00",
			EndTime:        "11:00",
		}

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPut,
			"/test/v1/consultants/availability",
			updateReq,
			"test_session_id=invalid_session_token_12345678901234567890",
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
		cookieHeader := registerAndLoginUser(t, harness, "noconsultant.update@example.com")

		updateReq := consultant_dto.UpdateAvailabilityRequest{
			AvailabilityID: "conav_nonexistent",
			DayOfWeek:      1,
			StartTime:      "09:00",
			EndTime:        "11:00",
		}

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPut,
			"/test/v1/consultants/availability",
			updateReq,
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
		cookieHeader := registerAndLoginUser(t, harness, "invalidbody.update@example.com")

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPut,
			"/test/v1/consultants/availability",
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

	t.Run("should fail when availability_id is empty", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "emptyid.update@example.com")

		// Register consultant
		regResp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/register", validConsultantReq, cookieHeader)
		require.NoError(t, err)
		regResp.Body.Close()

		updateReq := consultant_dto.UpdateAvailabilityRequest{
			AvailabilityID: "",
			DayOfWeek:      1,
			StartTime:      "09:00",
			EndTime:        "11:00",
		}

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPut,
			"/test/v1/consultants/availability",
			updateReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "availability_id is required", errRes.Message)
	})

	t.Run("should fail when day of week is invalid", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "badday.update@example.com")

		// Register consultant
		regResp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/register", validConsultantReq, cookieHeader)
		require.NoError(t, err)
		regResp.Body.Close()

		updateReq := consultant_dto.UpdateAvailabilityRequest{
			AvailabilityID: "conav_123",
			DayOfWeek:      7, // Invalid day of week (must be 0-6)
			StartTime:      "09:00",
			EndTime:        "11:00",
		}

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPut,
			"/test/v1/consultants/availability",
			updateReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "Invalid day of the week", errRes.Message)
	})

	t.Run("should fail when start time format is invalid", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "badstart.update@example.com")

		// 1. Register consultant
		regResp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/register", validConsultantReq, cookieHeader)
		require.NoError(t, err)
		regResp.Body.Close()

		// 2. Fetch consultant ID
		getPriv, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodGet, "/test/v1/consultants/user", nil, cookieHeader)
		require.NoError(t, err)
		getPrivRes := testhelpers.DecodeResponse(t, getPriv)
		getPriv.Body.Close()
		b, _ := json.Marshal(getPrivRes.Data)
		var cMap map[string]interface{}
		json.Unmarshal(b, &cMap)
		consultantID := cMap["id"].(string)

		// 3. Create slot
		createResp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/availability", validAvailabilityReq, cookieHeader)
		require.NoError(t, err)
		createResp.Body.Close()

		// 4. Get slot ID
		getAvail, err := testhelpers.PerformRequest(harness.App, http.MethodGet, "/test/v1/consultants/"+consultantID+"/availability", nil)
		require.NoError(t, err)
		getAvailRes := testhelpers.DecodeResponse(t, getAvail)
		getAvail.Body.Close()
		dBytes, _ := json.Marshal(getAvailRes.Data)
		var avails []consultant_dto.GetAvailabilityResponse
		json.Unmarshal(dBytes, &avails)
		availID := avails[0].AvailabilityID

		// 5. Update with invalid start time
		badReq := consultant_dto.UpdateAvailabilityRequest{
			AvailabilityID: availID,
			DayOfWeek:      1,
			StartTime:      "invalid_time",
			EndTime:        "11:00",
		}

		resp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPut, "/test/v1/consultants/availability", badReq, cookieHeader)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "invalid time format, expected HH:MM", errRes.Message)
	})

	t.Run("should fail when end time hour is invalid", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "badend.update@example.com")

		// 1. Register consultant
		regResp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/register", validConsultantReq, cookieHeader)
		require.NoError(t, err)
		regResp.Body.Close()

		// 2. Fetch consultant ID
		getPriv, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodGet, "/test/v1/consultants/user", nil, cookieHeader)
		require.NoError(t, err)
		getPrivRes := testhelpers.DecodeResponse(t, getPriv)
		getPriv.Body.Close()
		b, _ := json.Marshal(getPrivRes.Data)
		var cMap map[string]interface{}
		json.Unmarshal(b, &cMap)
		consultantID := cMap["id"].(string)

		// 3. Create slot
		createResp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/availability", validAvailabilityReq, cookieHeader)
		require.NoError(t, err)
		createResp.Body.Close()

		// 4. Get slot ID
		getAvail, err := testhelpers.PerformRequest(harness.App, http.MethodGet, "/test/v1/consultants/"+consultantID+"/availability", nil)
		require.NoError(t, err)
		getAvailRes := testhelpers.DecodeResponse(t, getAvail)
		getAvail.Body.Close()
		dBytes, _ := json.Marshal(getAvailRes.Data)
		var avails []consultant_dto.GetAvailabilityResponse
		json.Unmarshal(dBytes, &avails)
		availID := avails[0].AvailabilityID

		// 5. Update with invalid end time hour (24:00)
		badReq := consultant_dto.UpdateAvailabilityRequest{
			AvailabilityID: availID,
			DayOfWeek:      1,
			StartTime:      "09:00",
			EndTime:        "24:00",
		}

		resp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPut, "/test/v1/consultants/availability", badReq, cookieHeader)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "invalid hour", errRes.Message)
	})

	t.Run("should fail when start time is equal to end time", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "equaltime.update@example.com")

		// 1. Register consultant
		regResp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/register", validConsultantReq, cookieHeader)
		require.NoError(t, err)
		regResp.Body.Close()

		// 2. Fetch consultant ID
		getPriv, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodGet, "/test/v1/consultants/user", nil, cookieHeader)
		require.NoError(t, err)
		getPrivRes := testhelpers.DecodeResponse(t, getPriv)
		getPriv.Body.Close()
		b, _ := json.Marshal(getPrivRes.Data)
		var cMap map[string]interface{}
		json.Unmarshal(b, &cMap)
		consultantID := cMap["id"].(string)

		// 3. Create slot
		createResp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/availability", validAvailabilityReq, cookieHeader)
		require.NoError(t, err)
		createResp.Body.Close()

		// 4. Get slot ID
		getAvail, err := testhelpers.PerformRequest(harness.App, http.MethodGet, "/test/v1/consultants/"+consultantID+"/availability", nil)
		require.NoError(t, err)
		getAvailRes := testhelpers.DecodeResponse(t, getAvail)
		getAvail.Body.Close()
		dBytes, _ := json.Marshal(getAvailRes.Data)
		var avails []consultant_dto.GetAvailabilityResponse
		json.Unmarshal(dBytes, &avails)
		availID := avails[0].AvailabilityID

		// 5. Update with equal start and end times
		badReq := consultant_dto.UpdateAvailabilityRequest{
			AvailabilityID: availID,
			DayOfWeek:      1,
			StartTime:      "10:00",
			EndTime:        "10:00",
		}

		resp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPut, "/test/v1/consultants/availability", badReq, cookieHeader)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "invalid time range", errRes.Message)
	})

	t.Run("should fail when start time is after end time", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "reversedtime.update@example.com")

		// 1. Register consultant
		regResp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/register", validConsultantReq, cookieHeader)
		require.NoError(t, err)
		regResp.Body.Close()

		// 2. Fetch consultant ID
		getPriv, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodGet, "/test/v1/consultants/user", nil, cookieHeader)
		require.NoError(t, err)
		getPrivRes := testhelpers.DecodeResponse(t, getPriv)
		getPriv.Body.Close()
		b, _ := json.Marshal(getPrivRes.Data)
		var cMap map[string]interface{}
		json.Unmarshal(b, &cMap)
		consultantID := cMap["id"].(string)

		// 3. Create slot
		createResp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/availability", validAvailabilityReq, cookieHeader)
		require.NoError(t, err)
		createResp.Body.Close()

		// 4. Get slot ID
		getAvail, err := testhelpers.PerformRequest(harness.App, http.MethodGet, "/test/v1/consultants/"+consultantID+"/availability", nil)
		require.NoError(t, err)
		getAvailRes := testhelpers.DecodeResponse(t, getAvail)
		getAvail.Body.Close()
		dBytes, _ := json.Marshal(getAvailRes.Data)
		var avails []consultant_dto.GetAvailabilityResponse
		json.Unmarshal(dBytes, &avails)
		availID := avails[0].AvailabilityID

		// 5. Update with start time after end time
		badReq := consultant_dto.UpdateAvailabilityRequest{
			AvailabilityID: availID,
			DayOfWeek:      1,
			StartTime:      "15:00",
			EndTime:        "11:00",
		}

		resp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPut, "/test/v1/consultants/availability", badReq, cookieHeader)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "invalid time range", errRes.Message)
	})

	t.Run("should fail when availability_id does not exist", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "nonexistent.update@example.com")

		// Register consultant
		regResp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/register", validConsultantReq, cookieHeader)
		require.NoError(t, err)
		regResp.Body.Close()

		updateReq := consultant_dto.UpdateAvailabilityRequest{
			AvailabilityID: "conav_nonexistent_id_999",
			DayOfWeek:      1,
			StartTime:      "09:00",
			EndTime:        "11:00",
		}

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPut,
			"/test/v1/consultants/availability",
			updateReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "availability not found", errRes.Message)
	})

	t.Run("should fail when updating an availability belonging to another consultant", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookie1 := registerAndLoginUser(t, harness, "consultantA@example.com")
		cookie2 := registerAndLoginUser(t, harness, "consultantB@example.com")

		// Register consultant 1
		reg1, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/register", validConsultantReq, cookie1)
		require.NoError(t, err)
		reg1.Body.Close()

		// Register consultant 2
		req2 := validConsultantReq
		req2.DisplayName = "Consultant B"
		reg2, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/register", req2, cookie2)
		require.NoError(t, err)
		reg2.Body.Close()

		// Get consultant 1 ID
		getPriv1, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodGet, "/test/v1/consultants/user", nil, cookie1)
		require.NoError(t, err)
		getPrivRes1 := testhelpers.DecodeResponse(t, getPriv1)
		getPriv1.Body.Close()
		b1, _ := json.Marshal(getPrivRes1.Data)
		var cMap1 map[string]interface{}
		json.Unmarshal(b1, &cMap1)
		cID1 := cMap1["id"].(string)

		// Create slot for consultant 1
		createResp1, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/availability", validAvailabilityReq, cookie1)
		require.NoError(t, err)
		createResp1.Body.Close()

		// Get consultant 1's availability slot ID
		getAvail1, err := testhelpers.PerformRequest(harness.App, http.MethodGet, "/test/v1/consultants/"+cID1+"/availability", nil)
		require.NoError(t, err)
		getAvailRes1 := testhelpers.DecodeResponse(t, getAvail1)
		getAvail1.Body.Close()
		d1, _ := json.Marshal(getAvailRes1.Data)
		var avails1 []consultant_dto.GetAvailabilityResponse
		json.Unmarshal(d1, &avails1)
		slotID1 := avails1[0].AvailabilityID

		// Consultant 2 tries to update Consultant 1's availability slot
		badReq := consultant_dto.UpdateAvailabilityRequest{
			AvailabilityID: slotID1,
			DayOfWeek:      2,
			StartTime:      "14:00",
			EndTime:        "16:00",
		}

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPut,
			"/test/v1/consultants/availability",
			badReq,
			cookie2,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "availability not found", errRes.Message)
	})

	t.Run("should fail when updated availability overlaps with existing availability on same day", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "update.overlap@example.com")

		// 1. Register consultant
		regResp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/register", validConsultantReq, cookieHeader)
		require.NoError(t, err)
		regResp.Body.Close()

		// 2. Fetch consultant ID
		getPriv, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodGet, "/test/v1/consultants/user", nil, cookieHeader)
		require.NoError(t, err)
		getPrivRes := testhelpers.DecodeResponse(t, getPriv)
		getPriv.Body.Close()
		b, _ := json.Marshal(getPrivRes.Data)
		var cMap map[string]interface{}
		json.Unmarshal(b, &cMap)
		consultantID := cMap["id"].(string)

		// 3. Create Slot 1: Monday 09:00 - 11:00
		slot1Req := consultant_dto.CreateAvailabilityRequest{DayOfWeek: 1, StartTime: "09:00", EndTime: "11:00"}
		r1, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/availability", slot1Req, cookieHeader)
		require.NoError(t, err)
		r1.Body.Close()

		// 4. Create Slot 2: Monday 14:00 - 16:00
		slot2Req := consultant_dto.CreateAvailabilityRequest{DayOfWeek: 1, StartTime: "14:00", EndTime: "16:00"}
		r2, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultants/availability", slot2Req, cookieHeader)
		require.NoError(t, err)
		r2.Body.Close()

		// 5. Get slot IDs
		getAvail, err := testhelpers.PerformRequest(harness.App, http.MethodGet, "/test/v1/consultants/"+consultantID+"/availability", nil)
		require.NoError(t, err)
		getAvailRes := testhelpers.DecodeResponse(t, getAvail)
		getAvail.Body.Close()
		dBytes, _ := json.Marshal(getAvailRes.Data)
		var avails []consultant_dto.GetAvailabilityResponse
		json.Unmarshal(dBytes, &avails)
		require.Len(t, avails, 2)

		var slot2ID string
		for _, a := range avails {
			if a.StartTime == "14:00" {
				slot2ID = a.AvailabilityID
			}
		}
		require.NotEmpty(t, slot2ID)

		// 6. Attempt to update Slot 2 to Monday 10:00 - 12:00 (overlaps with Slot 1 09:00 - 11:00)
		overlapUpdateReq := consultant_dto.UpdateAvailabilityRequest{
			AvailabilityID: slot2ID,
			DayOfWeek:      1,
			StartTime:      "10:00",
			EndTime:        "12:00",
		}

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPut,
			"/test/v1/consultants/availability",
			overlapUpdateReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusConflict, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "availability overlaps with existing availability", errRes.Message)
	})
}

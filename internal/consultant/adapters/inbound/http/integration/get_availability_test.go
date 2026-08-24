package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	consultant_dto "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestGetAvailability_Integration(t *testing.T) {
	validConsultantReq := consultant_dto.RegisterConsultantDTO{
		ProfessionID:    "prof_9ee432d7-b672-40ae-b03f-c1f1fb696621",
		DisplayName:     "Jane Doe Tech",
		Bio:             "Experienced software engineer with 10 years of experience.",
		YearsExperience: 10,
	}

	t.Run("should get availabilities successfully without authentication for registered consultant", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "avail.get.public@example.com")

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

		// 2. Fetch private profile to obtain consultant ID
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

		privateBytes, err := json.Marshal(getPrivateRes.Data)
		require.NoError(t, err)
		var consultantMap map[string]interface{}
		err = json.Unmarshal(privateBytes, &consultantMap)
		require.NoError(t, err)
		consultantID := consultantMap["id"].(string)

		// 3. Create two availability slots
		slot1 := consultant_dto.CreateAvailabilityRequest{
			DayOfWeek: 1, // Monday
			StartTime: "09:00",
			EndTime:   "11:00",
		}
		resp1, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			slot1,
			cookieHeader,
		)
		require.NoError(t, err)
		resp1.Body.Close()
		require.Equal(t, http.StatusCreated, resp1.StatusCode)

		slot2 := consultant_dto.CreateAvailabilityRequest{
			DayOfWeek: 3, // Wednesday
			StartTime: "14:00",
			EndTime:   "16:00",
		}
		resp2, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			slot2,
			cookieHeader,
		)
		require.NoError(t, err)
		resp2.Body.Close()
		require.Equal(t, http.StatusCreated, resp2.StatusCode)

		// 4. Public endpoint GET /test/v1/consultants/:consultantID/availability (no cookie)
		resp, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodGet,
			"/test/v1/consultants/"+consultantID+"/availability",
			nil,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		res := testhelpers.DecodeResponse(t, resp)
		require.True(t, res.Success)
		require.Equal(t, "Availabilities fetched successfully", res.Message)

		dataBytes, err := json.Marshal(res.Data)
		require.NoError(t, err)

		var availabilities []consultant_dto.GetAvailabilityResponse
		err = json.Unmarshal(dataBytes, &availabilities)
		require.NoError(t, err)

		require.Len(t, availabilities, 2)
		require.Equal(t, "Monday", availabilities[0].DayOfWeek)
		require.Equal(t, "09:00", availabilities[0].StartTime)
		require.Equal(t, "11:00", availabilities[0].EndTime)

		require.Equal(t, "Wednesday", availabilities[1].DayOfWeek)
		require.Equal(t, "14:00", availabilities[1].StartTime)
		require.Equal(t, "16:00", availabilities[1].EndTime)
	})

	t.Run("should get availabilities successfully when requested by authenticated user", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		consultantCookie := registerAndLoginUser(t, harness, "avail.consultant@example.com")
		clientCookie := registerAndLoginUser(t, harness, "avail.client@example.com")

		// 1. Register consultant
		regResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			validConsultantReq,
			consultantCookie,
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
			consultantCookie,
		)
		require.NoError(t, err)
		getPrivateRes := testhelpers.DecodeResponse(t, getPrivateResp)
		getPrivateResp.Body.Close()

		privateBytes, err := json.Marshal(getPrivateRes.Data)
		require.NoError(t, err)
		var consultantMap map[string]interface{}
		err = json.Unmarshal(privateBytes, &consultantMap)
		require.NoError(t, err)
		consultantID := consultantMap["id"].(string)

		// 3. Create availability
		respCreate, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			consultant_dto.CreateAvailabilityRequest{
				DayOfWeek: 2,
				StartTime: "10:00",
				EndTime:   "12:00",
			},
			consultantCookie,
		)
		require.NoError(t, err)
		respCreate.Body.Close()
		require.Equal(t, http.StatusCreated, respCreate.StatusCode)

		// 4. Client user fetches consultant's availability using client session cookie
		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/consultants/"+consultantID+"/availability",
			nil,
			clientCookie,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		res := testhelpers.DecodeResponse(t, resp)
		require.True(t, res.Success)
		require.Equal(t, "Availabilities fetched successfully", res.Message)

		dataBytes, err := json.Marshal(res.Data)
		require.NoError(t, err)

		var availabilities []consultant_dto.GetAvailabilityResponse
		err = json.Unmarshal(dataBytes, &availabilities)
		require.NoError(t, err)

		require.Len(t, availabilities, 1)
		require.Equal(t, "Tuesday", availabilities[0].DayOfWeek)
		require.Equal(t, "10:00", availabilities[0].StartTime)
		require.Equal(t, "12:00", availabilities[0].EndTime)
	})

	t.Run("should return empty list when consultant has no availabilities", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "avail.empty@example.com")

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

		privateBytes, err := json.Marshal(getPrivateRes.Data)
		require.NoError(t, err)
		var consultantMap map[string]interface{}
		err = json.Unmarshal(privateBytes, &consultantMap)
		require.NoError(t, err)
		consultantID := consultantMap["id"].(string)

		// 3. Fetch availability without creating any slots
		resp, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodGet,
			"/test/v1/consultants/"+consultantID+"/availability",
			nil,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		res := testhelpers.DecodeResponse(t, resp)
		require.True(t, res.Success)
		require.Equal(t, "Availabilities fetched successfully", res.Message)

		dataBytes, err := json.Marshal(res.Data)
		require.NoError(t, err)

		var availabilities []consultant_dto.GetAvailabilityResponse
		err = json.Unmarshal(dataBytes, &availabilities)
		require.NoError(t, err)

		require.Empty(t, availabilities)
	})

	t.Run("should return empty list for non-existent consultant ID", func(t *testing.T) {
		harness := setUpConsultantApp(t)

		resp, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodGet,
			"/test/v1/consultants/non_existent_consultant_id/availability",
			nil,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		res := testhelpers.DecodeResponse(t, resp)
		require.True(t, res.Success)
		require.Equal(t, "Availabilities fetched successfully", res.Message)

		dataBytes, err := json.Marshal(res.Data)
		require.NoError(t, err)

		var availabilities []consultant_dto.GetAvailabilityResponse
		err = json.Unmarshal(dataBytes, &availabilities)
		require.NoError(t, err)

		require.Empty(t, availabilities)
	})

	t.Run("should return multiple non-overlapping slots on same day and different days", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "avail.multislots@example.com")

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

		privateBytes, err := json.Marshal(getPrivateRes.Data)
		require.NoError(t, err)
		var consultantMap map[string]interface{}
		err = json.Unmarshal(privateBytes, &consultantMap)
		require.NoError(t, err)
		consultantID := consultantMap["id"].(string)

		// 3. Create slots: Monday morning, Monday afternoon, Friday morning
		slots := []consultant_dto.CreateAvailabilityRequest{
			{DayOfWeek: 1, StartTime: "08:00", EndTime: "11:00"},
			{DayOfWeek: 1, StartTime: "13:00", EndTime: "17:00"},
			{DayOfWeek: 5, StartTime: "09:00", EndTime: "12:00"},
		}
		for _, s := range slots {
			r, err := testhelpers.PerformRequestWithCookie(
				harness.App,
				http.MethodPost,
				"/test/v1/consultants/availability",
				s,
				cookieHeader,
			)
			require.NoError(t, err)
			r.Body.Close()
			require.Equal(t, http.StatusCreated, r.StatusCode)
		}

		// 4. Fetch availability
		resp, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodGet,
			"/test/v1/consultants/"+consultantID+"/availability",
			nil,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		res := testhelpers.DecodeResponse(t, resp)
		require.True(t, res.Success)

		dataBytes, err := json.Marshal(res.Data)
		require.NoError(t, err)

		var availabilities []consultant_dto.GetAvailabilityResponse
		err = json.Unmarshal(dataBytes, &availabilities)
		require.NoError(t, err)

		require.Len(t, availabilities, 3)
	})

	t.Run("should isolate availabilities between different consultants", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookie1 := registerAndLoginUser(t, harness, "consultant1@example.com")
		cookie2 := registerAndLoginUser(t, harness, "consultant2@example.com")

		// Register consultant 1
		req1 := validConsultantReq
		req1.DisplayName = "Consultant One"
		regResp1, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			req1,
			cookie1,
		)
		require.NoError(t, err)
		regResp1.Body.Close()
		require.Equal(t, http.StatusCreated, regResp1.StatusCode)

		// Register consultant 2
		req2 := validConsultantReq
		req2.DisplayName = "Consultant Two"
		regResp2, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			req2,
			cookie2,
		)
		require.NoError(t, err)
		regResp2.Body.Close()
		require.Equal(t, http.StatusCreated, regResp2.StatusCode)

		// Get IDs
		getPriv1, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodGet, "/test/v1/consultants/user", nil, cookie1)
		require.NoError(t, err)
		res1 := testhelpers.DecodeResponse(t, getPriv1)
		getPriv1.Body.Close()
		b1, _ := json.Marshal(res1.Data)
		var m1 map[string]interface{}
		json.Unmarshal(b1, &m1)
		id1 := m1["id"].(string)

		getPriv2, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodGet, "/test/v1/consultants/user", nil, cookie2)
		require.NoError(t, err)
		res2 := testhelpers.DecodeResponse(t, getPriv2)
		getPriv2.Body.Close()
		b2, _ := json.Marshal(res2.Data)
		var m2 map[string]interface{}
		json.Unmarshal(b2, &m2)
		id2 := m2["id"].(string)

		// Create slot for consultant 1 (Monday 09:00-11:00)
		r1, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			consultant_dto.CreateAvailabilityRequest{DayOfWeek: 1, StartTime: "09:00", EndTime: "11:00"},
			cookie1,
		)
		require.NoError(t, err)
		r1.Body.Close()
		require.Equal(t, http.StatusCreated, r1.StatusCode)

		// Create slot for consultant 2 (Tuesday 14:00-16:00)
		r2, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/availability",
			consultant_dto.CreateAvailabilityRequest{DayOfWeek: 2, StartTime: "14:00", EndTime: "16:00"},
			cookie2,
		)
		require.NoError(t, err)
		r2.Body.Close()
		require.Equal(t, http.StatusCreated, r2.StatusCode)

		// Fetch availability for consultant 1
		respA, err := testhelpers.PerformRequest(harness.App, http.MethodGet, "/test/v1/consultants/"+id1+"/availability", nil)
		require.NoError(t, err)
		defer respA.Body.Close()
		require.Equal(t, http.StatusOK, respA.StatusCode)
		resA := testhelpers.DecodeResponse(t, respA)
		bA, _ := json.Marshal(resA.Data)
		var availsA []consultant_dto.GetAvailabilityResponse
		json.Unmarshal(bA, &availsA)
		require.Len(t, availsA, 1)
		require.Equal(t, "Monday", availsA[0].DayOfWeek)
		require.Equal(t, "09:00", availsA[0].StartTime)

		// Fetch availability for consultant 2
		respB, err := testhelpers.PerformRequest(harness.App, http.MethodGet, "/test/v1/consultants/"+id2+"/availability", nil)
		require.NoError(t, err)
		defer respB.Body.Close()
		require.Equal(t, http.StatusOK, respB.StatusCode)
		resB := testhelpers.DecodeResponse(t, respB)
		bB, _ := json.Marshal(resB.Data)
		var availsB []consultant_dto.GetAvailabilityResponse
		json.Unmarshal(bB, &availsB)
		require.Len(t, availsB, 1)
		require.Equal(t, "Tuesday", availsB[0].DayOfWeek)
		require.Equal(t, "14:00", availsB[0].StartTime)
	})
}

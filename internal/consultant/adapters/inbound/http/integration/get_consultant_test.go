package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	consultant_dto "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestGetConsultant_Integration(t *testing.T) {
	validConsultantReq := consultant_dto.RegisterConsultantDTO{
		ProfessionID:    "prof_9ee432d7-b672-40ae-b03f-c1f1fb696621",
		DisplayName:     "Jane Doe Tech",
		Bio:             "Experienced software engineer with 10 years of experience.",
		YearsExperience: 10,
	}

	t.Run("should get private consultant profile for authenticated user successfully", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "get.private@example.com")

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

		// 2. Fetch current user info to verify user ID
		meResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/identity/me",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, meResp.StatusCode)
		meRes := testhelpers.DecodeResponse(t, meResp)
		meResp.Body.Close()

		dataBytes, err := json.Marshal(meRes.Data)
		require.NoError(t, err)
		var userMap map[string]interface{}
		err = json.Unmarshal(dataBytes, &userMap)
		require.NoError(t, err)
		userID := userMap["id"].(string)

		// 3. Fetch private consultant profile via GET /consultants/user
		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/consultants/user",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		res := testhelpers.DecodeResponse(t, resp)
		require.True(t, res.Success)
		require.Equal(t, "Consultant fetched successfully", res.Message)

		consultantBytes, err := json.Marshal(res.Data)
		require.NoError(t, err)

		var consultantResp consultant_dto.PrivateConsultantResponseDTO
		err = json.Unmarshal(consultantBytes, &consultantResp)
		require.NoError(t, err)

		require.NotEmpty(t, consultantResp.ID)
		require.Equal(t, userID, consultantResp.UserID)
		require.Equal(t, "SOFTWARE_ENGINEER", consultantResp.Profession)
		require.Equal(t, "Jane Doe Tech", consultantResp.DisplayName)
		require.Equal(t, "Experienced software engineer with 10 years of experience.", consultantResp.Bio)
		require.Equal(t, 10, consultantResp.YearsExperience)
		require.True(t, consultantResp.IsAcceptingClients)
		require.NotEmpty(t, consultantResp.CreatedAt)
		require.NotEmpty(t, consultantResp.UpdatedAt)
	})

	t.Run("should get public consultant profile by ID successfully", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "get.byid@example.com")

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

		// 2. Fetch private consultant profile to retrieve generated consultant ID and User ID
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
		userID := consultantMap["user_id"].(string)

		// 3. Public endpoint GET /consultants/:id (no session cookie required)
		resp, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodGet,
			"/test/v1/consultants/"+consultantID,
			nil,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		res := testhelpers.DecodeResponse(t, resp)
		require.True(t, res.Success)
		require.Equal(t, "Consultant fetched successfully", res.Message)

		pubBytes, err := json.Marshal(res.Data)
		require.NoError(t, err)

		var pubResp consultant_dto.PublicConsultantResponseDTO
		err = json.Unmarshal(pubBytes, &pubResp)
		require.NoError(t, err)

		require.Equal(t, consultantID, pubResp.ID)
		require.Equal(t, userID, pubResp.UserID)
		require.Equal(t, "SOFTWARE_ENGINEER", pubResp.Profession)
		require.Equal(t, "Jane Doe Tech", pubResp.DisplayName)
		require.Equal(t, "Experienced software engineer with 10 years of experience.", pubResp.Bio)
		require.Equal(t, 10, pubResp.YearsExperience)
	})

	t.Run("should fail when consultant is not found by ID", func(t *testing.T) {
		harness := setUpConsultantApp(t)

		resp, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodGet,
			"/test/v1/consultants/non_existent_id_9999",
			nil,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "consultant not found", errRes.Message)
	})

	t.Run("should fail when private consultant profile is not found for authenticated user", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "notfound.user@example.com")

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/consultants/user",
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

	t.Run("should fail when getting private consultant profile unauthenticated", func(t *testing.T) {
		harness := setUpConsultantApp(t)

		resp, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodGet,
			"/test/v1/consultants/user",
			nil,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		errRes := testhelpers.DecodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "User is unauthorized", errRes.Message)
	})
}

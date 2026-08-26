package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	case_dto "github.com/AppeiYA/consultation-platform/internal/consultationcase/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestDeleteConsultationCase_Integration(t *testing.T) {
	initialCaseReq := case_dto.CreateConsultationCaseDTO{
		Title:       "Case to Delete",
		Description: "This case will be deleted in the test.",
		Category:    "GENERAL",
	}

	t.Run("should delete consultation case successfully", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "delete.success@example.com")

		// 1. Create case
		createResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultation-cases",
			initialCaseReq,
			cookieHeader,
		)
		require.NoError(t, err)
		createResp.Body.Close()
		require.Equal(t, http.StatusCreated, createResp.StatusCode)

		// 2. List to get ID
		listResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/consultation-cases",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer listResp.Body.Close()

		res := testhelpers.DecodeResponse(t, listResp)
		var cases []case_dto.ConsultationCasesDTO
		b, _ := json.Marshal(res.Data)
		_ = json.Unmarshal(b, &cases)
		require.Len(t, cases, 1)
		caseID := cases[0].ID

		// 3. Delete case
		delResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodDelete,
			"/test/v1/consultation-cases/"+caseID,
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer delResp.Body.Close()

		require.Equal(t, http.StatusOK, delResp.StatusCode)

		delRes := testhelpers.DecodeResponse(t, delResp)
		require.True(t, delRes.Success)
		require.Equal(t, "Case deleted successfully", delRes.Message)

		// 4. Verify case is gone (GET returns 404)
		getResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/consultation-cases/"+caseID,
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer getResp.Body.Close()

		require.Equal(t, http.StatusNotFound, getResp.StatusCode)
	})

	t.Run("should fail when deleting non-existent case", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "delnotfound@example.com")

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodDelete,
			"/test/v1/consultation-cases/case_nonexistent_12345",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should fail when deleting case belonging to another user", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)
		user1Cookie := registerAndLoginUser(t, harness, "user1.del@example.com")
		user2Cookie := registerAndLoginUser(t, harness, "user2.del@example.com")

		// User1 creates case
		createResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultation-cases",
			initialCaseReq,
			user1Cookie,
		)
		require.NoError(t, err)
		createResp.Body.Close()

		listResp, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodGet, "/test/v1/consultation-cases", nil, user1Cookie)
		require.NoError(t, err)
		defer listResp.Body.Close()

		res := testhelpers.DecodeResponse(t, listResp)
		var cases []case_dto.ConsultationCasesDTO
		b, _ := json.Marshal(res.Data)
		_ = json.Unmarshal(b, &cases)
		caseID := cases[0].ID

		// User2 attempts deletion
		delResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodDelete,
			"/test/v1/consultation-cases/"+caseID,
			nil,
			user2Cookie,
		)
		require.NoError(t, err)
		defer delResp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, delResp.StatusCode)
	})

	t.Run("should fail when unauthenticated", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)

		resp, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodDelete,
			"/test/v1/consultation-cases/case_123",
			nil,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

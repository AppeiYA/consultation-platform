package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	case_dto "github.com/AppeiYA/consultation-platform/internal/consultationcase/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestUpdateConsultationCase_Integration(t *testing.T) {
	initialCaseReq := case_dto.CreateConsultationCaseDTO{
		Title:       "Initial Consultation Title",
		Description: "Initial description for consultation case.",
		Category:    "FINANCE",
	}

	t.Run("should partially update consultation case successfully", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "update.success@example.com")

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

		// 2. Get case ID
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

		// 3. Update case (patch title & description)
		newTitle := "Updated Consultation Title"
		newDesc := "Updated detailed consultation description."
		updateReq := case_dto.UpdateConsultationCaseDTO{
			Title:       &newTitle,
			Description: &newDesc,
		}

		patchResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPatch,
			"/test/v1/consultation-cases/"+caseID,
			updateReq,
			cookieHeader,
		)
		require.NoError(t, err)
		defer patchResp.Body.Close()

		require.Equal(t, http.StatusOK, patchResp.StatusCode)

		updateRes := testhelpers.DecodeResponse(t, patchResp)
		require.True(t, updateRes.Success)
		require.Equal(t, "Case updated successfully", updateRes.Message)

		// 4. Verify updated fields with GET
		getResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/consultation-cases/"+caseID,
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer getResp.Body.Close()

		getRes := testhelpers.DecodeResponse(t, getResp)
		var fetchedCase case_dto.ConsultationCasesDTO
		gb, _ := json.Marshal(getRes.Data)
		_ = json.Unmarshal(gb, &fetchedCase)

		require.Equal(t, "Updated Consultation Title", fetchedCase.Title)
		require.Equal(t, "Updated detailed consultation description.", fetchedCase.Description)
		require.Equal(t, "FINANCE", fetchedCase.Category)
	})

	t.Run("should fail when update body provides no fields", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "emptyupdate@example.com")

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPatch,
			"/test/v1/consultation-cases/case_123",
			case_dto.UpdateConsultationCaseDTO{},
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail when updating non-existent case", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "notfoundupdate@example.com")

		newTitle := "New Title"
		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPatch,
			"/test/v1/consultation-cases/case_nonexistent_12345",
			case_dto.UpdateConsultationCaseDTO{Title: &newTitle},
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should fail when updating case belonging to another user", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)
		user1Cookie := registerAndLoginUser(t, harness, "user1.update@example.com")
		user2Cookie := registerAndLoginUser(t, harness, "user2.update@example.com")

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

		// User2 attempts update
		newTitle := "Hacked Title"
		patchResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPatch,
			"/test/v1/consultation-cases/"+caseID,
			case_dto.UpdateConsultationCaseDTO{Title: &newTitle},
			user2Cookie,
		)
		require.NoError(t, err)
		defer patchResp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, patchResp.StatusCode)
	})
}

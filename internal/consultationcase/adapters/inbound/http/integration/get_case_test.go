package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	case_dto "github.com/AppeiYA/consultation-platform/internal/consultationcase/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestGetConsultationCase_Integration(t *testing.T) {
	validCaseReq := case_dto.CreateConsultationCaseDTO{
		Title:       "Employment Law Advice",
		Description: "Seeking legal guidance regarding new employment contract terms.",
		Category:    "LEGAL",
	}

	t.Run("should get consultation case by ID successfully", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "client.get@example.com")

		// 1. Create case
		createResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultation-cases",
			validCaseReq,
			cookieHeader,
		)
		require.NoError(t, err)
		createResp.Body.Close()
		require.Equal(t, http.StatusCreated, createResp.StatusCode)

		// 2. List cases to get the case ID
		listResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/consultation-cases",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer listResp.Body.Close()
		require.Equal(t, http.StatusOK, listResp.StatusCode)

		res := testhelpers.DecodeResponse(t, listResp)
		var cases []case_dto.ConsultationCasesDTO
		bytesData, err := json.Marshal(res.Data)
		require.NoError(t, err)
		err = json.Unmarshal(bytesData, &cases)
		require.NoError(t, err)
		require.Len(t, cases, 1)

		caseID := cases[0].ID

		// 3. Get case by ID
		getResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/consultation-cases/"+caseID,
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer getResp.Body.Close()

		require.Equal(t, http.StatusOK, getResp.StatusCode)

		getRes := testhelpers.DecodeResponse(t, getResp)
		var fetchedCase case_dto.ConsultationCasesDTO
		fetchBytes, err := json.Marshal(getRes.Data)
		require.NoError(t, err)
		err = json.Unmarshal(fetchBytes, &fetchedCase)
		require.NoError(t, err)

		require.Equal(t, caseID, fetchedCase.ID)
		require.Equal(t, "Employment Law Advice", fetchedCase.Title)
		require.Equal(t, "Seeking legal guidance regarding new employment contract terms.", fetchedCase.Description)
		require.Equal(t, "LEGAL", fetchedCase.Category)
		require.Equal(t, "DRAFT", fetchedCase.Status)
	})

	t.Run("should fail when unauthenticated", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)

		resp, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodGet,
			"/test/v1/consultation-cases/case_random",
			nil,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should fail when case does not exist", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "notfound@example.com")

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/consultation-cases/case_nonexistent_12345",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should fail when case belongs to another client", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)
		user1Cookie := registerAndLoginUser(t, harness, "owner@example.com")
		user2Cookie := registerAndLoginUser(t, harness, "other@example.com")

		// 1. User1 creates case
		createResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultation-cases",
			validCaseReq,
			user1Cookie,
		)
		require.NoError(t, err)
		createResp.Body.Close()
		require.Equal(t, http.StatusCreated, createResp.StatusCode)

		// Get case ID from user1 list
		listResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/consultation-cases",
			nil,
			user1Cookie,
		)
		require.NoError(t, err)
		defer listResp.Body.Close()

		res := testhelpers.DecodeResponse(t, listResp)
		var cases []case_dto.ConsultationCasesDTO
		bytesData, _ := json.Marshal(res.Data)
		_ = json.Unmarshal(bytesData, &cases)
		require.Len(t, cases, 1)
		caseID := cases[0].ID

		// 2. User2 attempts to get User1's case
		getResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/consultation-cases/"+caseID,
			nil,
			user2Cookie,
		)
		require.NoError(t, err)
		defer getResp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, getResp.StatusCode)
	})
}

package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	case_dto "github.com/AppeiYA/consultation-platform/internal/consultationcase/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestListConsultationCases_Integration(t *testing.T) {
	t.Run("should return empty list when client has no cases", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "nocases@example.com")

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/consultation-cases",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		res := testhelpers.DecodeResponse(t, resp)
		require.True(t, res.Success)
		require.Equal(t, "Cases fetched successfully", res.Message)
	})

	t.Run("should list all consultation cases belonging to authenticated user", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)
		user1Cookie := registerAndLoginUser(t, harness, "user1.list@example.com")
		user2Cookie := registerAndLoginUser(t, harness, "user2.list@example.com")

		// Create 2 cases for user1
		c1 := case_dto.CreateConsultationCaseDTO{
			Title:       "User1 Case 1",
			Description: "First case description for user1.",
			Category:    "FINANCE",
		}
		c2 := case_dto.CreateConsultationCaseDTO{
			Title:       "User1 Case 2",
			Description: "Second case description for user1.",
			Category:    "LEGAL",
		}
		r1, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultation-cases", c1, user1Cookie)
		require.NoError(t, err)
		r1.Body.Close()
		require.Equal(t, http.StatusCreated, r1.StatusCode)

		r2, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultation-cases", c2, user1Cookie)
		require.NoError(t, err)
		r2.Body.Close()
		require.Equal(t, http.StatusCreated, r2.StatusCode)

		// Create 1 case for user2
		c3 := case_dto.CreateConsultationCaseDTO{
			Title:       "User2 Case 1",
			Description: "First case description for user2.",
			Category:    "HEALTH",
		}
		r3, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodPost, "/test/v1/consultation-cases", c3, user2Cookie)
		require.NoError(t, err)
		r3.Body.Close()
		require.Equal(t, http.StatusCreated, r3.StatusCode)

		// List user1 cases
		resp1, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodGet, "/test/v1/consultation-cases", nil, user1Cookie)
		require.NoError(t, err)
		defer resp1.Body.Close()
		require.Equal(t, http.StatusOK, resp1.StatusCode)

		res1 := testhelpers.DecodeResponse(t, resp1)
		var user1Cases []case_dto.ConsultationCasesDTO
		b1, _ := json.Marshal(res1.Data)
		_ = json.Unmarshal(b1, &user1Cases)
		require.Len(t, user1Cases, 2)

		// List user2 cases
		resp2, err := testhelpers.PerformRequestWithCookie(harness.App, http.MethodGet, "/test/v1/consultation-cases", nil, user2Cookie)
		require.NoError(t, err)
		defer resp2.Body.Close()
		require.Equal(t, http.StatusOK, resp2.StatusCode)

		res2 := testhelpers.DecodeResponse(t, resp2)
		var user2Cases []case_dto.ConsultationCasesDTO
		b2, _ := json.Marshal(res2.Data)
		_ = json.Unmarshal(b2, &user2Cases)
		require.Len(t, user2Cases, 1)
		require.Equal(t, "User2 Case 1", user2Cases[0].Title)
	})

	t.Run("should fail when unauthenticated", func(t *testing.T) {
		harness := setUpConsultationCaseApp(t)

		resp, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodGet,
			"/test/v1/consultation-cases",
			nil,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

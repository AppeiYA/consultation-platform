package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	consultant_dto "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestListProfessions_Integration(t *testing.T) {
	t.Run("should list all seeded professions successfully without authentication", func(t *testing.T) {
		harness := setUpConsultantApp(t)

		resp, err := testhelpers.PerformRequest(
			harness.App,
			http.MethodGet,
			"/test/v1/consultants/professions",
			nil,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		res := testhelpers.DecodeResponse(t, resp)
		require.True(t, res.Success)
		require.Equal(t, "Professions fetched successfully", res.Message)

		dataBytes, err := json.Marshal(res.Data)
		require.NoError(t, err)

		var professions []consultant_dto.ListProfessionsResponse
		err = json.Unmarshal(dataBytes, &professions)
		require.NoError(t, err)

		require.GreaterOrEqual(t, len(professions), 6)

		professionMap := make(map[string]string)
		for _, p := range professions {
			professionMap[p.Name] = p.ID
		}

		expectedProfessions := map[string]string{
			"SOFTWARE_ENGINEER": "prof_9ee432d7-b672-40ae-b03f-c1f1fb696621",
			"LAWYER":            "prof_12d965f5-e1f5-49aa-ac57-856772d236ce",
			"DOCTOR":            "prof_940f840d-617a-4ead-8a8c-873851762bc7",
			"ACCOUNTANT":        "prof_9a1e78f7-a053-4b4b-8802-fcac77e530ec",
			"THERAPIST":         "prof_d95d5c58-d5be-4bca-bf87-b84b3f2a2681",
			"CLERGY":            "prof_aef03e88-a0ee-49c9-b455-4b2210412b52",
		}

		for name, expectedID := range expectedProfessions {
			id, exists := professionMap[name]
			require.True(t, exists, "profession %s should exist", name)
			require.Equal(t, expectedID, id)
		}
	})

	t.Run("should list professions successfully when authenticated", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "prof.list@example.com")

		resp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/consultants/professions",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		res := testhelpers.DecodeResponse(t, resp)
		require.True(t, res.Success)
		require.Equal(t, "Professions fetched successfully", res.Message)

		dataBytes, err := json.Marshal(res.Data)
		require.NoError(t, err)

		var professions []consultant_dto.ListProfessionsResponse
		err = json.Unmarshal(dataBytes, &professions)
		require.NoError(t, err)

		require.GreaterOrEqual(t, len(professions), 6)
	})
}

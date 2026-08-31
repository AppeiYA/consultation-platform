package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	consultant_dto "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestConsultantExpertise_Integration(t *testing.T) {
	t.Run("should add, list, replace, and delete expertises", func(t *testing.T) {
		harness := setUpConsultantApp(t)
		cookieHeader := registerAndLoginUser(t, harness, "expertise.test@example.com")

		// 1. Register consultant with initial expertises
		regReq := consultant_dto.RegisterConsultantDTO{
			ProfessionID:    "prof_9ee432d7-b672-40ae-b03f-c1f1fb696621",
			DisplayName:     "Expertise Consultant",
			Bio:             "Senior cloud architect with deep distributed systems knowledge.",
			YearsExperience: 8,
			Expertises:      []string{"Go", "Distributed Systems"},
		}

		regResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/register",
			regReq,
			cookieHeader,
		)
		require.NoError(t, err)
		regResp.Body.Close()
		require.Equal(t, http.StatusCreated, regResp.StatusCode)

		// 2. Fetch private consultant profile to verify initial expertises
		profileResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/consultants/user",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer profileResp.Body.Close()
		require.Equal(t, http.StatusOK, profileResp.StatusCode)

		profileDecoded := testhelpers.DecodeResponse(t, profileResp)
		profileBytes, err := json.Marshal(profileDecoded.Data)
		require.NoError(t, err)

		var consultantProfile consultant_dto.PrivateConsultantResponseDTO
		err = json.Unmarshal(profileBytes, &consultantProfile)
		require.NoError(t, err)
		require.ElementsMatch(t, []string{"Go", "Distributed Systems"}, consultantProfile.Expertises)

		// 3. Add a new expertise
		addResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/consultants/me/expertises",
			consultant_dto.AddExpertiseDTO{Name: "PostgreSQL"},
			cookieHeader,
		)
		require.NoError(t, err)
		defer addResp.Body.Close()
		require.Equal(t, http.StatusCreated, addResp.StatusCode)

		// 4. List expertises
		listResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/consultants/me/expertises",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer listResp.Body.Close()
		require.Equal(t, http.StatusOK, listResp.StatusCode)

		listDecoded := testhelpers.DecodeResponse(t, listResp)
		listBytes, err := json.Marshal(listDecoded.Data)
		require.NoError(t, err)

		var expertisesList []consultant_dto.ExpertiseResponseDTO
		err = json.Unmarshal(listBytes, &expertisesList)
		require.NoError(t, err)
		require.Len(t, expertisesList, 3)

		var addedExpID string
		for _, exp := range expertisesList {
			if exp.Name == "PostgreSQL" {
				addedExpID = exp.ID
			}
		}
		require.NotEmpty(t, addedExpID)

		// 5. Delete the added expertise
		delResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodDelete,
			"/test/v1/consultants/me/expertises/"+addedExpID,
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		delResp.Body.Close()
		require.Equal(t, http.StatusOK, delResp.StatusCode)

		// 6. Bulk replace expertises
		replaceResp, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodPut,
			"/test/v1/consultants/me/expertises",
			consultant_dto.ReplaceExpertisesDTO{
				Expertises: []string{"Kubernetes", "gRPC", "Observability"},
			},
			cookieHeader,
		)
		require.NoError(t, err)
		defer replaceResp.Body.Close()
		require.Equal(t, http.StatusOK, replaceResp.StatusCode)

		// 7. Verify updated list
		listResp2, err := testhelpers.PerformRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/consultants/me/expertises",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer listResp2.Body.Close()
		require.Equal(t, http.StatusOK, listResp2.StatusCode)

		listDecoded2 := testhelpers.DecodeResponse(t, listResp2)
		listBytes2, err := json.Marshal(listDecoded2.Data)
		require.NoError(t, err)

		var expertisesList2 []consultant_dto.ExpertiseResponseDTO
		err = json.Unmarshal(listBytes2, &expertisesList2)
		require.NoError(t, err)
		require.Len(t, expertisesList2, 3)
		names := []string{expertisesList2[0].Name, expertisesList2[1].Name, expertisesList2[2].Name}
		require.ElementsMatch(t, []string{"Kubernetes", "gRPC", "Observability"}, names)
	})
}

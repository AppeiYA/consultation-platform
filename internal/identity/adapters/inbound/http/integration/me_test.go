package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/dto"
	"github.com/stretchr/testify/require"
)

func TestGetMe_Integration(t *testing.T) {
	t.Run("should fetch current user successfully", func(t *testing.T) {
		harness := setUpIdentityApp(t)

		// Register user
		regBody := dto.RegisterRequest{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane.doe@example.com",
			Password:  "SecurePassword123!",
		}

		regResp, err := performRequest(
			harness.App,
			http.MethodPost,
			"/test/v1/identity/register",
			regBody,
		)
		require.NoError(t, err)
		regResp.Body.Close()
		require.Equal(t, http.StatusCreated, regResp.StatusCode)

		// Login user to get session cookie
		loginBody := dto.LoginRequest{
			Email:    "jane.doe@example.com",
			Password: "SecurePassword123!",
		}

		loginResp, err := performRequest(
			harness.App,
			http.MethodPost,
			"/test/v1/identity/login",
			loginBody,
		)
		require.NoError(t, err)
		cookieHeader := extractCookieHeader(loginResp)
		loginResp.Body.Close()
		require.NotEmpty(t, cookieHeader)

		// Fetch current user
		resp, err := performRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/identity/me",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		res := decodeResponse(t, resp)
		require.True(t, res.Success)
		require.Equal(t, "user fetched successfully", res.Message)

		// Unmarshal Data into dto.GetMeResponse
		dataBytes, err := json.Marshal(res.Data)
		require.NoError(t, err)

		var meResp dto.GetMeResponse
		err = json.Unmarshal(dataBytes, &meResp)
		require.NoError(t, err)

		require.NotEmpty(t, meResp.ID)
		require.Equal(t, "Jane", meResp.FirstName)
		require.Equal(t, "Doe", meResp.LastName)
		require.Equal(t, "jane.doe@example.com", meResp.Email)
		require.Equal(t, "CLIENT", meResp.Role)
	})

	t.Run("should fail when unauthenticated (no session cookie)", func(t *testing.T) {
		harness := setUpIdentityApp(t)

		resp, err := performRequest(
			harness.App,
			http.MethodGet,
			"/test/v1/identity/me",
			nil,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		errRes := decodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "User is unauthorized", errRes.Message)
	})

	t.Run("should fail when session cookie is invalid", func(t *testing.T) {
		harness := setUpIdentityApp(t)

		resp, err := performRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/identity/me",
			nil,
			"test_session_id=invalid_token_12345678901234567890123456789012",
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		errRes := decodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "User is unauthorized", errRes.Message)
	})
}

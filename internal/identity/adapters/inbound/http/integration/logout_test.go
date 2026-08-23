package integration

import (
	"net/http"
	"strings"
	"testing"

	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/dto"
	"github.com/stretchr/testify/require"
)

func TestLogoutUser_Integration(t *testing.T) {
	t.Run("should logout user successfully", func(t *testing.T) {
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

		// Logout user
		resp, err := performRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/identity/logout",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		res := decodeResponse(t, resp)
		require.True(t, res.Success)
		require.Equal(t, "Logout Successful", res.Message)

		// Check session cookie is deleted/cleared in response header
		setCookie := resp.Header.Get("Set-Cookie")
		require.NotEmpty(t, setCookie)
		require.True(t, strings.Contains(setCookie, "session_id=;") || strings.Contains(setCookie, "Max-Age=-1") || strings.Contains(setCookie, "Expires=Thu, 01 Jan 1970"))

		// Verify that using the session cookie again fails authentication
		meResp, err := performRequestWithCookie(
			harness.App,
			http.MethodGet,
			"/test/v1/identity/me",
			nil,
			cookieHeader,
		)
		require.NoError(t, err)
		defer meResp.Body.Close()

		require.Equal(t, http.StatusNotFound, meResp.StatusCode)
	})

	t.Run("should fail when unauthenticated (no session cookie)", func(t *testing.T) {
		harness := setUpIdentityApp(t)

		resp, err := performRequest(
			harness.App,
			http.MethodPost,
			"/test/v1/identity/logout",
			nil,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should fail when session token format is invalid", func(t *testing.T) {
		harness := setUpIdentityApp(t)

		resp, err := performRequestWithCookie(
			harness.App,
			http.MethodPost,
			"/test/v1/identity/logout",
			nil,
			"test_session_id=invalid_token",
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := decodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "failed to logout", errRes.Message)
		require.Equal(t, "session token is too short", errRes.Error)
	})
}

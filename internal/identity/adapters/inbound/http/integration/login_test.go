package integration

import (
	"net/http"
	"strings"
	"testing"

	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/dto"
	"github.com/stretchr/testify/require"
)

func TestLoginUser_Integration(t *testing.T) {
	t.Run("should login user successfully", func(t *testing.T) {
		harness := setUpIdentityApp(t)

		// Register user first
		regBody := dto.RegisterRequest{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane.doe@example.com",
			Password:  "SecurePassword123!",
			Role:      "CLIENT",
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

		// Perform login
		loginBody := dto.LoginRequest{
			Email:    "jane.doe@example.com",
			Password: "SecurePassword123!",
		}

		resp, err := performRequest(
			harness.App,
			http.MethodPost,
			"/test/v1/identity/login",
			loginBody,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		res := decodeResponse(t, resp)
		require.True(t, res.Success)
		require.Equal(t, "login success", res.Message)

		// Verify session cookie is set
		setCookie := resp.Header.Get("Set-Cookie")
		require.NotEmpty(t, setCookie)
		require.True(t, strings.Contains(setCookie, "session_id="))
	})

	t.Run("should reject invalid email", func(t *testing.T) {
		harness := setUpIdentityApp(t)

		loginBody := dto.LoginRequest{
			Email:    "invalid-email",
			Password: "SecurePassword123!",
		}

		resp, err := performRequest(
			harness.App,
			http.MethodPost,
			"/test/v1/identity/login",
			loginBody,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := decodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "failed to login", errRes.Message)
		require.Equal(t, "invalid email address", errRes.Error)
	})

	t.Run("should reject weak password", func(t *testing.T) {
		harness := setUpIdentityApp(t)

		loginBody := dto.LoginRequest{
			Email:    "jane.doe@example.com",
			Password: "weak",
		}

		resp, err := performRequest(
			harness.App,
			http.MethodPost,
			"/test/v1/identity/login",
			loginBody,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errRes := decodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "failed to login", errRes.Message)
		require.Equal(t, "password must be at least 8 characters long", errRes.Error)
	})

	t.Run("should fail when user is not found", func(t *testing.T) {
		harness := setUpIdentityApp(t)

		loginBody := dto.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "SecurePassword123!",
		}

		resp, err := performRequest(
			harness.App,
			http.MethodPost,
			"/test/v1/identity/login",
			loginBody,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		errRes := decodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "failed to login", errRes.Message)
		require.Equal(t, "user not found", errRes.Error)
	})

	t.Run("should fail when password is incorrect", func(t *testing.T) {
		harness := setUpIdentityApp(t)

		// Register user first
		regBody := dto.RegisterRequest{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane.doe@example.com",
			Password:  "SecurePassword123!",
			Role:      "CLIENT",
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

		// Attempt login with wrong password
		loginBody := dto.LoginRequest{
			Email:    "jane.doe@example.com",
			Password: "WrongPassword123!",
		}

		resp, err := performRequest(
			harness.App,
			http.MethodPost,
			"/test/v1/identity/login",
			loginBody,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		errRes := decodeErrorResponse(t, resp)
		require.False(t, errRes.Success)
		require.Equal(t, "failed to login", errRes.Message)
		require.Equal(t, "invalid password", errRes.Error)
	})
}

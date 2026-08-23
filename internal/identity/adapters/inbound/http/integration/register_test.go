package integration

import (
	"net/http"
	"testing"

	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/dto"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser_Integration(t *testing.T) {
	t.Run("should register user successfully", func(t *testing.T) {
		harness := setUpIdentityApp(t)

		reqBody := dto.RegisterRequest{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane.doe@example.com",
			Password:  "SecurePassword123!",
		}

		resp, err := performRequest(
			harness.App,
			http.MethodPost,
			"/test/v1/identity/register",
			reqBody,
		)

		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		res := decodeResponse(t, resp)

		require.True(t, res.Success)
		require.Equal(t, "user registered successfully", res.Message)
	})

	t.Run("should reject invalid email", func(t *testing.T) {
		harness := setUpIdentityApp(t)

		reqBody := dto.RegisterRequest{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "invalid-email",
			Password:  "SecurePassword123!",
		}

		resp, err := performRequest(
			harness.App,
			http.MethodPost,
			"/test/v1/identity/register",
			reqBody,
		)

		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		res := decodeResponse(t, resp)

		require.False(t, res.Success)
		require.Equal(t, "invalid email address", res.Message)
	})

	t.Run("should reject weak password", func(t *testing.T) {
		harness := setUpIdentityApp(t)

		reqBody := dto.RegisterRequest{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane.doe@example.com",
			Password:  "weak",
		}

		resp, err := performRequest(
			harness.App,
			http.MethodPost,
			"/test/v1/identity/register",
			reqBody,
		)

		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		res := decodeResponse(t, resp)

		require.False(t, res.Success)
		require.Equal(t, "password must be at least 8 characters long", res.Message)
	})

	t.Run("should reject duplicate email", func(t *testing.T) {
		harness := setUpIdentityApp(t)

		reqBody := dto.RegisterRequest{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane.doe@example.com",
			Password:  "SecurePassword123!",
		}

		// First registration - should succeed
		resp1, err := performRequest(
			harness.App,
			http.MethodPost,
			"/test/v1/identity/register",
			reqBody,
		)
		require.NoError(t, err)
		resp1.Body.Close()
		require.Equal(t, http.StatusCreated, resp1.StatusCode)

		// Duplicate registration - should fail
		resp2, err := performRequest(
			harness.App,
			http.MethodPost,
			"/test/v1/identity/register",
			reqBody,
		)
		require.NoError(t, err)
		defer resp2.Body.Close()

		require.Equal(t, http.StatusConflict, resp2.StatusCode)

		res := decodeResponse(t, resp2)

		require.False(t, res.Success)
		require.Equal(t, "user already exists", res.Message)
	})
}

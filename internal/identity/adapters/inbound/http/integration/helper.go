package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func performRequest(
	app *fiber.App,
	method string,
	path string,
	body any,
) (*http.Response, error) {
	return performRequestWithCookie(app, method, path, body, "")
}

func performRequestWithCookie(
	app *fiber.App,
	method string,
	path string,
	body any,
	cookieHeader string,
) (*http.Response, error) {
	var bodyReader io.Reader

	if body != nil {
		switch v := body.(type) {
		case io.Reader:
			bodyReader = v
		case []byte:
			bodyReader = bytes.NewReader(v)
		case string:
			bodyReader = strings.NewReader(v)
		default:
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			bodyReader = bytes.NewReader(jsonBytes)
		}
	}

	req := httptest.NewRequest(method, path, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if cookieHeader != "" {
		req.Header.Set("Cookie", cookieHeader)
	}

	return app.Test(req, -1)
}

func extractCookieHeader(resp *http.Response) string {
	setCookie := resp.Header.Get("Set-Cookie")
	if setCookie == "" {
		return ""
	}
	parts := strings.Split(setCookie, ";")
	return parts[0]
}

func decodeResponse(
	t *testing.T,
	resp *http.Response,
) *response.Response {
	t.Helper()

	var res response.Response

	err := json.NewDecoder(resp.Body).Decode(&res)
	require.NoError(t, err)

	return &res
}

func decodeErrorResponse(
	t *testing.T,
	resp *http.Response,
) *response.ErrorResponse {
	t.Helper()

	var res response.ErrorResponse

	err := json.NewDecoder(resp.Body).Decode(&res)
	require.NoError(t, err)

	return &res
}




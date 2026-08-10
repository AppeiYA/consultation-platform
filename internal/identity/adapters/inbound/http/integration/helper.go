package integration

import (
	"net/http"
	"testing"

	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/AppeiYA/consultation-platform/internal/shared/testhelpers"
	"github.com/gofiber/fiber/v2"
)

func performRequest(
	app *fiber.App,
	method string,
	path string,
	body any,
) (*http.Response, error) {
	return testhelpers.PerformRequest(app, method, path, body)
}

func performRequestWithCookie(
	app *fiber.App,
	method string,
	path string,
	body any,
	cookieHeader string,
) (*http.Response, error) {
	return testhelpers.PerformRequestWithCookie(app, method, path, body, cookieHeader)
}

func extractCookieHeader(resp *http.Response) string {
	return testhelpers.ExtractCookieHeader(resp)
}

func decodeResponse(
	t *testing.T,
	resp *http.Response,
) *response.Response {
	return testhelpers.DecodeResponse(t, resp)
}

func decodeErrorResponse(
	t *testing.T,
	resp *http.Response,
) *response.ErrorResponse {
	return testhelpers.DecodeErrorResponse(t, resp)
}

package lampclient

import (
	"fmt"
	"net/http"
)

type colorMessage struct {
	Color uint32 `json:"color"`
}

// errorResponse represents an error reported by an API request.
type errorResponse struct {
	Response *http.Response

	Message string `json:"error"`
}

func (er errorResponse) Error() string {
	return fmt.Sprintf("%s %s: %d %s",
		er.Response.Request.Method, er.Response.Request.URL,
		er.Response.StatusCode, er.Message)
}

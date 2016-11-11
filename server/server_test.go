package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func getTestRequest(method string) *http.Request {
	return &http.Request{
		Method: method,
		Header: http.Header{
			"Sec-Websocket-Version": []string{"13"},
			"Connection":            []string{"Upgrade"},
			"Upgrade":               []string{"websocket"},
			"Sec-Websocket-Key":     []string{"key"},
		},
		URL: &url.URL{Host: "localhost:8000", Scheme: "ws", Path: "/ws"},
	}
}

func TestHandleRequest_405(t *testing.T) {
	w := httptest.NewRecorder()
	HandleRequest(w, getTestRequest("POST"))
	if w.Code != 405 {
		t.Error("Non-GET methods should not be allowed")
	}
}

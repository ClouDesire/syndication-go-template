package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

func TestPostEvent(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /subscription/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("CMW-Auth-Token") != "test-token" {
			t.Error("Expected authentication header to be present")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 1, "deploymentStatus": "PENDING", "paid": true}`))
	})
	mux.HandleFunc("PATCH /subscription/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("CMW-Auth-Token") != "test-token" {
			t.Error("Expected authentication header to be present")
		}
		w.WriteHeader(http.StatusNoContent)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	os.Setenv("CMW_BASE_URL", server.URL)
	os.Setenv("CMW_AUTH_TOKEN", "test-token")

	router := gin.Default()
	router = postEvent(router)

	w := httptest.NewRecorder()

	exampleEvent := Event{
		Entity: "Subscription",
		ID:     1,
		Type:   "CREATED",
	}
	eventJson, _ := json.Marshal(exampleEvent)
	req, _ := http.NewRequest("POST", "/event", strings.NewReader(string(eventJson)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
	assert.Empty(t, w.Body)
}

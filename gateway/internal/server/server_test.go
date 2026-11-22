package server

import (
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AFK068/skilloom/gateway/internal/config"
)

func TestServer(t *testing.T) {
	t.Parallel()

	t.Run("occupied port", func(t *testing.T) {
		t.Parallel()

		occupied := httptest.NewServer(http.NotFoundHandler())
		defer occupied.Close()

		srv := New(config.Config{Listen: occupied.Listener.Addr().String()})
		assert.ErrorIs(t, srv.Start(), syscall.EADDRINUSE)
	})

	t.Run("configured routes", func(t *testing.T) {
		t.Parallel()

		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Upstream-Request", r.Method+" "+r.URL.Path)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer upstream.Close()

		router := New(config.Config{Routes: []config.Route{
			{Prefix: "/openai/v1", Upstream: upstream.URL + "/v1", Endpoints: []config.Endpoint{
				{Method: "GET", Path: "/models"},
				{Method: "POST", Path: "/responses"},
			}},
			{Prefix: "/anthropic/v1", Upstream: upstream.URL + "/api", Endpoints: []config.Endpoint{
				{Method: "POST", Path: "/messages"},
			}},
			{Prefix: "/disabled/v1", Upstream: upstream.URL},
		}}).router

		for _, tc := range []struct {
			method, path string
			status       int
			upstream     string
		}{
			{"GET", "/openai/v1/models", http.StatusNoContent, "GET /v1/models"},
			{"POST", "/openai/v1/responses", http.StatusNoContent, "POST /v1/responses"},
			{"POST", "/anthropic/v1/messages", http.StatusNoContent, "POST /api/messages"},
			{"GET", "/openai/v1/responses", http.StatusNotFound, ""},
			{"POST", "/openai/v1/responses/", http.StatusNotFound, ""},
			{"POST", "/openai/v1/%72esponses", http.StatusNotFound, ""},
			{"POST", "/openai/v1/unknown", http.StatusNotFound, ""},
			{"POST", "/disabled/v1/responses", http.StatusNotFound, ""},
			{"GET", "/healthz", http.StatusNotFound, ""},
		} {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequestWithContext(t.Context(), tc.method, tc.path, nil))

			request := tc.method + " " + tc.path
			assert.Equal(t, tc.status, w.Code, request)
			assert.Equal(t, tc.upstream, w.Header().Get("X-Upstream-Request"), request)
		}
	})
}

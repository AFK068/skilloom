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

	t.Run("proxy prefixes", func(t *testing.T) {
		t.Parallel()

		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Upstream-Request", r.Method+" "+r.URL.RequestURI())
			w.WriteHeader(http.StatusNoContent)
		}))
		defer upstream.Close()

		router := New(config.Config{Routes: []config.Route{
			{Prefix: "/openai/v1", Upstream: upstream.URL + "/v1"},
			{Prefix: "/anthropic/v1", Upstream: upstream.URL + "/api"},
		}}).router

		for _, tc := range []struct {
			method, path string
			status       int
			upstream     string
		}{
			{"GET", "/openai/v1", http.StatusNoContent, "GET /v1/"},
			{"POST", "/openai/v1/new/path?trace=a%2Fb", http.StatusNoContent, "POST /v1/new/path?trace=a%2Fb"},
			{"DELETE", "/anthropic/v1/files/a%2Fb", http.StatusNoContent, "DELETE /api/files/a%2Fb"},
			{"GET", "/openai/v10/models", http.StatusNotFound, ""},
		} {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequestWithContext(t.Context(), tc.method, tc.path, nil))

			request := tc.method + " " + tc.path
			assert.Equal(t, tc.status, w.Code, request)
			assert.Equal(t, tc.upstream, w.Header().Get("X-Upstream-Request"), request)
		}
	})
}

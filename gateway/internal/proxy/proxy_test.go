package proxy_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AFK068/skilloom/gateway/internal/proxy"
)

func TestProxy(t *testing.T) {
	t.Parallel()

	t.Run("forwards request and response", func(t *testing.T) {
		t.Parallel()

		const payload = `{"model":"test"}`
		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST /v1/responses?trace=a%2Fb", r.Method+" "+r.URL.RequestURI())
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.JSONEq(t, payload, string(body))

			w.Header().Set("Retry-After", "10")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = io.WriteString(w, "rate limit")
		}))
		defer upstream.Close()

		req := httptest.NewRequestWithContext(
			t.Context(),
			"POST",
			"/gateway/responses?trace=a%2Fb",
			strings.NewReader(payload),
		)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		proxy.New(upstream.URL+"/v1", "/gateway").ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.Equal(t, "rate limit", w.Body.String())
		assert.Equal(t, "10", w.Header().Get("Retry-After"))
	})

	t.Run("rejects invalid upstreams", func(t *testing.T) {
		t.Parallel()

		for _, upstream := range []string{
			"https://%",
			"https:///v1",
			"http://user@127.0.0.1:0",
			"ftp://127.0.0.1",
		} {
			w := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(t.Context(), "GET", "/gateway/models", nil)
			proxy.New(upstream, "/gateway").ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadGateway, w.Code, upstream)
			assert.Equal(t, "invalid upstream URL\n", w.Body.String(), upstream)
		}
	})
}

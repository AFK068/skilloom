package server

import (
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"

	"github.com/AFK068/skilloom/gateway/internal/config"
)

func TestServer(t *testing.T) {
	t.Parallel()

	t.Run("health endpoint is not registered", func(t *testing.T) {
		t.Parallel()

		srv := New(config.Config{Listen: "127.0.0.1:0"})
		w := httptest.NewRecorder()
		srv.router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/healthz", nil))
		if w.Code != http.StatusNotFound {
			t.Errorf("GET /healthz status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("occupied port", func(t *testing.T) {
		t.Parallel()

		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if err := listener.Close(); err != nil {
				t.Error(err)
			}
		})

		srv := New(config.Config{Listen: listener.Addr().String()})
		if err := srv.Start(); !errors.Is(err, syscall.EADDRINUSE) {
			t.Fatalf("Start() error = %v, want address in use", err)
		}
	})
}

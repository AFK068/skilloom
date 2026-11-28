package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/AFK068/skilloom/ingest/internal/domain/principal"
	"github.com/AFK068/skilloom/ingest/internal/services/auth"
	"github.com/AFK068/skilloom/ingest/internal/transport/http/middleware"
	"github.com/AFK068/skilloom/ingest/internal/transport/http/middleware/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuthHandle(t *testing.T) {
	t.Parallel()

	wantPrincipal := principal.Principal{
		UserID:  uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		SpaceID: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
	}

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		authenticator := mock.NewMockAuthenticator(gomock.NewController(t))
		authenticator.EXPECT().
			Authenticate(gomock.Any(), "secret-token").
			Return(wantPrincipal, nil)

		mw := middleware.NewAuth(middleware.AuthDeps{
			Auth: authenticator,
			Log:  zap.NewNop(),
		})

		var (
			gotPrincipal principal.Principal
			gotOK        bool
		)
		router := gin.New()
		router.GET("/test", mw.Handle, func(c *gin.Context) {
			gotPrincipal, gotOK = middleware.Principal(c)
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
		req.Header.Set("Authorization", "Bearer secret-token")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, gotOK)
		assert.Equal(t, wantPrincipal, gotPrincipal)
	})

	t.Run("missing or invalid bearer header", func(t *testing.T) {
		t.Parallel()

		headers := []string{"", "Basic abc", "Bearer", "Bearer "}
		for _, hdr := range headers {
			t.Run("hdr="+hdr, func(t *testing.T) {
				t.Parallel()

				authenticator := mock.NewMockAuthenticator(gomock.NewController(t))
				mw := middleware.NewAuth(middleware.AuthDeps{
					Auth: authenticator,
					Log:  zap.NewNop(),
				})

				nextCalled := false
				router := gin.New()
				router.GET("/test", mw.Handle, func(c *gin.Context) {
					nextCalled = true
					c.Status(http.StatusOK)
				})

				req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
				if hdr != "" {
					req.Header.Set("Authorization", hdr)
				}

				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)

				assert.Equal(t, http.StatusUnauthorized, rec.Code)
				assert.JSONEq(t, `{"error":"missing or invalid bearer token"}`, rec.Body.String())
				assert.False(t, nextCalled)
			})
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		t.Parallel()

		authenticator := mock.NewMockAuthenticator(gomock.NewController(t))
		authenticator.EXPECT().
			Authenticate(gomock.Any(), "bad-token").
			Return(principal.Principal{}, auth.ErrUnauthorized)

		mw := middleware.NewAuth(middleware.AuthDeps{
			Auth: authenticator,
			Log:  zap.NewNop(),
		})

		nextCalled := false
		router := gin.New()
		router.GET("/test", mw.Handle, func(c *gin.Context) {
			nextCalled = true
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
		req.Header.Set("Authorization", "Bearer bad-token")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.JSONEq(t, `{"error":"invalid token"}`, rec.Body.String())
		assert.False(t, nextCalled)
	})

	t.Run("internal error", func(t *testing.T) {
		t.Parallel()

		authenticator := mock.NewMockAuthenticator(gomock.NewController(t))
		authenticator.EXPECT().
			Authenticate(gomock.Any(), "secret-token").
			Return(principal.Principal{}, errors.New("redis down"))

		mw := middleware.NewAuth(middleware.AuthDeps{
			Auth: authenticator,
			Log:  zap.NewNop(),
		})

		nextCalled := false
		router := gin.New()
		router.GET("/test", mw.Handle, func(c *gin.Context) {
			nextCalled = true
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
		req.Header.Set("Authorization", "Bearer secret-token")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"error":"internal server error"}`, rec.Body.String())
		assert.False(t, nextCalled)
	})
}

func TestPrincipalMissing(t *testing.T) {
	t.Parallel()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	got, ok := middleware.Principal(c)
	assert.False(t, ok)
	assert.Zero(t, got)
}

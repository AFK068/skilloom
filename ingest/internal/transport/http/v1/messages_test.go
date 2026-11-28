package v1_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/AFK068/skilloom/ingest/internal/domain/message"
	"github.com/AFK068/skilloom/ingest/internal/domain/principal"
	"github.com/AFK068/skilloom/ingest/internal/services/messages"
	v1 "github.com/AFK068/skilloom/ingest/internal/transport/http/v1"
	"github.com/AFK068/skilloom/ingest/internal/transport/http/v1/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestRouter(h *v1.MessagesHandler, p *principal.Principal) *gin.Engine {
	r := gin.New()
	group := r.Group("/v1")
	if p != nil {
		group.Use(func(c *gin.Context) {
			c.Set("principal", *p)
			c.Next()
		})
	}
	h.Register(group)
	return r
}

func TestMessagesHandlerCreate(t *testing.T) {
	t.Parallel()

	p := principal.Principal{
		UserID:  uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		SpaceID: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
	}
	eventID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	capturedAt := time.Date(2026, 11, 28, 12, 0, 0, 0, time.UTC)

	validBody := `{
		"event_id": "33333333-3333-3333-3333-333333333333",
		"harness": "claude",
		"session_id": "sess-1",
		"turn_id": "turn-1",
		"role": "user",
		"text": "hello world",
		"captured_at": "2026-11-28T12:00:00Z",
		"meta": {"source": "cli"}
	}`

	expectedMsg := message.Message{
		EventID:    eventID,
		Harness:    "claude",
		SessionID:  "sess-1",
		TurnID:     "turn-1",
		Role:       "user",
		Text:       "hello world",
		CapturedAt: capturedAt,
		Meta:       json.RawMessage(`{"source": "cli"}`),
	}

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		svc := mock.NewMockMessagesService(gomock.NewController(t))
		svc.EXPECT().
			SaveMessages(gomock.Any(), p, []message.Message{expectedMsg}).
			Return(nil)

		router := newTestRouter(v1.NewMessagesHandler(v1.MessagesDeps{
			Messages: svc,
			Log:      zap.NewNop(),
		}), &p)

		req := httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusAccepted, rec.Code)
	})

	t.Run("missing principal", func(t *testing.T) {
		t.Parallel()

		svc := mock.NewMockMessagesService(gomock.NewController(t))
		router := newTestRouter(v1.NewMessagesHandler(v1.MessagesDeps{
			Messages: svc,
			Log:      zap.NewNop(),
		}), nil)

		req := httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.JSONEq(t, `{"error":"missing principal"}`, rec.Body.String())
	})

	t.Run("invalid body", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			name string
			body string
		}{
			{name: "malformed json", body: "{"},
			{name: "missing required fields", body: `{"event_id":"33333333-3333-3333-3333-333333333333"}`},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				svc := mock.NewMockMessagesService(gomock.NewController(t))
				router := newTestRouter(v1.NewMessagesHandler(v1.MessagesDeps{
					Messages: svc,
					Log:      zap.NewNop(),
				}), &p)

				req := httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(tc.body))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)

				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.JSONEq(t, `{"error":"invalid request body"}`, rec.Body.String())
			})
		}
	})

	t.Run("invalid message from service", func(t *testing.T) {
		t.Parallel()

		svc := mock.NewMockMessagesService(gomock.NewController(t))
		svc.EXPECT().
			SaveMessages(gomock.Any(), p, gomock.Any()).
			Return(messages.ErrInvalidMessage)

		router := newTestRouter(v1.NewMessagesHandler(v1.MessagesDeps{
			Messages: svc,
			Log:      zap.NewNop(),
		}), &p)

		req := httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error":"invalid message"}`, rec.Body.String())
	})

	t.Run("internal error from service", func(t *testing.T) {
		t.Parallel()

		svc := mock.NewMockMessagesService(gomock.NewController(t))
		svc.EXPECT().
			SaveMessages(gomock.Any(), p, gomock.Any()).
			Return(errors.New("nats down"))

		router := newTestRouter(v1.NewMessagesHandler(v1.MessagesDeps{
			Messages: svc,
			Log:      zap.NewNop(),
		}), &p)

		req := httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"error":"failed to save messages"}`, rec.Body.String())
	})
}

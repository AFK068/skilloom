package natsclient_test

import (
	"errors"
	"testing"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/AFK068/skilloom/ingest/internal/natsclient"
	"github.com/AFK068/skilloom/ingest/internal/natsclient/mock"
)

func TestPublish(t *testing.T) {
	t.Parallel()

	t.Run("with msg id", func(t *testing.T) {
		t.Parallel()

		js := mock.NewMockJetStream(gomock.NewController(t))
		js.EXPECT().
			Publish(t.Context(), "ingest.messages", []byte("hi"), gomock.Any()).
			Return(&jetstream.PubAck{}, nil)

		c := &natsclient.Client{JS: js}
		assert.NoError(t, c.Publish(t.Context(), "ingest.messages", "id-1", []byte("hi")))
	})

	t.Run("without msg id", func(t *testing.T) {
		t.Parallel()

		js := mock.NewMockJetStream(gomock.NewController(t))
		js.EXPECT().
			Publish(t.Context(), "ingest.messages", []byte("hi")).
			Return(&jetstream.PubAck{}, nil)

		c := &natsclient.Client{JS: js}
		assert.NoError(t, c.Publish(t.Context(), "ingest.messages", "", []byte("hi")))
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		wantErr := errors.New("nats down")
		js := mock.NewMockJetStream(gomock.NewController(t))
		js.EXPECT().
			Publish(t.Context(), "ingest.messages", []byte("hi"), gomock.Any()).
			Return(nil, wantErr)

		c := &natsclient.Client{JS: js}
		assert.ErrorIs(t, c.Publish(t.Context(), "ingest.messages", "id-1", []byte("hi")), wantErr)
	})
}

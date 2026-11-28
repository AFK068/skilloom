package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/AFK068/skilloom/ingest/internal/domain/principal"
	"github.com/AFK068/skilloom/ingest/internal/store/redis"
)

const tokenKeyPrefix = "token:"

var ErrUnauthorized = errors.New("unauthorized")

type tokenRecord struct {
	UserID  uuid.UUID `json:"user_id"`
	SpaceID uuid.UUID `json:"space_id"`
}

func (s *Service) Authenticate(ctx context.Context, token string) (principal.Principal, error) {
	if token == "" {
		return principal.Principal{}, ErrUnauthorized
	}

	raw, err := s.Redis.Get(ctx, tokenKeyPrefix+hashToken(token))
	if err != nil {
		if errors.Is(err, redis.ErrNotFound) {
			return principal.Principal{}, ErrUnauthorized
		}
		return principal.Principal{}, fmt.Errorf("get token: %w", err)
	}

	var rec tokenRecord
	if err := json.Unmarshal([]byte(raw), &rec); err != nil {
		return principal.Principal{}, fmt.Errorf("unmarshal token record: %w", err)
	}
	if rec.UserID == uuid.Nil || rec.SpaceID == uuid.Nil {
		return principal.Principal{}, errors.New("token record missing user_id or space_id")
	}

	return principal.Principal{
		UserID:  rec.UserID,
		SpaceID: rec.SpaceID,
	}, nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

package principal

import "github.com/google/uuid"

type Principal struct {
	UserID  uuid.UUID
	SpaceID uuid.UUID
}

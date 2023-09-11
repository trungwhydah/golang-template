package base

import (
	"time"

	"github.com/google/uuid"
)

// Entity ....
type Entity struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

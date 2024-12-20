package favorite

import (
	"github.com/google/uuid"
	"time"
)

type Favorite struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	ProjectID  uuid.UUID  `db:"project_id" json:"project_id"`
	OwnerType  OwnerType  `db:"owner_type" json:"owner_type"`
	OwnerID    uuid.UUID  `db:"owner_id" json:"owner_id"`
	ObjectID   uuid.UUID  `db:"object_id" json:"object_id"`
	ObjectType ObjectType `db:"object_type" json:"object_type"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
}

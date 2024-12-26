package dto

import (
	"github.com/google/uuid"
)

type CreateFavoriteRequest struct {
	ProjectID  uuid.UUID `json:"project_id" binding:"required"`
	OwnerType  string    `json:"owner_type" binding:"required"`
	OwnerID    uuid.UUID `json:"owner_id" binding:"required"`
	ObjectID   uuid.UUID `json:"object_id" binding:"required"`
	ObjectType string    `json:"object_type" binding:"required"`
}

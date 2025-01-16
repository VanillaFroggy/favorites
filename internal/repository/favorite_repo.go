package repository

import (
	"favorites/internal/models/favorite"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type FavoriteRepository struct {
	db *sqlx.DB
}

func NewFavoriteRepository(db *sqlx.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) GetPageOfFavoritesByOwnerTypeAndOwnerID(
	ownerType favorite.OwnerType,
	ownerID uuid.UUID,
	limit uint64,
	cursorID uuid.UUID,
) ([]favorite.Favorite, uuid.UUID, error) {
	var favorites []favorite.Favorite
	var args []interface{}
	query := `
		SELECT *
		FROM favorites
		WHERE owner_type = $1
		  AND owner_id = $2
	`
	args = append(args, ownerType, ownerID)
	if cursorID != uuid.Nil {
		query += `
			AND created_at < (SELECT created_at FROM favorites WHERE id = $3)
			ORDER BY created_at DESC
			LIMIT $4
		`
		args = append(args, cursorID)
	} else {
		query += `
			ORDER BY created_at DESC
			LIMIT $3
		`
	}
	args = append(args, limit)
	err := r.db.Select(&favorites, query, args...)
	if err != nil {
		return nil, uuid.Nil, err
	}
	var nextCursor uuid.UUID
	if len(favorites) > 0 {
		nextCursor = favorites[len(favorites)-1].ID
	}
	return favorites, nextCursor, err
}

func (r *FavoriteRepository) CreateFavorite(f *favorite.Favorite) error {
	query := `INSERT INTO favorites (project_id, owner_type, owner_id, object_id, object_type)
	          VALUES ($1, $2, $3, $4, $5)
	          RETURNING id, project_id, owner_type, owner_id, object_id, object_type, created_at;`
	err := r.db.QueryRowx(
		query,
		f.ProjectID,
		f.OwnerType,
		f.OwnerID,
		f.ObjectID,
		f.ObjectType,
	).StructScan(f)
	return err
}

func (r *FavoriteRepository) DeleteFavorite(id uuid.UUID) error {
	query := `DELETE FROM favorites WHERE id = $1;`
	_, err := r.db.Exec(query, id)
	return err
}

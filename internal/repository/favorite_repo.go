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

func (r *FavoriteRepository) GetAllFavorites(ownerType favorite.OwnerType, ownerID uuid.UUID, limit uint64, offset uint64) ([]favorite.Favorite, error) {
	var favorites []favorite.Favorite
	query := `SELECT *
			  FROM favorites
			  WHERE owner_type = $1 AND owner_id = $2
			  GROUP BY owner_type, owner_id, id
			  ORDER BY created_at DESC
			  LIMIT $3 OFFSET $4;`
	err := r.db.Select(
		&favorites,
		query,
		ownerType,
		ownerID,
		limit,
		offset,
	)
	return favorites, err
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

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

func (r *FavoriteRepository) GetAllFavorites() ([]favorite.Favorite, error) {
	var favorites []favorite.Favorite
	query := `SELECT * FROM favorites`
	err := r.db.Select(&favorites, query)
	return favorites, err
}

func (r *FavoriteRepository) CreateFavorite(fav *favorite.Favorite) error {
	query := `INSERT INTO favorites (project_id, owner_type, owner_id, object_id, object_type)
	          VALUES ($1, $2, $3, $4, $5)
	          RETURNING id, project_id, owner_type, owner_id, object_id, object_type, created_at;`
	err := r.db.QueryRowx(
		query,
		fav.ProjectID,
		fav.OwnerType,
		fav.OwnerID,
		fav.ObjectID,
		fav.ObjectType,
	).StructScan(fav)
	return err
}

func (r *FavoriteRepository) DeleteFavorite(id uuid.UUID) error {
	query := `DELETE FROM favorites WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

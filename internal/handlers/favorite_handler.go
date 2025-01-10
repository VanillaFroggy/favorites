package handlers

import (
	"favorites/internal/handlers/dto"
	"favorites/internal/models/favorite"
	"favorites/internal/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
)

func RegisterRoutes(db *sqlx.DB, r *gin.Engine) {
	repo := repository.NewFavoriteRepository(db)
	r.GET("/favorites", func(c *gin.Context) {
		if !favorite.IsValidOwnerType(c.Query("owner_type")) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Incorrect owner_type '%s', correct is '%s'", c.Query("owner_type"), favorite.OwnerTypeUser)})
			return
		}
		ownerType := favorite.OwnerType(c.Query("owner_type"))
		ownerId, err := uuid.Parse(c.Query("owner_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		limit, err := strconv.ParseUint(c.Query("limit"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		offset, err := strconv.ParseUint(c.Query("offset"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		favorites, err := repo.GetAllFavorites(ownerType, ownerId, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(favorites) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No favorites found"})
			return
		}
		c.JSON(http.StatusOK, favorites)
	})
	r.POST("/favorites", func(c *gin.Context) {
		var request dto.CreateFavoriteRequest
		if err := c.ShouldBind(&request); err != nil || !favorite.IsValidObjectType(request.ObjectType) || !favorite.IsValidOwnerType(request.OwnerType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fav := favorite.Favorite{
			ProjectID:  request.ProjectID,
			OwnerType:  favorite.OwnerType(request.OwnerType),
			OwnerID:    request.OwnerID,
			ObjectID:   request.ObjectID,
			ObjectType: favorite.ObjectType(request.ObjectType),
		}
		err := repo.CreateFavorite(&fav)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, fav)
	})
	r.DELETE("/favorites/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		err = repo.DeleteFavorite(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})
}

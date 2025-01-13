package handlers

import (
	_ "favorites/docs"
	"favorites/internal/handlers/dto"
	"favorites/internal/models/favorite"
	"favorites/internal/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"strconv"
)

var repo *repository.FavoriteRepository

func RegisterRoutes(db *sqlx.DB, r *gin.Engine) {
	repo = repository.NewFavoriteRepository(db)
	r.GET("/favorites", GetFavorites)
	r.POST("/favorites", CreateFavorite)
	r.DELETE("/favorites/:id", DeleteFavorite)
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// GetFavorites godoc
// @Summary       Get favorites array
// @Description   Responds with the page of favorites by owner_type, owner_id, limit and offset as JSON.
// @Tags          favorites
// @Produce       json
// @Param		  owner_type  query    favorite.OwnerType  true  "type of owner"
// @Param		  owner_id  query    string  true  "ID of owner in uuid format"
// @Param		  limit  query    number  true  "size of page"
// @Param		  offset  query    number  true  "count of entries to skip before page"
// @Success       200  {array}  favorite.Favorite
// @Failure       400       {object}  gin.H
// @Failure       404       {object}  gin.H
// @Failure       500       {object}  gin.H
// @Router        /favorites [get]
func GetFavorites(c *gin.Context) {
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
}

// CreateFavorite godoc
// @Summary       Create new favorite
// @Description   Creates a new favorite entry and responses with it as JSON.
// @Tags          favorites
// @Produce       json
// @Param		  request  body    dto.CreateFavoriteRequest  true  "Favorite to create"
// @Success       200  {object}  favorite.Favorite
// @Failure       400       {object}  gin.H
// @Failure       500       {object}  gin.H
// @Router        /favorites [post]
func CreateFavorite(c *gin.Context) {
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
}

// DeleteFavorite godoc
// @Summary       Delete favorite by id
// @Description   Deletes favorite entry and responses with NoContent Code.
// @Tags          favorites
// @Produce       json
// @Param		  id  path    string  true  "ID of favorite to delete in uuid format"
// @Success       201
// @Failure       400       {object}  gin.H
// @Failure       500       {object}  gin.H
// @Router        /favorites/id [delete]
func DeleteFavorite(c *gin.Context) {
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
}

package httputil

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func ParseUUIDFromBase64(c *gin.Context, key string) (uuid.UUID, error) {
	cursorBase64 := c.Query(key)
	var cursorID uuid.UUID
	if cursorBase64 != "" {
		decodedCursor, err := base64.URLEncoding.DecodeString(cursorBase64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cursor"})
			return uuid.UUID{}, err
		}
		cursorID, err = uuid.Parse(string(decodedCursor))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return uuid.UUID{}, err
		}
	}
	return cursorID, nil
}

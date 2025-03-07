package routes

import (
	"context"
	"net/http"
	"tinyurl/database"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func directUrl(c *gin.Context) {
	param := c.Param("url")
	r2 := database.RedisClient(1)
	defer r2.Close()
	val, err := r2.Get(context.Background(), param).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No url present"})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in connecting the client"})
		return
	}
	c.Redirect(301, val)
}

package controllers

import (
	"context"
	"net/http"
	"time"
	"tinyurl/database"
	"tinyurl/models"
	"tinyurl/sendgrid"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Site is working fine"})
}

func VerifyEmail(c *gin.Context) {
	var User models.Verification
	err := c.BindJSON(&User)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "In Binding JSON"})
		return
	}
	r1 := database.RedisClient(2)
	defer r1.Close()

	value, err := r1.Get(context.Background(), User.Email).Result()
	if err == redis.Nil {
		/*_, errs := r1.Set(context.Background(), User.Email, User, 0).Result()
		if errs != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Connectivity in database"})
			return
		}*/
		User, errs := sendgrid.SendEmails(User)
		if errs != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error in sending emails"})
			return
		}
		User.CreatedAT = time.Now().Unix()
		_, errs = r1.Set(context.Background(), User.Email, User, 0).Result()
		if errs != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Connectivity to databse"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Connection in database"})
		return
	}
	if value != "" {
		c.JSON(http.StatusOK, gin.H{"message": "Email already verified"})
	}
}

func VerifyOtp(c *gin.Context) {

}

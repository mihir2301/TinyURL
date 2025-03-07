package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"tinyurl/constants"
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
		User.ISverified = false
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
	var verify models.VerifyOtp
	err := c.BindJSON(&verify)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error at binding json"})
		return
	}
	r1 := database.RedisClient(2)
	defer r1.Close()
	value, err := r1.Get(context.Background(), verify.Email).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email Not VErified first verify your email"})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Connectivity to database"})
		return
	}
	var vals models.Verification
	err = json.Unmarshal([]byte(value), &vals)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in Unmarshalling"})
		return
	}
	sec := vals.CreatedAT + constants.Otpvalidation
	if sec < time.Now().Unix() {
		c.JSON(http.StatusBadRequest, gin.H{"message": "OTP expired"})
		return
	}
	if vals.ISverified {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Otp already verified"})
		return
	} else {
		if verify.Otp == vals.Otp {
			vals.ISverified = true
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Wrong otp, Please provide with correct otp"})
			return
		}
	}
	_, errs := r1.Set(context.Background(), vals.Email, vals, 0).Result()
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Database connectivity issue"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Otp verified successfully"})
}

func ResendOTP(c *gin.Context) {
	var verify models.Verification
	err := c.BindJSON(&verify)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error in binding json"})
		return
	}
	r1 := database.RedisClient(2)
	defer r1.Close()
	value, err := r1.Get(context.Background(), verify.Email).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email not verified, First verify ur email"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error at connectivity"})
		return
	}
	var val models.Verification
	err = json.Unmarshal([]byte(value), &val)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error at unmarshalling json"})
		return
	}
	if val.ISverified {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Otp verified NO need to verify again"})
		return
	}
	val, err = sendgrid.SendEmails(val)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sending email"})
		return
	}
	val.CreatedAT = time.Now().Unix()
	val.ISverified = false
	_, err = r1.Set(context.Background(), val.Email, val, 0).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database connectivity issue"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Otp Sent"})
}

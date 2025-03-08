package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"
	"tinyurl/auth"
	"tinyurl/constants"
	"tinyurl/database"
	"tinyurl/helper"
	"tinyurl/models"
	"tinyurl/sendgrid"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
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

		User, errs := sendgrid.SendEmails(User)
		if errs != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error in sending emails"})
			return
		}

		User.CreatedAT = time.Now().Unix()
		jsondata, err := json.Marshal(User)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "error at marshalling"})
			return
		}
		_, errs = r1.Set(context.Background(), User.Email, jsondata, 0).Result()
		if errs != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Connectivity to databse"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Connection in database 2"})
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
	jsondata, err := json.Marshal(vals)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error at marshalling"})
		return
	}
	_, errs := r1.Set(context.Background(), vals.Email, jsondata, 0).Result()
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Database connectivity issue"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Otp verified successfully"})
}

func Resend(c *gin.Context) {
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
	} else {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "error at connectivity"})
			return
		}
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
	jsondata, err := json.Marshal(val)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error at marshalling"})
		return
	}
	_, err = r1.Set(context.Background(), val.Email, jsondata, 0).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database connectivity issue"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Otp Sent"})
}
func RegisterUser(c *gin.Context) {
	var user models.UserClient
	var dbuser models.Users
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in Binding register user"})
		return
	}
	r1 := database.RedisClient(2) //verification table
	r2 := database.RedisClient(3) // User table
	defer r2.Close()
	defer r1.Close()

	val, err := r1.Get(context.Background(), user.Email).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "first verify your email"})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Connectivityy to Databse Issue"})
		return
	}
	var value models.Verification
	err = json.Unmarshal([]byte(val), &value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error at unmarshalling"})
		return
	}
	if !value.ISverified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please verify your otp"})
		return
	}
	err = helper.CheckDetails(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dbuser.Email = user.Email
	dbuser.Name = user.Name
	dbuser.Password = helper.GenPassHash(user.Password)
	dbuser.Phone = user.Phone
	dbuser.CreatedAT = time.Now().Unix()
	dbuser.UpdatedAT = time.Now().Unix()
	jsondata, err := json.Marshal(dbuser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error at marshalling"})
		return
	}
	_, err = r2.Set(context.Background(), user.Email, jsondata, 0).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	dbuser.Password = ""
	c.JSON(http.StatusOK, gin.H{"message": "user registerd successfuly", "data": dbuser})
}
func UserLogin(c *gin.Context) {
	var user models.Login
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "In binding json"})
		return
	}
	r1 := database.RedisClient(3)
	defer r1.Close()

	value, err := r1.Get(context.Background(), user.Email).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No user found, Please enter correct user details"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error at connecting to database"})
		return
	}
	var val models.Users
	err = json.Unmarshal([]byte(value), &val)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error at Unmarshalling data"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(val.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error at comparing password"})
		return
	}
	jwtwrapper := auth.JWTwrapper{
		SecretKey:      os.Getenv("JwtSecrets"),
		Issuer:         os.Getenv("JwtIssuer"),
		ExpirationTime: 48,
	}
	token, err := jwtwrapper.GenerateToken(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error at generating token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user, "token": token})
}

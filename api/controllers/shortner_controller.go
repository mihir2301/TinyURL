package controllers

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"
	"tinyurl/database"
	"tinyurl/helper"
	"tinyurl/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func Shorten(c *gin.Context) {
	var req models.Request
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Binding error"})
		return
	}
	//enforce http
	req.Url = helper.EnforceUrl(req.Url)
	//check for authentic url
	isURl := helper.CheckUrl(req.Url)
	if !isURl {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please provide a valid url"})
		return
	}
	//check for domain error
	err = helper.CheckDomain(req.Url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		return
	}
	//check for rate limit
	r1 := database.RedisClient(0)
	defer r1.Close()
	value, err := r1.Get(context.Background(), c.ClientIP()).Result()
	if err == redis.Nil {
		_, err = r1.Set(context.Background(), c.ClientIP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Result()
		value = os.Getenv("API_QUOTA")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
	} else {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
	}
	valInt, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errror": "error at converting"})
		return
	}
	if valInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Api Quota exceeded"})
		return
	}
	if req.Custom_short == "" {
		req.Custom_short = uuid.New().String()[:6]
	}
	r3 := database.RedisClient(1)
	defer r3.Close()

	if req.Expiry == 0 {
		req.Expiry = 24
	}

	val, err := r3.Get(context.Background(), req.Custom_short).Result()

	if err == redis.Nil {
		r3.Set(context.Background(), req.Custom_short, req.Url, req.Expiry*3600*time.Second).Result()
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Connectivity to database"})
		return
	}
	if val != "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "URl Already taken with that name"})
		return
	}

	resp := models.Response{
		URL:          "",
		Custom_short: "",
		Expiry:       0,
		RateLimit:    10,
	}
	_ = r1.Decr(context.Background(), c.ClientIP())
	vals, errs := r1.Get(context.Background(), c.ClientIP()).Result()
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "databse connectivity problem"})
		return
	}
	valInt, _ = strconv.Atoi(vals)
	resp.Custom_short = os.Getenv("DOMAIN") + "/" + req.Custom_short
	resp.Expiry = req.Expiry
	resp.RateLimit = valInt
	resp.URL = req.Url

	c.JSON(http.StatusOK, gin.H{"Data is": resp})
}

func DirectUrl(c *gin.Context) {
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

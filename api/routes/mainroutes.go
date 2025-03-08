package routes

import (
	"log"
	"net/http"
	"os"
	"tinyurl/auth"

	"github.com/gin-gonic/gin"
)

type Route struct {
	name       string
	method     string
	pattern    string
	handlefunc func(*gin.Context)
}

type Router []Route

type Routes struct {
	r *gin.Engine
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("content-type", "application/json")
		c.Writer.Header().Set("access-control-allow-origin", "*")
		c.Writer.Header().Set("access-control-allow-credentials", "true")
		c.Writer.Header().Set("access-control-allow-methods", "PUT,GET,POST,DELETE,OPTIONS,PATCH")
		c.Writer.Header().Set("access-control-allow-headers", "authorization,content-type,X-requested-width")
		if c.Request.Method == "OPTION" {
			c.Status(http.StatusNoContent)
			c.Abort()
			return
		}
	}
}

func (r *Routes) WebsiteHealthChecker(rg *gin.RouterGroup) {
	group := rg.Group("/urlshortner")
	group.Use(CORSMiddleware())
	for _, value := range healthcheck {
		switch value.method {
		case "GET":
			group.GET(value.pattern, value.handlefunc)

		case "POST":
			group.POST(value.pattern, value.handlefunc)

		case "PUT":
			group.PUT(value.pattern, value.handlefunc)

		case "OPTIONS":
			group.OPTIONS(value.pattern, value.handlefunc)

		case "DELETE":
			group.DELETE(value.pattern, value.handlefunc)

		default:
			group.GET(value.pattern, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Please select a correct method"})
			})
		}
	}
}

func (r *Routes) USER(rg *gin.RouterGroup) {
	group := rg.Group("/urlshortner")
	group.Use(CORSMiddleware())
	for _, value := range User {
		switch value.method {
		case "GET":
			group.GET(value.pattern, value.handlefunc)
		case "POST":
			group.POST(value.pattern, value.handlefunc)
		case "OPTIONS":
			group.OPTIONS(value.pattern, value.handlefunc)
		case "PUT":
			group.OPTIONS(value.pattern, value.handlefunc)
		case "DELETE":
			group.DELETE(value.pattern, value.handlefunc)
		default:
			group.GET(value.pattern, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Please provide a valid method"})
			})
		}
	}
}

func (r *Routes) URLShortner(rg *gin.RouterGroup) {
	group := rg.Group("/urlshortner")
	group.Use(CORSMiddleware())
	for _, value := range Shortner {
		switch value.method {
		case "GET":
			group.GET(value.pattern, value.handlefunc)
		case "POST":
			group.POST(value.pattern, value.handlefunc)
		case "OPTIONS":
			group.OPTIONS(value.pattern, value.handlefunc)
		case "PUT":
			group.PUT(value.pattern, value.handlefunc)
		case "DELETE":
			group.DELETE(value.pattern, value.handlefunc)
		default:
			group.GET(value.pattern, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "please enter correct method"})
			})
		}
	}
}

func Client() {
	routes := Routes{
		r: gin.Default(),
	}
	rg := routes.r.Group("/api/v1")
	routes.WebsiteHealthChecker(rg)
	routes.USER(rg)
	rg.Use(auth.Auth())
	routes.URLShortner(rg)
	log.Fatal(routes.r.Run(os.Getenv("APP_PORT")))
}

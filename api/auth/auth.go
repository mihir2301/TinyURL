package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
)

type JWTwrapper struct {
	SecretKey      string
	Issuer         string
	ExpirationTime int64
}

type JWTclaims struct {
	Email string
	jwt.StandardClaims
}

func (j *JWTwrapper) GenerateToken(email string) (token string, err error) {
	claims := &JWTclaims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: &jwt.Time{time.Now().Add(time.Hour * time.Duration(j.ExpirationTime))},
			Issuer:    j.Issuer,
		},
	}
	token1 := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = token1.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (j *JWTwrapper) ValidateToken(signedToken string) (claims *JWTclaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTclaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTclaims)
	if !ok {
		err = errors.New("could not parse claims")
		return nil, err
	}
	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		err = errors.New("token is expired")
		return
	}
	return claims, nil
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You Are Not Authorized"})
			c.Abort()
			return
		}
		extractedToken := strings.Split(token, "Bearer ")
		if len(extractedToken) == 2 {
			token = strings.TrimSpace(extractedToken[1])
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized"})
			c.Abort()
			return
		}
		JWTwrapper := JWTwrapper{
			SecretKey: os.Getenv("JwtSecrets"),
			Issuer:    os.Getenv("Jwtissuer"),
		}

		claims, err := JWTwrapper.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Next()
	}
}

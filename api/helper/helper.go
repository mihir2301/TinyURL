package helper

import (
	"errors"
	"os"
	"strings"
	"tinyurl/models"

	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
)

func CheckUrl(url string) bool {
	return govalidator.IsURL(url)
}

func EnforceUrl(url string) string {
	if url[:4] != "http" {
		return "https://" + url
	} else {
		return url
	}
}

func CheckDomain(url string) error {
	url = strings.Replace(url, "http://", "", 1)
	url = strings.Replace(url, "https://", "", 1)
	url = strings.Replace(url, "www.", "", 1)
	url = strings.Split(url, "/")[0]

	if url == os.Getenv("DOMAIN") {
		return errors.New("nice try diddy")
	}
	return nil
}
func GenPassHash(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(bytes)
}
func CheckDetails(user models.UserClient) error {
	if user.Email == "" {
		return errors.New("email caanot be empty")
	}
	if user.Password == "" {
		return errors.New("password cannot be empty")
	}
	if user.Phone == "" {
		return errors.New("Phone cannotbe Epmty")
	}
	if user.Name == "" {
		return errors.New("Name cannot be empty")
	}
	return nil
}

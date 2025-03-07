package helper

import (
	"errors"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
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

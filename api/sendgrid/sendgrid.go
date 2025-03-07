package sendgrid

import (
	"errors"
	"math/rand"
	"os"
	"strconv"
	"time"
	"tinyurl/constants"
	"tinyurl/models"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmails(req models.Verification) (models.Verification, error) {
	apikey := os.Getenv("SENDGRID_API_KEY")
	if apikey == "" {
		return req, errors.New("apikey of sendgrid is empty")
	}
	client := sendgrid.NewSendClient(apikey)

	from := mail.NewEmail("Sender Name", constants.Sender)
	to := mail.NewEmail("Recipient Name", req.Email)
	subject := "OTP verification mail"
	otp := Randomnum()
	req.Otp = strconv.Itoa(otp)
	htmlContent := "<p>Your login otp is <strong>" + req.Otp + "</strong><p>"
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	_, err := client.Send(message)
	if err != nil {
		return req, err
	}
	return req, nil
}
func Randomnum() int {
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(1000) + 1000
	return randomInt
}

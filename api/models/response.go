package models

import "time"

type Response struct {
	URL          string        `json:"url"`
	Custom_short string        `json:"short"`
	Expiry       time.Duration `json:"expiry"`
	RateLimit    int           `json:"limit"`
}

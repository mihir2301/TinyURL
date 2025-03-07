package models

import "time"

type Request struct {
	Url          string        `json:"url" binding:"required"`
	Custom_short string        `json:"short"`
	Expiry       time.Duration `json:"expiry"`
}

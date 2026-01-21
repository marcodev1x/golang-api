package structs

import "time"

type ById struct {
	Id string `json:"id" binding:"required"`
}

type CreateUrl struct {
	Url       string     `json:"url" binding:"required"`
	ExpiresAt *time.Time `json:"expires_at"`
}

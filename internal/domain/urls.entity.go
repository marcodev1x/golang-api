package domain

import "time"

type Urls struct {
	Id            int        `json:"id" gorm:"primaryKey"`
	HashedDomain  string     `json:"hashed_domain"`
	ShortenedUrl  string     `json:"shortened_url"`
	CreatedAt     time.Time  `json:"created_at"`
	ExpiresAt     *time.Time `json:"expires_at"`
	CountedClicks int        `json:"counted_clicks"`
}

func (u *Urls) TableName() string {
	return "urls"
}

package domain

import "time"

type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

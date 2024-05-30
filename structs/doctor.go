package structs

import "time"

type Doctor struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" binding:"required"`
	Password  string    `json:"password" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

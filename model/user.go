package model

import "time"

type User struct {
	Id        int       `json:"id"`
	Name      string    `json:"name" validate:"required,min=2,max=100"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=6"`
	Role      string    `json:"role" validate:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

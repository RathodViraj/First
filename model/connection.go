package model

import "time"

type Connection struct {
	FollowerID  int       `json:"follower_id" db:"follower_id" validate:"required"`
	FollowingID int       `json:"following_id" db:"following_id" validate:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

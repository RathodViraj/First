package model

import "time"

type Post struct {
	Id        int       `json:"id"`
	Uid       int       `json:"uid"`
	Content   string    `json:"content"`
	Likes     int       `json:"likes"`
	ParentId  *int      `json:"parent_id,omitempty"` // nil for top-level posts
	CreatedAt time.Time `json:"created_at"`
	Comments  []Post    `json:"comments,omitempty"` // Populated when fetching with joins
}

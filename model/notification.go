package model

type Notification struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	FromUser  int    `json:"from_user"`
	ToUser    int    `json:"to_user"`
	PostID    *int   `json:"post_id"`
	Message   string `json:"message"`
	Seen      bool   `json:"seen"`
	Timestamp int64  `json:"timestamp"`
}

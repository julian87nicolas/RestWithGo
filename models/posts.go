package models

type Post struct {
	Id          string `json:"id"`
	PostContent string `json:"post_content"`
	CreatedAt   string `json:"created_at"`
	UserId      string `json:"user_id"`
}



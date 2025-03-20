package models

// Post represents the post model
type Post struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Category     string `json:"category"`
	TextContent  string `json:"text_content"`
	DateCreation string `json:"date_creation"`
	UserID       int    `json:"user_id"`
	Likes        int    `json:"likes"`
	Dislikes     int    `json:"dislikes"`
	UserNickname string `json:"user_nickname,omitempty"` // For join queries
}

// CreatePostRequest captures the data needed to create a post
type CreatePostRequest struct {
	Title       string `json:"title"`
	Category    string `json:"category"`
	TextContent string `json:"text_content"`
}

// Comment represents the comment model
type Comment struct {
	ID           int    `json:"id"`
	TextContent  string `json:"text_content"`
	CreateDate   string `json:"create_date"`
	UserID       int    `json:"user_id"`
	PostID       int    `json:"post_id"`
	Likes        int    `json:"likes"`
	Dislikes     int    `json:"dislikes"`
	UserNickname string `json:"user_nickname,omitempty"` // For join queries
}

// CreateCommentRequest captures the data needed to create a comment
//
//	type CreateCommentRequest struct {
//		TextContent string `json:"text_content"`
//		PostID      int    `json:"post_id"`
//	}
//
//	type CreateCommentRequest struct {
//		TextContent string `json:"TextContent"`
//		PostID      int    `json:"PostID"`
//	}
type CreateCommentRequest struct {
	TextContent string `json:"text_content"` // Use "text_content" for JSON
	PostID      int    `json:"post_id"`      // Keep "PostID" as is
}

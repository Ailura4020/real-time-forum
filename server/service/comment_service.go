package service

import (
	"real-time-forum/models"
	"real-time-forum/repository"
	"time"
)

// CommentService provides methods for comment-related operations
type CommentService struct {
	CommentRepo *repository.CommentRepository
}

// NewCommentService creates a new CommentService
func NewCommentService(commentRepo *repository.CommentRepository) *CommentService {
	return &CommentService{CommentRepo: commentRepo}
}

// CreateComment handles comment creation
func (s *CommentService) CreateComment(req models.CreateCommentRequest, userID int) (models.Comment, error) {
	comment := models.Comment{
		TextContent: req.TextContent,
		PostID:      req.PostID,
		UserID:      userID,
		CreateDate:  time.Now().Format(time.RFC3339),
	}

	id, err := s.CommentRepo.CreateComment(comment)
	if err != nil {
		return models.Comment{}, err
	}

	comment.ID = int(id)
	return comment, nil
}

// GetCommentsByPostID retrieves all comments for a specific post
func (s *CommentService) GetCommentsByPostID(postID int) ([]models.Comment, error) {
	return s.CommentRepo.GetCommentsByPostID(postID)
}

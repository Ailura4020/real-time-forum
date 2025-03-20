package service

import (
	"real-time-forum/models"
	"real-time-forum/repository"
	"time"
)

// PostService provides methods for post-related operations
type PostService struct {
	PostRepo *repository.PostRepository
}

// NewPostService creates a new PostService
func NewPostService(postRepo *repository.PostRepository) *PostService {
	return &PostService{PostRepo: postRepo}
}

// CreatePost handles post creation
func (s *PostService) CreatePost(req models.CreatePostRequest, userID int) (models.Post, error) {
	post := models.Post{
		Title:        req.Title,
		Category:     req.Category,
		TextContent:  req.TextContent,
		UserID:       userID,
		DateCreation: time.Now().Format(time.RFC3339),
	}

	id, err := s.PostRepo.CreatePost(post)
	if err != nil {
		return models.Post{}, err
	}

	post.ID = int(id)
	return post, nil
}

// GetAllPosts retrieves all posts
func (s *PostService) GetAllPosts() ([]models.Post, error) {
	return s.PostRepo.GetAllPosts()
}

// GetPostByID retrieves a post by its ID
func (s *PostService) GetPostByID(id int) (models.Post, error) {
	return s.PostRepo.GetPostByID(id)
}

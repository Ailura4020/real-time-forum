package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-forum/models"
	"real-time-forum/repository"
	"real-time-forum/service"
	"real-time-forum/utils"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// GetPostsHandler returns a handler for retrieving all posts
func GetPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postRepo := repository.NewPostRepository(db)
		postService := service.NewPostService(postRepo)

		posts, err := postService.GetAllPosts()
		if err != nil {
			SendErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		SendResponse(w, true, "Posts retrieved successfully", posts, "")
	}
}

// GetPostHandler returns a handler for retrieving a specific post with its comments
func GetPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			SendErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		postRepo := repository.NewPostRepository(db)
		postService := service.NewPostService(postRepo)

		post, err := postService.GetPostByID(id)
		if err != nil {
			SendErrorResponse(w, err, http.StatusNotFound)
			return
		}

		commentRepo := repository.NewCommentRepository(db)
		commentService := service.NewCommentService(commentRepo)

		comments, err := commentService.GetCommentsByPostID(id)
		if err != nil {
			SendErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		response := struct {
			Post     models.Post      `json:"post"`
			Comments []models.Comment `json:"comments"`
		}{
			Post:     post,
			Comments: comments,
		}

		SendResponse(w, true, "Post retrieved successfully", response, "")
	}
}

// CreatePostHandler returns a handler for creating a new post
func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from JWT token
		userID, err := ExtractUserIDFromRequest(r)
		if err != nil {
			SendErrorResponse(w, err, http.StatusUnauthorized)
			return
		}

		var req models.CreatePostRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			SendErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		// Validate request
		if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Category) == "" {
			SendErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		postRepo := repository.NewPostRepository(db)
		postService := service.NewPostService(postRepo)

		post, err := postService.CreatePost(req, userID)
		if err != nil {
			SendErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		SendResponse(w, true, "Post created successfully", post, "")
	}
}

// ExtractUserIDFromRequest extracts the user ID from the request using the JWT token
func ExtractUserIDFromRequest(r *http.Request) (int, error) {
	// Get the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("authorization header is required")
	}

	// Split the header to get the token
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return 0, fmt.Errorf("invalid token format")
	}

	tokenString := splitToken[1]

	// Validate the token using the existing utils function
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		return 0, err
	}

	return claims.ID, nil
}

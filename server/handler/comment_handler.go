package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"real-time-forum/models"
	"real-time-forum/repository"
	"real-time-forum/service"
	"strings"
)

// AddCommentHandler returns a handler for adding a comment to a post
func AddCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from JWT token
		userID, err := ExtractUserIDFromRequest(r)
		if err != nil {
			SendErrorResponse(w, err, http.StatusUnauthorized)
			return
		}

		var req models.CreateCommentRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			SendErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		// Validate request
		//if strings.TrimSpace(req.TextContent) == "" || req.PostID <= 0 {
		//	SendErrorResponse(w, err, http.StatusBadRequest)
		//	return
		//}

		log.Printf("Decoded request: %+v", req)

		if strings.TrimSpace(req.TextContent) == "" || req.PostID <= 0 {
			SendErrorResponse(w, errors.New("invalid request: text content is required and post ID must be greater than 0"), http.StatusBadRequest)
			return
		}

		// Verify that the post exists
		postRepo := repository.NewPostRepository(db)
		postService := service.NewPostService(postRepo)

		_, err = postService.GetPostByID(req.PostID)
		if err != nil {
			SendErrorResponse(w, err, http.StatusNotFound)
			return
		}

		commentRepo := repository.NewCommentRepository(db)
		commentService := service.NewCommentService(commentRepo)

		comment, err := commentService.CreateComment(req, userID)
		if err != nil {
			SendErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		SendResponse(w, true, "Comment added successfully", comment, "")
	}
}

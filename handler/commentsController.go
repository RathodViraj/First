package handler

import (
	chachingservice "First/chachingService"
	"First/model"
	"First/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentsHandler struct {
	service *service.CommentsService
}

func NewCommentsHandler(service *service.CommentsService) *CommentsHandler {
	return &CommentsHandler{service}
}

func (h *CommentsHandler) AddComment(ctx *gin.Context) {
	uidVal, exists := ctx.Get("userID")
	if !exists {
		JSONError(ctx, http.StatusUnauthorized, "User not authenticated")
		return
	}
	userID, ok := uidVal.(int)
	if !ok {
		JSONError(ctx, http.StatusInternalServerError, "Failed to extract user ID")
		return
	}

	postID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var comment model.Post
	if err := ctx.BindJSON(&comment); err != nil {
		JSONError(ctx, http.StatusBadRequest, "Invalid request")
		return
	}

	comment.Uid = userID
	comment.ParentId = &postID

	if err := h.service.CreateComment(&comment); err != nil {
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.IndentedJSON(http.StatusCreated, comment)

}

func (h *CommentsHandler) GetAllComments(ctx *gin.Context) {
	_, exists := ctx.Get("userID")
	if !exists {
		JSONError(ctx, http.StatusUnauthorized, "User not authenticated")
		return
	}

	postID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "Invalid post ID")
		return
	}

	comments, err := h.service.GetComments(postID)
	if err != nil {
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.IndentedJSON(http.StatusOK, comments)
}

func (h *CommentsHandler) UpdateComment(ctx *gin.Context) {
	postID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "Invalid post ID")
		return
	}
	commentID, err := strconv.Atoi(ctx.Param("comment_id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "Invalid comment ID")
		return
	}
	uidVal, exists := ctx.Get("userID")
	if !exists {
		JSONError(ctx, http.StatusUnauthorized, "User not authenticated")
		return
	}
	userID, ok := uidVal.(int)
	if !ok {
		JSONError(ctx, http.StatusInternalServerError, "Failed to extract user ID")
		return
	}

	var comment model.Post
	if err := ctx.BindJSON(&comment); err != nil {
		JSONError(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}
	comment.Id = commentID
	comment.Uid = userID
	comment.ParentId = &postID

	if err := h.service.UpdateComment(&comment); err != nil {
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	chachingservice.InvalUserIDateUserProfileChahe(userID, ctx)
	ctx.Status(http.StatusOK)
}

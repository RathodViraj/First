package handler

import (
	"First/model"
	"First/notification"
	"First/repository"
	"First/service"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentsHandler struct {
	service *service.CommentsService
	hub     *notification.Hub
	gr      repository.Graph
	nr      repository.NotificationRepository
}

func NewCommentsHandler(service *service.CommentsService, h *notification.Hub, gr repository.Graph, nr repository.NotificationRepository) *CommentsHandler {
	return &CommentsHandler{
		service,
		h,
		gr,
		nr,
	}
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

	ownerID, _ := h.gr.GetUserIDByPostID(postID)
	if ownerID == -1 || ownerID == userID {
		ctx.IndentedJSON(http.StatusCreated, comment)
		return
	}

	notif := model.Notification{
		Type:      "commnet",
		FromUser:  userID,
		ToUser:    ownerID,
		PostID:    &postID,
		Message:   "Someone comment on your post",
		Timestamp: time.Now().Unix(),
	}
	h.hub.Broadcast <- notif
	log.Println("Sending notification via hub")

	_ = h.nr.SaveNotification(model.Notification{
		Type:      "comment",
		FromUser:  userID,
		ToUser:    ownerID,
		PostID:    &postID,
		Message:   "Someone comment on your post",
		Seen:      false,
		Timestamp: time.Now().Unix(),
	})

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

	ctx.Status(http.StatusOK)
}
